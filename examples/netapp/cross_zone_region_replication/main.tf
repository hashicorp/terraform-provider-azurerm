# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = var.location
}

resource "azurerm_virtual_network" "example_primary" {
  name                = "${var.prefix}-virtualnetwork-primary"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example_primary" {
  name                 = "${var.prefix}-subnet-primary"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example_primary.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "testdelegation"

    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_virtual_network" "example_cross_zone" {
  name                = "${var.prefix}-virtualnetwork-cross-zone"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.1.0.0/16"]
}

resource "azurerm_subnet" "example_cross_zone" {
  name                 = "${var.prefix}-subnet-cross-zone"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example_cross_zone.name
  address_prefixes     = ["10.1.2.0/24"]

  delegation {
    name = "testdelegation"

    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_virtual_network" "example_cross_region" {
  name                = "${var.prefix}-virtualnetwork-cross-region"
  location            = var.alt_location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.2.0.0/16"]
}

resource "azurerm_subnet" "example_cross_region" {
  name                 = "${var.prefix}-subnet-cross-region"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example_cross_region.name
  address_prefixes     = ["10.2.2.0/24"]

  delegation {
    name = "testdelegation"

    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_netapp_account" "example_primary" {
  name                = "${var.prefix}-netappaccount-primary"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_netapp_account" "example_secondary" {
  name                = "${var.prefix}-netappaccount-secondary"
  location            = var.alt_location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_netapp_pool" "example_primary" {
  name                = "${var.prefix}-netapppool-primary"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example_primary.name
  service_level       = "Premium"
  size_in_tb          = 4
}

resource "azurerm_netapp_pool" "example_cross_zone" {
  name                = "${var.prefix}-netapppool-cross-zone"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example_primary.name
  service_level       = "Premium"
  size_in_tb          = 4
}

resource "azurerm_netapp_pool" "example_cross_region" {
  name                = "${var.prefix}-netapppool-cross-region"
  location            = var.alt_location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example_secondary.name
  service_level       = "Premium"
  size_in_tb          = 4
}

# Primary volume (source) with availability zone
resource "azurerm_netapp_volume" "example_primary" {
  lifecycle {
    prevent_destroy = true
  }

  name                = "${var.prefix}-netappvolume-primary"
  location            = azurerm_resource_group.example.location
  zone                = "1"
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example_primary.name
  pool_name           = azurerm_netapp_pool.example_primary.name
  volume_path         = "${var.prefix}-netappvolume-primary"
  service_level       = "Premium"
  protocols           = ["NFSv3"]
  subnet_id           = azurerm_subnet.example_primary.id
  storage_quota_in_gb = 100

  export_policy_rule {
    rule_index      = 1
    allowed_clients = ["0.0.0.0/0"]
    protocol        = ["NFSv3"]
    unix_read_only  = false
    unix_read_write = true
  }
}

# Cross-zone replication destination volume
resource "azurerm_netapp_volume" "example_cross_zone" {
  lifecycle {
    prevent_destroy = true
  }

  depends_on = [azurerm_netapp_volume.example_primary]

  name                = "${var.prefix}-netappvolume-cross-zone"
  location            = azurerm_resource_group.example.location
  zone                = "2"
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example_primary.name
  pool_name           = azurerm_netapp_pool.example_cross_zone.name
  volume_path         = "${var.prefix}-netappvolume-cross-zone"
  service_level       = "Premium"
  protocols           = ["NFSv3"]
  subnet_id           = azurerm_subnet.example_cross_zone.id
  storage_quota_in_gb = 100

  export_policy_rule {
    rule_index      = 1
    allowed_clients = ["0.0.0.0/0"]
    protocol        = ["NFSv3"]
    unix_read_only  = false
    unix_read_write = true
  }

  data_protection_replication {
    endpoint_type             = "dst"
    remote_volume_location    = azurerm_resource_group.example.location
    remote_volume_resource_id = azurerm_netapp_volume.example_primary.id
    replication_frequency     = "10minutes"
  }
}

# Cross-region replication destination volume
resource "azurerm_netapp_volume" "example_cross_region" {
  lifecycle {
    prevent_destroy = true
  }

  depends_on = [
    azurerm_netapp_volume.example_primary,
    azurerm_netapp_volume.example_cross_zone
  ]

  name                = "${var.prefix}-netappvolume-cross-region"
  location            = var.alt_location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example_secondary.name
  pool_name           = azurerm_netapp_pool.example_cross_region.name
  volume_path         = "${var.prefix}-netappvolume-cross-region"
  service_level       = "Premium"
  protocols           = ["NFSv3"]
  subnet_id           = azurerm_subnet.example_cross_region.id
  storage_quota_in_gb = 100

  export_policy_rule {
    rule_index      = 1
    allowed_clients = ["0.0.0.0/0"]
    protocol        = ["NFSv3"]
    unix_read_only  = false
    unix_read_write = true
  }

  data_protection_replication {
    endpoint_type             = "dst"
    remote_volume_location    = azurerm_resource_group.example.location
    remote_volume_resource_id = azurerm_netapp_volume.example_primary.id
    replication_frequency     = "10minutes"
  }
}
