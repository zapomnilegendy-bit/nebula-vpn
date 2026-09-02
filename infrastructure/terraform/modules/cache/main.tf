variable "environment" { type = string }
variable "vpc_id" { type = string }
variable "subnet_ids" { type = list(string) }
variable "redis_node_type" { type = string }
variable "redis_num_cache_nodes" { type = number }
variable "security_group_ids" { type = list(string) }

resource "aws_security_group" "redis" {
  name        = "nebula-${var.environment}-redis"
  description = "Redis security group"
  vpc_id      = var.vpc_id
  ingress {
    from_port       = 6379
    to_port         = 6379
    protocol        = "tcp"
    security_groups = var.security_group_ids
  }
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
  tags = { Name = "nebula-${var.environment}-redis", Environment = var.environment }
}

resource "aws_elasticache_subnet_group" "main" {
  name       = "nebula-${var.environment}"
  subnet_ids = var.subnet_ids
}

resource "aws_elasticache_replication_group" "main" {
  replication_group_id          = "nebula-${var.environment}"
  description                   = "NebulaVPN Redis ${var.environment}"
  node_type                     = var.redis_node_type
  num_cache_clusters            = var.redis_num_cache_nodes
  port                          = 6379
  parameter_group_name          = "default.redis7"
  automatic_failover_enabled    = var.redis_num_cache_nodes > 1
  multi_az_enabled              = var.redis_num_cache_nodes > 1
  subnet_group_name             = aws_elasticache_subnet_group.main.name
  security_group_ids            = [aws_security_group.redis.id]
  at_rest_encryption_enabled    = true
  transit_encryption_enabled    = true
  auto_minor_version_upgrade    = true
  tags                          = { Name = "nebula-${var.environment}", Environment = var.environment }
}

output "redis_host" { value = aws_elasticache_replication_group.main.primary_endpoint_address }
output "redis_port" { value = aws_elasticache_replication_group.main.port }
