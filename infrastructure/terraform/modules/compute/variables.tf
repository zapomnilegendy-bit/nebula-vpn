variable "environment" { type = string }
variable "vpc_id" { type = string }
variable "public_subnet_ids" { type = list(string) }
variable "private_subnet_ids" { type = list(string) }
variable "certificate_arn" { type = string }
variable "domain_name" { type = string }
variable "db_host" { type = string }
variable "db_port" { type = number }
variable "db_name" { type = string }
variable "db_username" { type = string }
variable "db_password" { type = string; sensitive = true }
variable "redis_host" { type = string }
variable "redis_port" { type = number }
variable "api_container_image" { type = string }
variable "user_panel_image" { type = string }
variable "admin_panel_image" { type = string }
variable "allowed_cidrs" { type = list(string) }
variable "security_group_ids" { type = list(string) }
variable "ecs_cpu" { type = number; default = 1024 }
variable "ecs_memory" { type = number; default = 2048 }
variable "desired_count" { type = number; default = 2 }
