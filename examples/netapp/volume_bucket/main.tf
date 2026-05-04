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

# The first bucket created on a set of volumes sharing the same backing IP
# must supply a server FQDN and a server certificate (base64-encoded PEM
# containing both certificate and private key). This example generates a
# self-signed certificate purely for demonstration - replace it with a
# CA-signed certificate in production.
resource "tls_private_key" "bucket" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "bucket" {
  private_key_pem = tls_private_key.bucket.private_key_pem

  subject {
    common_name = var.server_fqdn
  }

  dns_names = [var.server_fqdn]

  validity_period_hours = 8760 # 1 year

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
  ]
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

  server {
    fqdn = var.server_fqdn

    # certificate_pem must be the base64-encoded concatenation of the PEM
    # certificate and the PEM private key.
    certificate_pem = base64encode("${tls_self_signed_cert.bucket.cert_pem}${tls_private_key.bucket.private_key_pem}")
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
