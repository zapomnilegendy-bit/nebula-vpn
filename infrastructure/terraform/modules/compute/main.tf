locals {
  environment = var.environment
  name_prefix = "nebula-${local.environment}"
}

# ECS Cluster
resource "aws_ecs_cluster" "main" {
  name = local.name_prefix
  setting { name = "containerInsights", value = "enabled" }
  tags = { Environment = local.environment }
}

# CloudWatch Log Group
resource "aws_cloudwatch_log_group" "ecs" {
  name              = "/ecs/${local.name_prefix}"
  retention_in_days = 30
}

# Service Discovery
resource "aws_service_discovery_private_dns_namespace" "main" {
  name        = "nebula.internal"
  vpc         = var.vpc_id
  description = "NebulaVPN internal service discovery"
}

# ALB
resource "aws_lb" "main" {
  name               = "${local.name_prefix}-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets            = var.public_subnet_ids
  enable_deletion_protection = var.environment == "production"
  tags               = { Environment = local.environment }
}

resource "aws_security_group" "alb" {
  name        = "${local.name_prefix}-alb"
  description = "ALB security group"
  vpc_id      = var.vpc_id
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
  tags = { Environment = local.environment }
}

resource "aws_lb_target_group" "api" {
  name        = "${local.name_prefix}-api"
  port        = 8080
  protocol    = "HTTP"
  vpc_id      = var.vpc_id
  target_type = "ip"
  health_check {
    enabled             = true
    path                = "/v1/health"
    port                = 8080
    healthy_threshold   = 2
    unhealthy_threshold = 3
    interval            = 30
    timeout             = 5
    matcher             = "200"
  }
  tags = { Environment = local.environment }
}

resource "aws_lb_target_group" "user_panel" {
  name        = "${local.name_prefix}-user"
  port        = 3000
  protocol    = "HTTP"
  vpc_id      = var.vpc_id
  target_type = "ip"
  health_check {
    path                = "/"
    healthy_threshold   = 2
    unhealthy_threshold = 3
    interval            = 30
    matcher             = "200"
  }
  tags = { Environment = local.environment }
}

resource "aws_lb_target_group" "admin_panel" {
  name        = "${local.name_prefix}-admin"
  port        = 3000
  protocol    = "HTTP"
  vpc_id      = var.vpc_id
  target_type = "ip"
  health_check {
    path                = "/"
    healthy_threshold   = 2
    unhealthy_threshold = 3
    interval            = 30
    matcher             = "200"
  }
  tags = { Environment = local.environment }
}

resource "aws_lb_listener" "https" {
  load_balancer_arn = aws_lb.main.arn
  port              = 443
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-TLS13-1-2-2021-06"
  certificate_arn   = var.certificate_arn
  default_action {
    type = "fixed-response"
    fixed_response {
      content_type = "application/json"
      status_code  = "404"
      message_body = "{\"error\":\"not_found\"}"
    }
  }
}

resource "aws_lb_listener_rule" "api" {
  listener_arn = aws_lb_listener.https.arn
  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.api.arn
  }
  condition { path_pattern { values = ["/v1/*", "/api/*"] } }
}

resource "aws_lb_listener_rule" "user_panel" {
  listener_arn = aws_lb_listener.https.arn
  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.user_panel.arn
  }
  condition { host_header { values = ["app.${var.domain_name}"] } }
}

resource "aws_lb_listener_rule" "admin_panel" {
  listener_arn = aws_lb_listener.https.arn
  action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.admin_panel.arn
  }
  condition { host_header { values = ["admin.${var.domain_name}"] } }
}

# HTTP redirect to HTTPS
resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.main.arn
  port              = 80
  protocol          = "HTTP"
  default_action {
    type = "redirect"
    redirect {
      port        = "443"
      protocol    = "HTTPS"
      status_code = "HTTP_301"
    }
  }
}

# ECS Task Definitions
resource "aws_ecs_task_definition" "api" {
  family                   = "${local.name_prefix}-api"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = var.ecs_cpu
  memory                   = var.ecs_memory
  execution_role_arn       = aws_iam_role.ecs_execution.arn
  task_role_arn            = aws_iam_role.ecs_task.arn

  container_definitions = jsonencode([
    {
      name  = "api-gateway"
      image = var.api_container_image
      portMappings = [{ containerPort = 8080, protocol = "tcp" }]
      environment = [
        { name = "NEBULA_ENV", value = local.environment },
        { name = "NEBULA_DB_HOST", value = var.db_host },
        { name = "NEBULA_DB_PORT", value = tostring(var.db_port) },
        { name = "NEBULA_DB_NAME", value = var.db_name },
        { name = "NEBULA_DB_USER", value = var.db_username },
        { name = "NEBULA_REDIS_ADDR", value = "${var.redis_host}:${var.redis_port}" },
        { name = "NEBULA_JWT_SECRET", value = "from-ssm:/nebula/${local.environment}/jwt-secret" },
      ]
      secrets = [
        { name = "NEBULA_DB_PASSWORD", valueFrom = "arn:aws:ssm:us-east-1:ACCOUNT:parameter/nebula/${local.environment}/db-password" },
        { name = "STRIPE_SECRET_KEY", valueFrom = "arn:aws:ssm:us-east-1:ACCOUNT:parameter/nebula/${local.environment}/stripe-secret" },
        { name = "STRIPE_WEBHOOK_SECRET", valueFrom = "arn:aws:ssm:us-east-1:ACCOUNT:parameter/nebula/${local.environment}/stripe-webhook" },
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.ecs.name
          "awslogs-region"        = "us-east-1"
          "awslogs-stream-prefix" = "api"
        }
      }
      healthCheck = {
        command     = ["CMD-SHELL", "curl -f http://localhost:8080/v1/health || exit 1"]
        interval    = 30
        timeout     = 5
        retries     = 3
      }
    }
  ])
  tags = { Environment = local.environment }
}

resource "aws_ecs_task_definition" "user_panel" {
  family                   = "${local.name_prefix}-user-panel"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = 256
  memory                   = 512
  execution_role_arn       = aws_iam_role.ecs_execution.arn

  container_definitions = jsonencode([
    {
      name  = "user-panel"
      image = var.user_panel_image
      portMappings = [{ containerPort = 3000, protocol = "tcp" }]
      environment = [
        { name = "NEXT_PUBLIC_API_URL", value = "https://api.${var.domain_name}" },
        { name = "NEXT_PUBLIC_STRIPE_KEY", value = "pk_live_..." },
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.ecs.name
          "awslogs-region"        = "us-east-1"
          "awslogs-stream-prefix" = "user-panel"
        }
      }
    }
  ])
  tags = { Environment = local.environment }
}

resource "aws_ecs_task_definition" "admin_panel" {
  family                   = "${local.name_prefix}-admin-panel"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = 256
  memory                   = 512
  execution_role_arn       = aws_iam_role.ecs_execution.arn

  container_definitions = jsonencode([
    {
      name  = "admin-panel"
      image = var.admin_panel_image
      portMappings = [{ containerPort = 3000, protocol = "tcp" }]
      environment = [
        { name = "NEXT_PUBLIC_API_URL", value = "https://api.${var.domain_name}" },
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.ecs.name
          "awslogs-region"        = "us-east-1"
          "awslogs-stream-prefix" = "admin-panel"
        }
      }
    }
  ])
  tags = { Environment = local.environment }
}

# ECS Services
resource "aws_ecs_service" "api" {
  name                              = "${local.name_prefix}-api"
  cluster                           = aws_ecs_cluster.main.id
  task_definition                   = aws_ecs_task_definition.api.arn
  desired_count                     = var.desired_count
  launch_type                       = "FARGATE"
  platform_version                  = "LATEST"
  health_check_grace_period_seconds = 30

  network_configuration {
    subnets          = var.private_subnet_ids
    security_groups  = var.security_group_ids
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.api.arn
    container_name   = "api-gateway"
    container_port   = 8080
  }

  service_connect_configuration {
    enabled   = true
    namespace = aws_service_discovery_private_dns_namespace.main.arn
    service {
      discovery_name = "api"
      port_name      = "api-gateway"
      client_alias {
        port     = 8080
        dns_name = "api.nebula.internal"
      }
    }
  }

  deployment_circuit_breaker {
    enable   = true
    rollback = true
  }

  tags = { Environment = local.environment }
}

resource "aws_ecs_service" "user_panel" {
  name            = "${local.name_prefix}-user-panel"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.user_panel.arn
  desired_count   = var.desired_count
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = var.private_subnet_ids
    security_groups  = var.security_group_ids
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.user_panel.arn
    container_name   = "user-panel"
    container_port   = 3000
  }

  tags = { Environment = local.environment }
}

resource "aws_ecs_service" "admin_panel" {
  name            = "${local.name_prefix}-admin-panel"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.admin_panel.arn
  desired_count   = var.desired_count
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = var.private_subnet_ids
    security_groups  = var.security_group_ids
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.admin_panel.arn
    container_name   = "admin-panel"
    container_port   = 3000
  }

  tags = { Environment = local.environment }
}

# IAM
resource "aws_iam_role" "ecs_execution" {
  name = "${local.name_prefix}-ecs-execution"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Principal = { Service = "ecs-tasks.amazonaws.com" }
      Action = "sts:AssumeRole"
    }]
  })
  managed_policy_arns = [
    "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy",
    "arn:aws:iam::aws:policy/CloudWatchLogsFullAccess",
    "arn:aws:iam::aws:policy/AmazonSSMReadOnlyAccess",
  ]
}

resource "aws_iam_role" "ecs_task" {
  name = "${local.name_prefix}-ecs-task"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Principal = { Service = "ecs-tasks.amazonaws.com" }
      Action = "sts:AssumeRole"
    }]
  })
}

# Auto Scaling
resource "aws_appautoscaling_target" "api" {
  service_namespace  = "ecs"
  resource_id        = "service/${aws_ecs_cluster.main.name}/${aws_ecs_service.api.name}"
  scalable_dimension = "ecs:service:DesiredCount"
  min_capacity       = var.desired_count
  max_capacity       = 10
}

resource "aws_appautoscaling_policy" "api_cpu" {
  name               = "${local.name_prefix}-api-cpu"
  service_namespace  = aws_appautoscaling_target.api.service_namespace
  resource_id        = aws_appautoscaling_target.api.resource_id
  scalable_dimension = aws_appautoscaling_target.api.scalable_dimension
  policy_type        = "TargetTrackingScaling"
  target_tracking_scaling_policy_configuration {
    predefined_metric_specification {
      predefined_metric_type = "ECSServiceAverageCPUUtilization"
    }
    target_value       = 70
    scale_in_cooldown  = 300
    scale_out_cooldown = 60
  }
}

resource "aws_appautoscaling_policy" "api_memory" {
  name               = "${local.name_prefix}-api-memory"
  service_namespace  = aws_appautoscaling_target.api.service_namespace
  resource_id        = aws_appautoscaling_target.api.resource_id
  scalable_dimension = aws_appautoscaling_target.api.scalable_dimension
  policy_type        = "TargetTrackingScaling"
  target_tracking_scaling_policy_configuration {
    predefined_metric_specification {
      predefined_metric_type = "ECSServiceAverageMemoryUtilization"
    }
    target_value       = 80
    scale_in_cooldown  = 300
    scale_out_cooldown = 60
  }
}

# Outputs
output "alb_dns_name" { value = aws_lb.main.dns_name }
output "ecs_cluster_name" { value = aws_ecs_cluster.main.name }
