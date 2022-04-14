## This code demonstrate how to setup Azure Redis Cache monitors

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "main" {
  name     = "example_rg_rediscache"
  location = "east us"
}

resource "azurerm_monitor_action_group" "main" {
  name                = "example-actiongroup"
  resource_group_name = azurerm_resource_group.main.name
  short_name          = "exampleact"

  email_receiver {
    name                    = "ishantdevops"
    email_address           = "devops@example.com"
    use_common_alert_schema = true
  }

  webhook_receiver {
    name        = "callmyapi"
    service_uri = "http://example.com/alert"
  }
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
    action_group_id = azurerm_monitor_action_group.main.id
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
    action_group_id = azurerm_monitor_action_group.main.id
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
    action_group_id = azurerm_monitor_action_group.main.id
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
    action_group_id = azurerm_monitor_action_group.main.id
  }
}
