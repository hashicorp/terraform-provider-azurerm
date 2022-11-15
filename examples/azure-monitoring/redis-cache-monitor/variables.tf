variable "cache" {
  description = "The Azure Redis Cache alerts"

  default = {
    cache_name                        = "<replace this with cache name>"
    service_name                      = "<replace this project name>"
    environment                       = "<Stage/Production>"
    scope                             = "/subscriptions/<subscription_id>/resourceGroups/<resource_group_name>/providers/Microsoft.Cache/Redis/<azure_redis_cache_name>"
    cache_hit_threshold               = 5000
    cache_misses_threshold            = 5000
    cache_cpu_threshold               = 50
    cache_connected_clients_threshold = 5000

  }

}
