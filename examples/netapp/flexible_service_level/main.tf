# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location

  tags = {
    "SkipNRMSNSG" = true
  }
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

# NetApp Pool with Flexible Service Level
resource "azurerm_netapp_pool" "example" {
  name                    = "${var.prefix}-netapppool"
  location                = azurerm_resource_group.example.location
  resource_group_name     = azurerm_resource_group.example.name
  account_name            = azurerm_netapp_account.example.name
  service_level           = "Flexible"
  size_in_tb              = 4
  qos_type                = "Manual"
  custom_throughput_mibps = 400
}

# High-throughput volume example
resource "azurerm_netapp_volume" "high_throughput" {
  lifecycle {
    prevent_destroy = true
  }

  name                = "${var.prefix}-netappvolume-highthroughput"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  pool_name           = azurerm_netapp_pool.example.name
  volume_path         = "high-throughput-volume"
  service_level       = "Flexible"
  subnet_id           = azurerm_subnet.example.id
  protocols           = ["NFSv4.1"]
  storage_quota_in_gb = 500
  throughput_in_mibps = 200

  export_policy_rule {
    rule_index      = 1
    allowed_clients = ["0.0.0.0/0"]
    protocol        = ["NFSv4.1"]
    unix_read_only  = false
    unix_read_write = true
  }

  tags = {
    "Environment" = "Example"
    "Purpose"     = "High Throughput Volume"
  }
}

# Low-throughput, high-capacity volume example
resource "azurerm_netapp_volume" "high_capacity" {
  lifecycle {
    prevent_destroy = true
  }

  name                = "${var.prefix}-netappvolume-highcapacity"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  pool_name           = azurerm_netapp_pool.example.name
  volume_path         = "high-capacity-volume"
  service_level       = "Flexible"
  subnet_id           = azurerm_subnet.example.id
  protocols           = ["NFSv4.1"]
  storage_quota_in_gb = 2000
  throughput_in_mibps = 64

  export_policy_rule {
    rule_index      = 1
    allowed_clients = ["10.0.0.0/16"]
    protocol        = ["NFSv4.1"]
    unix_read_only  = false
    unix_read_write = true
  }

  tags = {
    "Environment" = "Example"
    "Purpose"     = "High Capacity Volume"
  }
}
