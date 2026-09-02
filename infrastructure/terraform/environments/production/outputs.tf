output "vpc_id" {
  value = module.networking.vpc_id
}

output "db_host" {
  value = module.database.db_host
}

output "redis_host" {
  value = module.cache.redis_host
}

output "alb_dns_name" {
  value = module.compute.alb_dns_name
}

output "ecs_cluster_name" {
  value = module.compute.ecs_cluster_name
}
