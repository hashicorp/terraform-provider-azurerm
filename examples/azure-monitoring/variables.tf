variable "cache" {
	description = "The Azure Redis Cache alerts"

	default = {
      		cache_name                       = "<replace this with cache name>"
      		service_name					= "<replace this project name>"
      		environment						= "<Stage/Production>"
			scope                   = "/subscriptions/<subscription_id>/resourceGroups/<resource_group_name>/providers/Microsoft.Cache/Redis/<azure_redis_cache_name>"
			cache_hit_threshold          = 5000
			cache_misses_threshold    = 5000
			cache_cpu_threshold = 50
			cache_connected_clients_threshold       = 5000
		
    }

	}

variable "action_group" {
  description = "The azure action group"
  type        = string
  default = "/subscriptions/269e5205-71e4-43fc-b2cf-a2f76843e431/resourceGroups/<resource_group_name>/providers/microsoft.insights/actionGroups/<action_group_name>"
}