---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_account"
description: |-
  Lists Storage Account resources.
---

# List resource: azurerm_storage_account

Lists Storage Account resources.

## Example Usage

### List all Storage Accounts in the subscription

```hcl
list "azurerm_storage_account" "test" {
  provider = azurerm
  config {}
}
```

### List all Storage Accounts in a specific resource group

```hcl
list "azurerm_storage_account" "test" {
  provider = azurerm
  config {
    resource_group_name = "example-rg"
  }
}
```

## Argument Reference

This list resource supports the following attributes:

* `resource_group_name` - (Optional) The name of the resource group to query.

* `subscription_id` - (Optional) The Subscription ID to query. Defaults to the value used in the Provider Config.