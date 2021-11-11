## Example: Azure Monitor Diagnostic Setting

This example provisions an Azure Monitor Diagnostic Setting using event hub as the destination

## Example Usage of azurerm_monitor_diagnostic_setting integration with eventhub

!> **Note:** Azure Monitor (Diagnostic Settings) can't access Event Hubs resources when virtual networks are enabled. You have to enable the allow trusted Microsoft services to bypass this firewall setting in Event Hub, so that Azure Monitor (Diagnostic Settings) service is granted access to your Event Hubs resources.


```hcl
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
```