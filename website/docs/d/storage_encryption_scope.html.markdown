---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_storage_encryption_scope"
description: |-
  Gets information about an existing Storage Encryption Scope.
---

# Data Source: azurerm_storage_encryption_scope

Use this data source to access information about an existing Storage Encryption Scope.

## Example Usage

```hcl
data "azurerm_storage_account" "example" {
  name                = "storageaccountname"
  resource_group_name = "resourcegroupname"
}

data "azurerm_storage_encryption_scope" "example" {
  name               = "existingStorageES"
  storage_account_id = data.azurerm_storage_account.example.id
}

output "id" {
  value = data.azurerm_storage_encryption_scope.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Storage Encryption Scope.

* `storage_account_id` - (Required) The ID of the Storage Account where this Storage Encryption Scope exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Storage Encryption Scope.

* `key_vault_key_id` - The ID of the Key Vault Key.

* `source` - The source of the Storage Encryption Scope.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Encryption Scope.
