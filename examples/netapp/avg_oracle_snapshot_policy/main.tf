# Copyright (c) HashiCorp, Inc.
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
  name                = "${var.prefix}-vnet"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.6.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "${var.prefix}-delegated-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.6.2.0/24"]

  delegation {
    name = "exampledelegation"

    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_netapp_account" "example" {
  name                = "${var.prefix}-netapp-account"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  depends_on = [
    azurerm_subnet.example
  ]
}

resource "azurerm_netapp_pool" "example" {
  name                = "${var.prefix}-netapp-pool"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  service_level       = "Standard"
  size_in_tb          = 4
  qos_type            = "Manual"
}

resource "azurerm_netapp_snapshot_policy" "example" {
  name                = "${var.prefix}-snapshot-policy"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  enabled             = true

  monthly_schedule {
    snapshots_to_keep = 1
    days_of_month     = [15, 30]
    hour              = 23
    minute            = 30
  }

  tags = {
    "CreatedOnDate"    = "2022-07-08T23:50:21Z",
    "SkipASMAzSecPack" = "true"
  }
}

resource "azurerm_netapp_volume_group_oracle" "example" {
  name                   = "${var.prefix}-volume-group-oracle"
  location               = azurerm_resource_group.example.location
  resource_group_name    = azurerm_resource_group.example.name
  account_name           = azurerm_netapp_account.example.name
  group_description      = "Example volume group for Oracle"
  application_identifier = "TST"

  volume {
    name                       = "${var.prefix}-volume-ora1"
    volume_path                = "${var.prefix}-my-unique-file-ora-path-1"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.example.id
    subnet_id                  = azurerm_subnet.example.id
    zone                       = "1"
    volume_spec_name           = "ora-data1"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv4.1"]
    security_style             = "unix"
    snapshot_directory_visible = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }
  }

  volume {
    name                       = "${var.prefix}-volume-oraLog"
    volume_path                = "${var.prefix}-my-unique-file-oralog-path"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.example.id
    subnet_id                  = azurerm_subnet.example.id
    zone                       = "1"
    volume_spec_name           = "ora-log"
    storage_quota_in_gb        = 1024
    throughput_in_mibps        = 24
    protocols                  = ["NFSv4.1"]
    security_style             = "unix"
    snapshot_directory_visible = false

    export_policy_rule {
      rule_index          = 1
      allowed_clients     = "0.0.0.0/0"
      nfsv3_enabled       = false
      nfsv41_enabled      = true
      unix_read_only      = false
      unix_read_write     = true
      root_access_enabled = false
    }
  }
}
