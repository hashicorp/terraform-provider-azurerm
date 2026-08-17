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
  address_space       = ["10.48.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "${var.prefix}-delegated-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.48.2.0/24"]

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

# Self-signed bucket server certificate. Replace this with a CA-signed
# certificate in production. The Subject Alternative Name must include
# `server_fqdn`.
resource "tls_private_key" "bucket" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "bucket" {
  private_key_pem = tls_private_key.bucket.private_key_pem

  subject {
    common_name = var.server_fqdn
  }

  dns_names             = [var.server_fqdn]
  validity_period_hours = 8760

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
  ]
}

# The FIRST bucket on the volume establishes the bucket server (FQDN +
# certificate). It must be created with the
# `azurerm_netapp_volume_bucket_with_server` resource.
resource "azurerm_netapp_volume_bucket_with_server" "first" {
  name      = "${var.prefix}-bucket-first"
  volume_id = azurerm_netapp_volume.example.id

  file_system_nfs_user {
    group_id = 1000
    user_id  = 1000
  }

  server {
    fqdn            = var.server_fqdn
    certificate_pem = base64encode("${tls_self_signed_cert.bucket.cert_pem}${tls_private_key.bucket.private_key_pem}")
  }
}

# SUBSEQUENT buckets reuse the server configuration established by the first
# bucket, so they are created with the server-less
# `azurerm_netapp_volume_bucket` resource. Declaring a `server` block on more
# than one bucket would overwrite the shared server configuration.
resource "azurerm_netapp_volume_bucket" "second" {
  name      = "${var.prefix}-bucket-second"
  volume_id = azurerm_netapp_volume.example.id

  file_system_nfs_user {
    group_id = 2000
    user_id  = 2000
  }

  depends_on = [
    azurerm_netapp_volume_bucket_with_server.first,
  ]
}
