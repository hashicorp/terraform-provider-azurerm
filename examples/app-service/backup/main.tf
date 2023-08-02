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
  name                     = "${var.prefix}sa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name                  = "${var.prefix}-sc"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

resource "azurerm_service_plan" "example" {
  name                = "${var.prefix}-sp-zipdeploy"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  os_type             = "Windows"
  sku_name            = "S1"
}

data "azurerm_storage_account_sas" "example" {
  connection_string = azurerm_storage_account.example.primary_connection_string
  https_only        = true

  resource_types {
    service   = false
    container = false
    object    = true
  }

  services {
    blob  = true
    queue = false
    table = false
    file  = false
  }

  start  = "2019-03-21"
  expiry = "2020-03-21"

  permissions {
    read    = false
    write   = true
    delete  = false
    list    = false
    add     = false
    create  = false
    update  = false
    process = false
    tag     = false
    filter  = false
  }
}

resource "azurerm_windows_web_app" "example" {
  name                = "${var.prefix}-appservice"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  service_plan_id     = azurerm_service_plan.example.id

  backup {
    name                = "Example"
    storage_account_url = "https://${azurerm_storage_account.example.name}.blob.core.windows.net/${azurerm_storage_container.example.name}${data.azurerm_storage_account_sas.example.sas}&sr=b"
    schedule {
      frequency_interval = 30
      frequency_unit     = "Day"
    }
  }

  site_config {
    application_stack {
      dotnet_version = "v6.0"
      current_stack  = "dotnet"
    }
  }
}