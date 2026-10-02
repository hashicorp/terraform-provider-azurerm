# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {
    netapp {
      prevent_volume_destruction = true
    }
  }
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location
  tags = {
    "SkipNRMSNSG" = "true"
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
    name = "netapp"

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
  name                = "${var.prefix}-netapppool"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  service_level       = "Standard"
  size_in_tb          = var.pool_size_in_tb
  qos_type            = "Manual"
}

# Large volume running in Breakthrough Mode. Breakthrough Mode places the volume on dedicated
# capacity, which delivers higher throughput and allows the volume to grow up to 2,400 TiB.
resource "azurerm_netapp_volume" "example" {
  lifecycle {
    prevent_destroy = true
  }

  name                = "${var.prefix}-netappvolume"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  pool_name           = azurerm_netapp_pool.example.name
  volume_path         = "${var.prefix}-netappvolume"
  service_level       = "Standard"
  subnet_id           = azurerm_subnet.example.id
  protocols           = ["NFSv4.1"]
  security_style      = "unix"
  throughput_in_mibps = var.throughput_in_mibps

  # Breakthrough Mode is only supported on large volumes and cannot be changed after creation
  large_volume_enabled      = true
  breakthrough_mode_enabled = true

  # Breakthrough Mode volumes are sized between 2,400 GiB and 2,400 TiB (2,457,600 GiB)
  storage_quota_in_gb = var.storage_quota_in_gb

  export_policy_rule {
    rule_index          = 1
    allowed_clients     = ["10.0.0.0/16"]
    protocol            = ["NFSv4.1"]
    unix_read_only      = false
    unix_read_write     = true
    root_access_enabled = true
  }
}
