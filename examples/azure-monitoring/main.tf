## This code demonstrate how to setup Azure Redis Cache monitors


provider "azurerm" {
  version = "=2.0.0"
  features {}
}


### Cache Hits Alert
resource "azurerm_monitor_metric_alert" "cache_hit_alert" {
  name                = "${var.cache.service_name} ${var.cache.environment} - Cache Hits Alert"
  resource_group_name = var.cache.cache_name
  scopes              = [var.cache.scope]
  description         = "${var.cache.service_name} Cache Hits Alert"

  criteria {
    metric_namespace = "Microsoft.Cache/redis"
    metric_name      = "cachehits"
    aggregation      = "Total"
    operator         = "GreaterThan"
    threshold        = var.cache.cache_hit_threshold

    
  }

  action {
    action_group_id = var.action_group
  }
}


### Cache Misses Alert
resource "azurerm_monitor_metric_alert" "cache_miss_alert" {
  name                = "${var.cache.service_name} ${var.cache.environment} - Cache Miss Alert"
  resource_group_name = var.cache.cache_name
  scopes              = [var.cache.scope]
  description         = "${var.cache.service_name} - Cache Miss Alert"

  criteria {
    metric_namespace = "Microsoft.Cache/redis"
    metric_name      = "cachemisses"
    aggregation      = "Total"
    operator         = "GreaterThan"
    threshold        = var.cache.cache_misses_threshold

    
  }

  action {
    action_group_id = var.action_group
  }
}

### Cache Connection Alert
resource "azurerm_monitor_metric_alert" "cache_connected_clients" {
  name                = "${var.cache.service_name} ${var.cache.environment} - Cache Connected Clients"
  resource_group_name = var.cache.cache_name
  scopes              = [var.cache.scope]
  description         = "${var.cache.service_name} - Cache Connected Clients"

  criteria {
    metric_namespace = "Microsoft.Cache/redis"
    metric_name      = "connectedclients"
    aggregation      = "Total"
    operator         = "GreaterThan"
    threshold        = var.cache.cache_connected_clients_threshold

    
  }

  action {
    action_group_id = var.action_group
  }
}


### Cache CPU Alert
resource "azurerm_monitor_metric_alert" "cache_cpu" {
  name                = "${var.cache.service_name} ${var.cache.environment} - Cache CPU"
  resource_group_name = var.cache.cache_name
  scopes              = [var.cache.scope]
  description         = "${var.cache.service_name} - Cache CPU"

  criteria {
    metric_namespace = "Microsoft.Cache/redis"
    metric_name      = "percentProcessorTime"
    aggregation      = "Total"
    operator         = "GreaterThan"
    threshold        = var.cache.cache_cpu_threshold

    
  }

  action {
    action_group_id = var.action_group
  }
}