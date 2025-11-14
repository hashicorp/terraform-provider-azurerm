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

resource "azurerm_storage_account" "example" {
  name                     = "${var.prefix}sa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name                 = "${var.prefix}csc"
  storage_account_id   = azurerm_storage_account.example.id
  container_access_type = "private"
}

resource "azurerm_cognitive_account_connection" "example" {
  name                 = "${var.prefix}-connection"
  cognitive_account_id = azurerm_cognitive_account.example.id
  auth_type            = "OAuth2"
  category             = "AzureBlob"
  target               = azurerm_storage_account.example.primary_blob_endpoint

  metadata = {
    containerName = azurerm_storage_container.example.name
    accountName   = azurerm_storage_account.example.name
  }

  oauth2 {
    auth_url        = "https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000/oauth2/v2.0/token"
    client_id       = "00000000-0000-0000-0000-000000000000"
    client_secret   = "placeHolderClientSecret"
    tenant_id       = "00000000-0000-0000-0000-000000000000"
    developer_token = "placeHolderDevToken"
    refresh_token   = "placeRefreshToken"
    username        = "placeHolderUsername"
    password        = "placeHolderPassword"
  }
}
