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

# Primary region networking
resource "azurerm_virtual_network" "example_primary" {
  name                = "${var.prefix}-vnet-primary"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.47.0.0/16"]
}

resource "azurerm_subnet" "example_primary" {
  name                 = "${var.prefix}-delegated-subnet-primary"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example_primary.name
  address_prefixes     = ["10.47.2.0/24"]

  delegation {
    name = "exampledelegation"

    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

# Secondary region networking
resource "azurerm_virtual_network" "example_secondary" {
  name                = "${var.prefix}-vnet-secondary"
  location            = var.alt_location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.48.0.0/16"]
}

resource "azurerm_subnet" "example_secondary" {
  name                 = "${var.prefix}-delegated-subnet-secondary"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example_secondary.name
  address_prefixes     = ["10.48.2.0/24"]

  delegation {
    name = "exampledelegation"

    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

# Primary region NetApp infrastructure
resource "azurerm_netapp_account" "example_primary" {
  name                = "${var.prefix}-netapp-account-primary"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  depends_on = [
    azurerm_subnet.example_primary
  ]
}

resource "azurerm_netapp_pool" "example_primary" {
  name                = "${var.prefix}-netapp-pool-primary"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example_primary.name
  service_level       = "Standard"
  size_in_tb          = 4
  qos_type            = "Manual"
}

# Secondary region NetApp infrastructure
resource "azurerm_netapp_account" "example_secondary" {
  name                = "${var.prefix}-netapp-account-secondary"
  location            = var.alt_location
  resource_group_name = azurerm_resource_group.example.name

  depends_on = [
    azurerm_subnet.example_secondary
  ]
}

resource "azurerm_netapp_pool" "example_secondary" {
  name                = "${var.prefix}-netapp-pool-secondary"
  location            = var.alt_location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example_secondary.name
  service_level       = "Standard"
  size_in_tb          = 4
  qos_type            = "Manual"
}

# Primary Oracle volume group
resource "azurerm_netapp_volume_group_oracle" "example_primary" {
  name                   = "${var.prefix}-NetAppVolumeGroupOracle-primary"
  location               = azurerm_resource_group.example.location
  resource_group_name    = azurerm_resource_group.example.name
  account_name           = azurerm_netapp_account.example_primary.name
  group_description      = "Primary Oracle volume group for CRR"
  application_identifier = "TST"

  volume {
    name                       = "${var.prefix}-volume-ora1-primary"
    volume_path                = "${var.prefix}-my-unique-file-ora-path-1-primary"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.example_primary.id
    subnet_id                  = azurerm_subnet.example_primary.id
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
    name                       = "${var.prefix}-volume-oraLog-primary"
    volume_path                = "${var.prefix}-my-unique-file-oralog-path-primary"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.example_primary.id
    subnet_id                  = azurerm_subnet.example_primary.id
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

# Secondary Oracle volume group with CRR
resource "azurerm_netapp_volume_group_oracle" "example_secondary" {
  depends_on = [azurerm_netapp_volume_group_oracle.example_primary]

  name                   = "${var.prefix}-NetAppVolumeGroupOracle-secondary"
  location               = var.alt_location
  resource_group_name    = azurerm_resource_group.example.name
  account_name           = azurerm_netapp_account.example_secondary.name
  group_description      = "Secondary Oracle volume group for CRR"
  application_identifier = "TST"

  volume {
    name                       = "${var.prefix}-volume-ora1-secondary"
    volume_path                = "${var.prefix}-my-unique-file-ora-path-1-secondary"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.example_secondary.id
    subnet_id                  = azurerm_subnet.example_secondary.id
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

    data_protection_replication {
      endpoint_type             = "dst"
      remote_volume_location    = azurerm_resource_group.example.location
      remote_volume_resource_id = azurerm_netapp_volume_group_oracle.example_primary.volume[0].id
      replication_frequency     = "10minutes"
    }
  }

  volume {
    name                       = "${var.prefix}-volume-oraLog-secondary"
    volume_path                = "${var.prefix}-my-unique-file-oralog-path-secondary"
    service_level              = "Standard"
    capacity_pool_id           = azurerm_netapp_pool.example_secondary.id
    subnet_id                  = azurerm_subnet.example_secondary.id
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

    data_protection_replication {
      endpoint_type             = "dst"
      remote_volume_location    = azurerm_resource_group.example.location
      remote_volume_resource_id = azurerm_netapp_volume_group_oracle.example_primary.volume[1].id
      replication_frequency     = "10minutes"
    }
  }
}
