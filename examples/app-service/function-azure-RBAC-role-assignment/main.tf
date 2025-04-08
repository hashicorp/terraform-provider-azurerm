# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location
}

resource "azurerm_storage_account" "example" {
  name                     = "${var.prefix}storageacct"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "${var.prefix}-law"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_application_insights" "example" {
  name                = "${var.prefix}-appinsights"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  application_type    = "web"
  workspace_id        = azurerm_log_analytics_workspace.example.id
}

resource "azurerm_service_plan" "example" {
  name                = "${var.prefix}-sp"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  os_type             = "Linux"
  sku_name            = "Y1"
}

resource "azurerm_linux_function_app" "example" {
  name                = "${var.prefix}-LFA"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  service_plan_id     = azurerm_service_plan.example.id

  storage_account_name       = azurerm_storage_account.example.name
  storage_account_access_key = azurerm_storage_account.example.primary_access_key

  site_config {
    application_insights_connection_string = azurerm_application_insights.example.connection_string
  }

  identity {
    type = "SystemAssigned"
  }
}

data "azurerm_subscription" "primary" {
}

resource "azurerm_role_assignment" "example" {
  scope                = data.azurerm_subscription.primary.id
  role_definition_name = var.role_definition_name
  principal_id         = azurerm_linux_function_app.example.identity.0.principal_id
}
