---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_account_customer_managed_key"
description: |-
  Manages the customer managed key of an Azure Storage Account.
---

# azurerm_storage_account_customer_managed_key

Manages the customer managed key of an Azure Storage Account.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

provider "azurerm" {
  alias = "keyVault"

  features {
    key_vault {
      purge_soft_delete_on_destroy = false
    }
  }
}

resource "azurerm_resource_group" "tfex" {
  name     = "tfex-RG"
  location = "westeurope"
}

resource "azurerm_key_vault" "tfex" {
  name                = "tfex-key-vault"
  provider            = azurerm.keyVault
  location            = azurerm_resource_group.tfex.location
  resource_group_name = azurerm_resource_group.tfex.name
  tenant_id           = data.azurerm_client_config.current.tenant_id

  soft_delete_enabled      = true
  purge_protection_enabled = true

  sku_name = "standard"

  tags = {
    environment = "testing"
  }
}

resource "azurerm_key_vault_key" "tfex" {
  name                       = "tfex-key"
  key_vault_id               = azurerm_key_vault.tfex.id
  key_vault_access_policy_id = azurerm_key_vault_access_policy.storage.id
  key_type                   = "RSA"
  key_size                   = 2048
  key_opts                   = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]
}

resource "azurerm_key_vault_access_policy" "storage" {
  key_vault_id = azurerm_key_vault.tfex.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_storage_account.tfex.identity.0.principal_id

  key_permissions    = ["get", "create", "delete", "list", "restore", "recover", "unwrapkey", "wrapkey", "purge", "encrypt", "decrypt", "sign", "verify"]
  secret_permissions = ["get"]
}

resource "azurerm_key_vault_access_policy" "client" {
  key_vault_id = azurerm_key_vault.tfex.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_storage_account.tfex.identity.0.principal_id

  key_permissions    = ["get", "create", "delete", "list", "restore", "recover", "unwrapkey", "wrapkey", "purge", "encrypt", "decrypt", "sign", "verify"]
  secret_permissions = ["get"]
}

resource "azurerm_storage_account" "tfex" {
  name                     = "tfexstorageaccount"
  resource_group_name      = azurerm_resource_group.tfex.name
  location                 = azurerm_resource_group.tfex.location
  account_tier             = "Standard"
  account_replication_type = "GRS"

  identity {
    type = "SystemAssigned"
  }

  tags = {
    environment = "testing"
  }
}

resource "azurerm_storage_account_customer_managed_key" "tfex" {
  storage_account_id = azurerm_storage_account.tfex.id
  key_vault_id       = azurerm_key_vault.tfex.id
  key_name           = azurerm_key_vault_key.tfex.name
  key_version        = azurerm_key_vault_key.tfex.version
}
```

## Argument Reference

The following arguments are supported:

* `storage_account_id` - (Required) The id of the storage account to manage the encryption settings for.
* `key_vault_id` - (Required) The ID of the Key Vault.
* `key_name` - (Required) The name of Key Vault key.
* `key_version` - (Required) The version of Key Vault key.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The storage account encryption settings Resource ID.
* `key_vault_uri` - The base URI of the Key Vault.

---

## Import

Storage Accounts Encryption Settings can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_account_customer_managed_key.tfex /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Storage/storageAccounts/myaccount
```