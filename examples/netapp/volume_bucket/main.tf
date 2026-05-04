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
  name                = "${var.prefix}-vnet"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.47.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "${var.prefix}-delegated-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.47.2.0/24"]

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
}

resource "azurerm_netapp_pool" "example" {
  name                = "${var.prefix}-netapp-pool"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  service_level       = "Standard"
  size_in_tb          = 4
  qos_type            = "Auto"
}

resource "azurerm_netapp_volume" "example" {
  name                = "${var.prefix}-netapp-volume"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  pool_name           = azurerm_netapp_pool.example.name
  volume_path         = "${var.prefix}-bucket-vol"
  service_level       = "Standard"
  subnet_id           = azurerm_subnet.example.id
  storage_quota_in_gb = 100
  protocols           = ["NFSv3"]
}

resource "azurerm_netapp_volume_bucket" "example" {
  name      = "${var.prefix}-bucket"
  volume_id = azurerm_netapp_volume.example.id
  path      = "/"

  file_system_user {
    nfs_user {
      group_id = 1000
      user_id  = 1000
    }
  }
}

# Generates the S3 access key / secret key pair used to authenticate against
# the bucket. Because `store_in_key_vault` is omitted (default `false`), the
# Azure NetApp Files API returns the keys inline and the provider stores them
# as the sensitive `access_key` and `secret_key` attributes on this resource.
# Consume them from `outputs.tf` (e.g. `terraform output -raw bucket_access_key`)
# or by referencing them in another resource / provider.
#
# NOTE: this puts the bucket access / secret keys into Terraform state. For
# production workloads use the Key Vault-backed example
# (`examples/netapp/volume_bucket_akv`) which keeps the credentials in Azure
# Key Vault instead.
#
# `key_pair_expiry_days` is `ForceNew` - changing it (or tainting this
# resource) generates a new key pair and immediately invalidates the
# previous one.
resource "azurerm_netapp_volume_bucket_credentials" "example" {
  bucket_id            = azurerm_netapp_volume_bucket.example.id
  key_pair_expiry_days = 30
}
