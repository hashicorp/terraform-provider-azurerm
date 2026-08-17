# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {
    netapp {
      prevent_volume_destruction = true
    }
    key_vault {
      purge_soft_delete_on_destroy    = true
      recover_soft_deleted_key_vaults = true
    }
  }
}

data "azurerm_client_config" "current" {}

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

# A system-assigned managed identity is enabled on the NetApp account so it
# can read the bucket server certificate from Key Vault and write the
# generated bucket credentials to Key Vault.
resource "azurerm_netapp_account" "example" {
  name                = "${var.prefix}-netapp-account"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  identity {
    type = "SystemAssigned"
  }
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

# Two separate Key Vaults are recommended by the Azure NetApp Files
# documentation: one read-mostly vault for the server certificate, and one
# write vault for the bucket credentials.

resource "azurerm_key_vault" "certificate" {
  name                       = "${var.prefix}-cert-kv"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  rbac_authorization_enabled = false
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7
  purge_protection_enabled   = false
}

resource "azurerm_key_vault" "credentials" {
  name                       = "${var.prefix}-creds-kv"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  rbac_authorization_enabled = false
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7
  purge_protection_enabled   = false
}

# Allow the principal running Terraform to import the bucket server certificate.
resource "azurerm_key_vault_access_policy" "deployer_certificate" {
  key_vault_id = azurerm_key_vault.certificate.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  certificate_permissions = [
    "Get",
    "List",
    "Create",
    "Import",
    "Update",
    "Delete",
    "Purge",
    "Recover",
  ]

  secret_permissions = [
    "Get",
    "List",
    "Set",
    "Delete",
    "Purge",
    "Recover",
  ]
}

# Permissions required by the NetApp account's system-assigned managed
# identity on the certificate vault (per the Object REST API documentation).
resource "azurerm_key_vault_access_policy" "anf_certificate" {
  key_vault_id = azurerm_key_vault.certificate.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_netapp_account.example.identity[0].principal_id

  certificate_permissions = [
    "Get",
    "List",
    "Update",
    "Create",
    "Import",
    "ManageContacts",
    "GetIssuers",
    "ListIssuers",
    "SetIssuers",
    "DeleteIssuers",
  ]

  secret_permissions = [
    "Get",
    "List",
    "Set",
    "Delete",
  ]
}

# Permissions required by the NetApp account's system-assigned managed
# identity on the credentials vault.
resource "azurerm_key_vault_access_policy" "anf_credentials" {
  key_vault_id = azurerm_key_vault.credentials.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_netapp_account.example.identity[0].principal_id

  secret_permissions = [
    "Get",
    "List",
    "Set",
    "Delete",
  ]
}

# Self-signed bucket server certificate. Replace this with an
# `azurerm_key_vault_certificate` block that imports a CA-signed certificate
# in production. The Subject Alternative Name must include `server_fqdn`.
resource "azurerm_key_vault_certificate" "bucket" {
  name         = "${var.prefix}-bucket-cert"
  key_vault_id = azurerm_key_vault.certificate.id

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = false
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      key_usage = [
        "digitalSignature",
        "keyEncipherment",
      ]

      extended_key_usage = [
        "1.3.6.1.5.5.7.3.1", # Server Authentication
      ]

      subject = "CN=${var.server_fqdn}"

      subject_alternative_names {
        dns_names = [var.server_fqdn]
      }

      validity_in_months = 12
    }
  }

  depends_on = [
    azurerm_key_vault_access_policy.deployer_certificate,
  ]
}

resource "azurerm_netapp_volume_bucket_with_server" "example" {
  name      = "${var.prefix}-bucket-kv"
  volume_id = azurerm_netapp_volume.example.id

  file_system_nfs_user {
    group_id = 1000
    user_id  = 1000
  }

  server {
    fqdn = var.server_fqdn
  }

  key_vault {
    certificate_key_vault_uri = azurerm_key_vault.certificate.vault_uri
    certificate_name          = azurerm_key_vault_certificate.bucket.name
    credentials_key_vault_uri = azurerm_key_vault.credentials.vault_uri
    credentials_secret_name   = "${var.prefix}-bucket-creds"
  }

  depends_on = [
    azurerm_key_vault_access_policy.anf_certificate,
    azurerm_key_vault_access_policy.anf_credentials,
  ]
}

# Generates the S3 access key / secret key pair used to authenticate against
# the bucket. The Azure NetApp Files API writes the keys as a JSON secret
# (`{"access_key_id": "...", "secret_access_key": "..."}`) into the
# credentials Key Vault configured on the parent bucket
# (`key_vault.0.credentials_key_vault_uri` / `credentials_secret_name`).
# Consumers fetch the keys from Key Vault (e.g. via `azurerm_key_vault_secret`).
#
# The action is invoked once on the lifecycle of the `terraform_data`
# trigger below; to rotate credentials, taint the trigger (or change its
# `input`) and re-apply.
action "azurerm_netapp_volume_bucket_credentials" "example" {
  config {
    bucket_id            = azurerm_netapp_volume_bucket_with_server.example.id
    key_pair_expiry_days = 30
  }
}

resource "terraform_data" "bucket_credentials" {
  input = azurerm_netapp_volume_bucket_with_server.example.id

  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.azurerm_netapp_volume_bucket_credentials.example]
    }
  }
}
