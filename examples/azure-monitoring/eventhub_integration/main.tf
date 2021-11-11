resource "azurerm_resource_group" "example" {
  name     = "example_rg_eventhub"
  location = "West Europe"
}

resource "azurerm_eventhub_namespace" "example" {
  name                = "example_eventhubns"
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
  name                = "example_kustocluster"
  resource_group_name = "example_rg_kusto_cluster"
}

resource "azurerm_monitor_diagnostic_setting" "example" {
  name                           = "example_monitor_diag_setting"
  target_resource_id             = data.azurerm_kusto_cluster.example.id
  eventhub_name                  = azurerm_eventhub.example.name
  eventhub_authorization_rule_id = azurerm_eventhub_namespace_authorization_rule.example.id

  log {
    category = "Journal"
    enabled  = false
    retention_policy {
      enabled = false
    }
  }
  metric {
    category = "AllMetrics"
    retention_policy {
      enabled = false
    }
  }
}