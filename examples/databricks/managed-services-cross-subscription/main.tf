# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  features {}
}

provider "azurerm" {
  features {}
  alias           = "keyVaultSubscription"
  subscription_id = "00000000-0000-0000-0000-000000000000" # Subscription where the Key Vault should be hosted
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-databricks-managed-services"
  location = "West Europe"
}

resource "azurerm_resource_group" "keyVault" {
  provider = azurerm.keyVaultSubscription

  name     = "${var.prefix}-databricks-managed-services"
  location = "West Europe"
}

resource "azurerm_databricks_workspace" "example" {
  depends_on = [azurerm_key_vault_access_policy.managed]

  name                        = "${var.prefix}-DBW"
  resource_group_name         = azurerm_resource_group.example.name
  location                    = azurerm_resource_group.example.location
  sku                         = "premium"
  managed_resource_group_name = "${var.prefix}-DBW-managed-services"

  managed_services_cmk_key_vault_id     = azurerm_key_vault.example.id
  managed_services_cmk_key_vault_key_id = azurerm_key_vault_key.services.id

  managed_disk_cmk_key_vault_id     = azurerm_key_vault.example.id
  managed_disk_cmk_key_vault_key_id = azurerm_key_vault_key.disk.id

  tags = {
    Environment = "Sandbox"
  }
}

resource "azurerm_key_vault" "example" {
  provider = azurerm.keyVaultSubscription

  name                = "${var.prefix}-keyvault"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"

  soft_delete_retention_days = 7
}

resource "azurerm_key_vault_key" "services" {
  depends_on = [azurerm_key_vault_access_policy.terraform]

  provider = azurerm.keyVaultSubscription

  name         = "${var.prefix}-certificate"
  key_vault_id = azurerm_key_vault.example.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_key_vault_key" "disk" {
  depends_on = [azurerm_key_vault_access_policy.terraform]

  provider = azurerm.keyVaultSubscription

  name         = "${var.prefix}-certificate"
  key_vault_id = azurerm_key_vault.example.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_key_vault_access_policy" "terraform" {
  provider = azurerm.keyVaultSubscription

  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = azurerm_key_vault.example.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Get",
    "List",
    "Create",
    "Decrypt",
    "Encrypt",
    "Sign",
    "UnwrapKey",
    "Verify",
    "WrapKey",
    "Delete",
    "Restore",
    "Recover",
    "Update",
    "Purge",
    "GetRotationPolicy",
    "SetRotationPolicy",
  ]
}

resource "azurerm_key_vault_access_policy" "managed" {
  provider = azurerm.keyVaultSubscription

  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = azurerm_key_vault.example.tenant_id
  object_id    = "00000000-0000-0000-0000-000000000000" # See the README.md file for instructions on how to lookup the correct value to enter here.

  key_permissions = [
    "Get",
    "UnwrapKey",
    "WrapKey",
  ]
}
