# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "${var.prefix}-rg-dbws"
  location = "southeastasia"
}

resource "azurerm_virtual_network" "test" {
  name                = "${var.prefix}-vnet"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/24"]
}

resource "azurerm_subnet" "test" {
  name                              = "${var.prefix}-snet"
  resource_group_name               = azurerm_resource_group.test.name
  virtual_network_name              = azurerm_virtual_network.test.name
  address_prefixes                  = ["10.0.0.0/28"]
  private_endpoint_network_policies = "Enabled"
}

resource "azurerm_databricks_workspace_serverless" "test" {
  name                          = "${var.prefix}-dbws"
  resource_group_name           = azurerm_resource_group.test.name
  location                      = azurerm_resource_group.test.location
  public_network_access_enabled = false
}

resource "azurerm_private_dns_zone" "test" {
  name                = "privatelink.azuredatabricks.net"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_zone_virtual_network_link" "test" {
  name                  = "${var.prefix}-pdzvnl"
  resource_group_name   = azurerm_resource_group.test.name
  private_dns_zone_name = azurerm_private_dns_zone.test.name
  virtual_network_id    = azurerm_virtual_network.test.id
}

resource "azurerm_private_endpoint" "test" {
  name                          = "${var.prefix}-pe"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  subnet_id                     = azurerm_subnet.test.id
  custom_network_interface_name = "${var.prefix}-nic"

  private_dns_zone_group {
    name                 = "${var.prefix}-pdzg"
    private_dns_zone_ids = [azurerm_private_dns_zone.test.id]
  }

  private_service_connection {
    name                           = "${var.prefix}-psc"
    is_manual_connection           = false
    private_connection_resource_id = azurerm_databricks_workspace_serverless.test.id
    subresource_names              = ["databricks_ui_api"]
  }
}

resource "azurerm_private_endpoint" "test2" {
  name                          = "${var.prefix}-pe2"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  subnet_id                     = azurerm_subnet.test.id
  custom_network_interface_name = "${var.prefix}-nic2"

  private_dns_zone_group {
    name                 = "${var.prefix}-pdzg2"
    private_dns_zone_ids = [azurerm_private_dns_zone.test.id]
  }

  private_service_connection {
    name                           = "${var.prefix}-psc2"
    is_manual_connection           = false
    private_connection_resource_id = azurerm_databricks_workspace_serverless.test.id
    subresource_names              = ["browser_authentication"]
  }

  depends_on = [azurerm_private_endpoint.test]
}