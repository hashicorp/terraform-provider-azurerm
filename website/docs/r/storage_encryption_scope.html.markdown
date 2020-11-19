---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_encryption_scope"
description: |-
  Manages a Storage Encryption Scope.
---

# azurerm_storage_encryption_scope

Manages a Storage Encryption Scope.

~> **Note:** Storage Encryption Scopes are in Preview [more information can be found here](https://docs.microsoft.com/en-us/azure/storage/blobs/encryption-scope-manage).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplesa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_storage_encryption_scope" "example" {
  name               = "microsoftmanaged"
  storage_account_id = azurerm_storage_account.example.id
  source             = "Microsoft.Storage"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Storage Encryption Scope. Changing this forces a new Storage Encryption Scope to be created.

* `source` - (Required) The source of the Storage Encryption Scope. Possible values are `Microsoft.KeyVault` and `Microsoft.Storage`.

* `storage_account_id` - (Required) The ID of the Storage Account where this Storage Encryption Scope is created. Changing this forces a new Storage Encryption Scope to be created.

---

* `key_vault_key_id` - (Optional) The ID of the Key Vault Key. Required when `source` is `Microsoft.KeyVault`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Storage Encryption Scope.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Encryption Scope.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Encryption Scope.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Encryption Scope.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Encryption Scope.

## Import

Storage Encryption Scopes can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_encryption_scope.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Storage/storageAccounts/account1/encryptionScopes/scope1
```
