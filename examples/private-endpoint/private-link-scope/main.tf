# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {}
}

locals {
  private_dns_zones_names = toset([
    "privatelink.agentsvc.azure-automation.net",
    "privatelink.blob.core.windows.net",
    "privatelink.monitor.azure.com",
    "privatelink.ods.opinsights.azure.com",
    "privatelink.oms.opinsights.azure.com",
  ])
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location
}

resource "azurerm_virtual_network" "example" {
  name                = "${var.prefix}-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "${var.prefix}-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]

  private_link_service_network_policies_enabled = false
}

resource "azurerm_private_dns_zone" "example" {
  for_each = local.private_dns_zones_names

  name                = each.value
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_monitor_private_link_scope" "example" {
  name                = "${var.prefix}-ampls"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_private_endpoint" "this" {
  name                = "${var.prefix}-ape"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  subnet_id           = azurerm_subnet.example.id

  private_dns_zone_group {
    name                 = "default"
    private_dns_zone_ids = [for _, v in azurerm_private_dns_zone.example : v.id]
  }

  private_service_connection {
    name                           = "${var.prefix}-psc"
    is_manual_connection           = false
    private_connection_resource_id = azurerm_monitor_private_link_scope.example.id
    subresource_names              = ["azuremonitor"]
  }
}
