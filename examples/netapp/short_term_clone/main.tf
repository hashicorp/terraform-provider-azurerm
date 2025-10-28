# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location
}

resource "azurerm_virtual_network" "example" {
  name                = "${var.prefix}-virtualnetwork"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "${var.prefix}-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "testdelegation"

    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_netapp_account" "example" {
  name                = "${var.prefix}-netappaccount"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_netapp_pool" "example" {
  name                = "${var.prefix}-pool"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  service_level       = "Standard"
  size_in_tb          = 4
}

resource "azurerm_netapp_volume" "source" {
  name                = "${var.prefix}-source-volume"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  pool_name           = azurerm_netapp_pool.example.name
  volume_path         = "${var.prefix}-source-volume"
  service_level       = "Standard"
  protocols           = ["NFSv3"]
  subnet_id           = azurerm_subnet.example.id
  storage_quota_in_gb = 100

  export_policy_rule {
    rule_index        = 1
    allowed_clients   = ["0.0.0.0/0"]
    protocols_enabled = ["NFSv3"]
    unix_read_write   = true
  }
}

resource "azurerm_netapp_snapshot" "example" {
  name                = "${var.prefix}-snapshot"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  pool_name           = azurerm_netapp_pool.example.name
  volume_name         = azurerm_netapp_volume.source.name
}

resource "azurerm_netapp_volume" "short_term_clone" {
  name                                                 = "${var.prefix}-clone-volume"
  location                                             = azurerm_resource_group.example.location
  resource_group_name                                  = azurerm_resource_group.example.name
  account_name                                         = azurerm_netapp_account.example.name
  pool_name                                            = azurerm_netapp_pool.example.name
  volume_path                                          = "${var.prefix}-clone-volume"
  service_level                                        = "Standard"
  protocols                                            = ["NFSv3"]
  subnet_id                                            = azurerm_subnet.example.id
  storage_quota_in_gb                                  = 100
  create_from_snapshot_resource_id                     = azurerm_netapp_snapshot.example.id
  accept_grow_capacity_pool_for_short_term_clone_split = "Accepted"

  export_policy_rule {
    rule_index        = 1
    allowed_clients   = ["0.0.0.0/0"]
    protocols_enabled = ["NFSv3"]
    unix_read_write   = true
  }
}
