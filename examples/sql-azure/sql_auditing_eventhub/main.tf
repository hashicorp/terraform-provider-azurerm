# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location
}

resource "azurerm_mssql_server" "example" {
  name                         = "${var.prefix}-server-primary"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = var.location
  version                      = "12.0"
  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"
}

resource "azurerm_mssql_database" "example" {
  name      = "${var.prefix}-db-primary"
  server_id = azurerm_mssql_server.example.id
}

resource "azurerm_eventhub_namespace" "example" {
  name                = "${var.prefix}-EHN"
  location            = var.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"
}

resource "azurerm_eventhub" "example" {
  name                = "${var.prefix}-EH"
  namespace_name      = azurerm_eventhub_namespace.example.name
  resource_group_name = azurerm_resource_group.example.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_eventhub_namespace_authorization_rule" "example" {
  name                = "${var.prefix}EHRule"
  namespace_name      = azurerm_eventhub_namespace.example.name
  resource_group_name = azurerm_resource_group.example.name
  listen              = true
  send                = true
  manage              = true
}

resource "azurerm_monitor_diagnostic_setting" "example" {
  name                           = "${var.prefix}-DS"
  target_resource_id             = azurerm_mssql_database.example.id
  eventhub_authorization_rule_id = azurerm_eventhub_namespace_authorization_rule.example.id
  eventhub_name                  = azurerm_eventhub.example.name

  enabled_log {
    category = "SQLSecurityAuditEvents"
  }

  metric {
    category = "AllMetrics"
  }

  lifecycle {
    ignore_changes = [metric]
  }
}

resource "azurerm_mssql_database_extended_auditing_policy" "example" {
  database_id            = azurerm_mssql_database.example.id
  log_monitoring_enabled = true
}

resource "azurerm_mssql_server_extended_auditing_policy" "example" {
  server_id              = azurerm_mssql_server.example.id
  log_monitoring_enabled = true
}