variable "environment" { type = string }
variable "vpc_id" { type = string }
variable "subnet_ids" { type = list(string) }
variable "db_instance_class" { type = string }
variable "db_allocated_storage" { type = number }
variable "db_name" { type = string }
variable "db_username" { type = string }
variable "db_password" { type = string; sensitive = true }
variable "security_group_ids" { type = list(string) }

resource "aws_db_subnet_group" "main" {
  name       = "nebula-${var.environment}"
  subnet_ids = var.subnet_ids
  tags       = { Name = "nebula-${var.environment}", Environment = var.environment }
}

resource "aws_security_group" "rds" {
  name        = "nebula-${var.environment}-rds"
  description = "RDS security group"
  vpc_id      = var.vpc_id
  ingress {
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = var.security_group_ids
  }
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
  tags = { Name = "nebula-${var.environment}-rds", Environment = var.environment }
}

resource "random_password" "db" {
  length  = 16
  special = false
}

resource "aws_db_instance" "main" {
  identifier             = "nebula-${var.environment}"
  engine                 = "postgres"
  engine_version         = "16"
  instance_class         = var.db_instance_class
  allocated_storage      = var.db_allocated_storage
  storage_type           = "gp3"
  storage_encrypted      = true
  db_name                = var.db_name
  username               = var.db_username
  password               = var.db_password != "" ? var.db_password : random_password.db.result
  db_subnet_group_name   = aws_db_subnet_group.main.name
  vpc_security_group_ids = [aws_security_group.rds.id]
  backup_retention_period = 30
  backup_window          = "03:00-04:00"
  maintenance_window     = "sun:04:00-sun:05:00"
  skip_final_snapshot    = false
  final_snapshot_identifier = "nebula-${var.environment}-final-${formatdate("YYYY-MM-DD-hhmm", timestamp())}"
  deletion_protection    = var.environment == "production"
  performance_insights_enabled = true
  monitoring_interval    = 60
  tags                   = { Name = "nebula-${var.environment}", Environment = var.environment }
}

output "db_host" { value = aws_db_instance.main.address }
output "db_port" { value = aws_db_instance.main.port }
output "db_name" { value = aws_db_instance.main.db_name }
output "db_username" { value = aws_db_instance.main.username }
