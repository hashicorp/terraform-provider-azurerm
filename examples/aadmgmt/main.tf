provider "azurerm" {
  features {}
}


resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources2"
  location = var.location
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "${var.prefix}-la"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_storage_account" "example" {
  name                     = "${var.prefix}stgacc"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_eventhub_namespace" "example" {
  name                = "${var.prefix}eventhubns"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"
  capacity            = 1

  tags = {
    environment = "Production"
  }
}

resource "azurerm_eventhub" "example" {
  name                = "${var.prefix}eventhub"
  namespace_name      = azurerm_eventhub_namespace.example.name
  resource_group_name = azurerm_resource_group.example.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_aad_diagnostic_settings" "example" {
    name = "${var.prefix}-diag-settings"
    storage_account_id = azurerm_storage_account.example.id
    workspace_id = azurerm_log_analytics_workspace.example.id
    event_hub_name = azurerm_eventhub.example.name
    event_hub_auth_rule_id = "${azurerm_eventhub_namespace.example.id}/authorizationRules/RootManageSharedAccessKey"
    logs  {
        category = "AuditLogs"
        enabled = false
    }

    logs  {
        enabled = true
        category = "SignInLogs"
        retention_policy {
        retention_policy_days = 20
        retention_policy_enabled = true
        }
    }
}


