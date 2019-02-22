---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_account_encryption_settings"
sidebar_current: "docs-azurerm-resource-storage-account-encryption-settings"
description: |-
  Manages the encryption settings of an Azure Storage Account.
---

# azurerm_storage_account_encryption_settings

Manages the encryption settings of an Azure Storage Account.

## Example Usage

```hcl
resource "azurerm_storage_account_encryption_settings" "custom" {
  storage_account_id     = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/tfex-RG/providers/Microsoft.Storage/storageAccounts/tfexstorageaccount"

  key_vault {
    key_vault_policy_id  = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/tfex-RG/providers/Microsoft.KeyVault/vaults/tfex-key-vault/objectId/00000000-0000-0000-0000-000000000000"
    key_vault_id         = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/tfex-RG/providers/Microsoft.KeyVault/vaults/tfex-key-vault"
    key_name             = "tfex-key"
    key_version          = "955b9ad9579e4501a311df5493bacd02"
  }
}
```

## Example Usage with User Managed Key Vault Key

```hcl
resource "azurerm_resource_group" "tfex" {
  name     = "tfex-RG"
  location = "westeurope"
}

resource "azurerm_key_vault" "tfex" {
  name                        = "tfex-key-vault"
  location                    = "${azurerm_resource_group.tfex.location}"
  resource_group_name         = "${azurerm_resource_group.tfex.name}"
  enabled_for_disk_encryption = true
  tenant_id                   = "00000000-0000-0000-0000-000000000000"

  sku {
    name = "standard"
  }

  tags {
    environment = "testing"
  }
}

resource "azurerm_key_vault_key" "tfex" {
  name         = "tfex-key"
  key_vault_id = "${azurerm_key_vault.tfex.id}"
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]
}

resource "azurerm_key_vault_access_policy" "tfex" {
  key_vault_id       = "${azurerm_key_vault.tfex.id}"
  tenant_id          = "00000000-0000-0000-0000-000000000000"
  object_id          = "${azurerm_storage_account.tfex.identity.0.principal_id}"

  key_permissions    = ["get","create","delete","list","restore","recover","unwrapkey","wrapkey","purge","encrypt","decrypt","sign","verify"]
  secret_permissions = ["get"]
}

resource "azurerm_storage_account" "tfex" {
  name                      = "tfexstorageaccount"
  resource_group_name       = "${azurerm_resource_group.tfex.name}"
  location                  = "${azurerm_resource_group.tfex.location}"
  account_tier              = "Standard"
  account_replication_type  = "GRS"

  identity {
    type = "SystemAssigned"
  }

  tags {
    environment = "testing"
  }
}

resource "azurerm_storage_account_encryption_settings" "custom" {
  storage_account_id     = "${azurerm_storage_account.tfex.id}"

  key_vault {
    key_vault_policy_id  = "${azurerm_key_vault_access_policy.tfex.id}"
    key_vault_id         = "${azurerm_key_vault.tfex.id}"
    key_name             = "${azurerm_key_vault_key.tfex.name}"
    key_version          = "${azurerm_key_vault_key.tfex.version}"
  }
}
```

## Argument Reference

The following arguments are supported:

* `storage_account_id` - (Required) The id of the storage account to manage the encryption settings for.

* `enable_blob_encryption` - (Optional) Boolean flag which controls if Encryption Services are enabled for Blob storage, see [here](https://azure.microsoft.com/en-us/documentation/articles/storage-service-encryption/) for more information. Defaults to `true`.

* `enable_file_encryption` - (Optional) Boolean flag which controls if Encryption Services are enabled for File storage, see [here](https://azure.microsoft.com/en-us/documentation/articles/storage-service-encryption/) for more information. Defaults to `true`.

* `key_vault` - (Optional) A `key_vault` block as documented below.

---

* `key_vault` supports the following:

* `key_vault_id` - (Required) The ID of the Key Vault.
* `key_vault_policy_id` - (Required) The resource ID of the `azurerm_key_vault_access_policy` granting the storage account access to the key vault.
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
terraform import azurerm_storage_account_encryption_settings.tfex /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Storage/storageAccounts/myaccount
```
