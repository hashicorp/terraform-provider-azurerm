# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "${var.prefix}-laworkspace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
}


resource "azurerm_network_function_azure_traffic_collector" "example" {
  name                = "${var.prefix}-traffic-collector"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_monitor_diagnostic_setting" "example" {
  name                           = "${var.prefix}-diag"
  target_resource_id             = azurerm_network_function_azure_traffic_collector.example.id
  log_analytics_workspace_id     = azurerm_log_analytics_workspace.example.id
  log_analytics_destination_type = "Dedicated" # 'Dedicated' = resource specific tables

  enabled_log {
    category = "ExpressRouteCircuitIpfix"
  }
}