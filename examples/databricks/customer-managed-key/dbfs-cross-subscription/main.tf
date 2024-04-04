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
  name     = "${var.prefix}-databricks-cmk"
  location = "West Europe"
}

resource "azurerm_resource_group" "keyVault" {
  provider = azurerm.keyVaultSubscription

  name     = "${var.prefix}-databricks-cmk"
  location = "West Europe"
}

resource "azurerm_databricks_workspace" "example" {
  name                        = "${var.prefix}-DBW"
  resource_group_name         = azurerm_resource_group.example.name
  location                    = azurerm_resource_group.example.location
  sku                         = "premium"
  managed_resource_group_name = "${var.prefix}-DBW-managed-dbfs"

  customer_managed_key_enabled = true

  tags = {
    Environment = "Sandbox"
  }
}

resource "azurerm_databricks_workspace_root_dbfs_customer_managed_key" "example" {
  depends_on = [azurerm_key_vault_access_policy.databricks]

  workspace_id     = azurerm_databricks_workspace.example.id
  key_vault_id     = azurerm_key_vault.example.id
  key_vault_key_id = azurerm_key_vault_key.dbfs.id
}

resource "azurerm_key_vault" "example" {
  provider = azurerm.keyVaultSubscription
  
  name                = "${var.prefix}-keyvault"
  location            = azurerm_resource_group.keyVault.location
  resource_group_name = azurerm_resource_group.keyVault.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"

  purge_protection_enabled   = true
  soft_delete_retention_days = 7
}

resource "azurerm_key_vault_key" "dbfs" {
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

resource "azurerm_key_vault_access_policy" "databricks" {
  depends_on = [azurerm_databricks_workspace.example]

  provider = azurerm.keyVaultSubscription

  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = azurerm_databricks_workspace.example.storage_account_identity.0.tenant_id
  object_id    = azurerm_databricks_workspace.example.storage_account_identity.0.principal_id

  key_permissions = [
    "Get",
    "UnwrapKey",
    "WrapKey",
  ]
}
