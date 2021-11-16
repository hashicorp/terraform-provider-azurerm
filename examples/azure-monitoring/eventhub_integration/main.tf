resource "azurerm_resource_group" "example" {
  name     = "example_rg_eventhub"
  location = "West Europe"
}

resource "azurerm_eventhub_namespace" "example" {
  name                = "exampleEventhubns"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku = "Standard"
}

resource "azurerm_eventhub" "example" {
  name                = "example_eventhub"
  namespace_name      = azurerm_eventhub_namespace.example.name
  resource_group_name = azurerm_resource_group.example.name

  partition_count   = 2
  message_retention = 1
}

resource "azurerm_eventhub_namespace_authorization_rule" "example" {
  name                = "example_eventhubns_auth_rule"
  namespace_name      = azurerm_eventhub_namespace.example.name
  resource_group_name = azurerm_resource_group.example.name

  listen = true
  send   = true
  manage = true
}

data "azurerm_kusto_cluster" "example" {
  name                = "examplekustocluster"
  resource_group_name = "example_rg_kusto_cluster"
}

data "azurerm_monitor_diagnostic_categories" "example" {
  resource_id = data.azurerm_kusto_cluster.example.id
}

resource "azurerm_monitor_diagnostic_setting" "example" {
  name                           = "example_monitor_diag_setting"
  target_resource_id             = data.azurerm_kusto_cluster.example.id
  eventhub_name                  = azurerm_eventhub.example.name
  eventhub_authorization_rule_id = azurerm_eventhub_namespace_authorization_rule.example.id

  dynamic "log" {
    for_each = data.azurerm_monitor_diagnostic_categories.example.logs
    content {
      category = log.key

      retention_policy {
        enabled = false
        days    = 0
      }
    }
  }
  dynamic "metric" {
    for_each = data.azurerm_monitor_diagnostic_categories.example.metrics

    content {
      category = metric.key

      retention_policy {
        enabled = false
        days    = 0
      }
    }
  }
}
