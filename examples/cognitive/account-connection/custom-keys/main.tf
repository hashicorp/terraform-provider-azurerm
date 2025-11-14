# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location
}

resource "azurerm_cognitive_account" "example" {
  name                       = "${var.prefix}-aiservices"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  kind                       = "AIServices"
  sku_name                   = "S0"
  project_management_enabled = true
  custom_subdomain_name      = "${var.prefix}aiservices"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_cognitive_account" "openai" {
  name                = "${var.prefix}-openai"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  kind                = "OpenAI"
  sku_name            = "S0"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_cognitive_account_connection" "example" {
  name                 = "${var.prefix}-connection"
  cognitive_account_id = azurerm_cognitive_account.example.id
  auth_type            = "CustomKeys"
  category             = "CustomKeys"
  target               = azurerm_cognitive_account.openai.endpoint

  metadata = {
    apiType    = "Azure"
    resourceId = azurerm_cognitive_account.openai.id
    location   = azurerm_cognitive_account.openai.location
  }

  custom_keys = {
    primaryKey   = azurerm_cognitive_account.openai.primary_access_key
    secondaryKey = azurerm_cognitive_account.openai.secondary_access_key
  }
}
