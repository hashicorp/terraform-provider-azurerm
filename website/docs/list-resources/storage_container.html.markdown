---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_container"
description: |-
  Lists Storage Container resources.
---

# List resource: azurerm_storage_container

Lists Storage Container resources.

## Example Usage

### List all Storage Accounts in the subscription

```hcl
list "azurerm_storage_container" "example" {
  provider = azurerm
  config {
    storage_account_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Storage/storageAccounts/example"
  }
}
```

## Argument Reference

This list resource supports the following attributes:

* `storage_account_id` - (Required) The ID of the Storage Account to query.
