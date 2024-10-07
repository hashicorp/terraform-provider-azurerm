# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {}
}

module "naming" {
  source  = "Azure/naming/azurerm"
  version = "0.4.0"
  prefix  = [var.prefix]
}

resource "azurerm_resource_group" "example" {
  name     = module.naming.resource_group.name
  location = "westeurope"
}

resource "azurerm_iothub_dps" "example" {
  name                          = module.naming.iothub_dps.name
  resource_group_name           = azurerm_resource_group.example.name
  location                      = azurerm_resource_group.example.location
  allocation_policy             = "Hashed"
  data_residency_enabled        = false
  public_network_access_enabled = true

  linked_hub {
    connection_string       = var.iot_hub_connection_string
    location                = azurerm_resource_group.example.location
    allocation_weight       = 1
    apply_allocation_policy = true
  }

  sku {
    name     = "S1"
    capacity = 1
  }

  tags = {
    environment = "Development"
    region      = azurerm_resource_group.example.location
  }
}
