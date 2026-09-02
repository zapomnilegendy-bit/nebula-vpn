terraform {
  required_version = ">= 1.6"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
  backend "s3" {
    bucket = "nebula-terraform-state"
    key    = "staging/terraform.tfstate"
    region = "us-east-1"
  }
}

provider "aws" {
  region = var.aws_region
}

module "networking" {
  source = "../../modules/networking"
  environment = var.environment
  vpc_cidr = var.vpc_cidr
  availability_zones = var.availability_zones
  public_subnet_cidrs = var.public_subnet_cidrs
  private_subnet_cidrs = var.private_subnet_cidrs
}

module "database" {
  source = "../../modules/database"
  environment = var.environment
  vpc_id = module.networking.vpc_id
  subnet_ids = module.networking.private_subnet_ids
  db_instance_class = "db.t4g.small"
  db_allocated_storage = 20
  db_name = var.db_name
  db_username = var.db_username
  db_password = var.db_password
  security_group_ids = [module.networking.default_security_group_id]
}

module "cache" {
  source = "../../modules/cache"
  environment = var.environment
  vpc_id = module.networking.vpc_id
  subnet_ids = module.networking.private_subnet_ids
  redis_node_type = "cache.t4g.micro"
  redis_num_cache_nodes = 1
  security_group_ids = [module.networking.default_security_group_id]
}

module "compute" {
  source = "../../modules/compute"
  environment = var.environment
  vpc_id = module.networking.vpc_id
  public_subnet_ids = module.networking.public_subnet_ids
  private_subnet_ids = module.networking.private_subnet_ids
  certificate_arn = var.certificate_arn
  domain_name = var.domain_name
  db_host = module.database.db_host
  db_port = module.database.db_port
  db_name = var.db_name
  db_username = var.db_username
  db_password = var.db_password
  redis_host = module.cache.redis_host
  redis_port = module.cache.redis_port
  api_container_image = var.api_container_image
  user_panel_image = var.user_panel_image
  admin_panel_image = var.admin_panel_image
  allowed_cidrs = var.allowed_cidrs
  security_group_ids = [module.networking.default_security_group_id]
  ecs_cpu = 512
  ecs_memory = 1024
  desired_count = 1
}
