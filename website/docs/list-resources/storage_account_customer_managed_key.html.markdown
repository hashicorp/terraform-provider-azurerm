---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_account_customer_managed_key"
description: |-
  Lists Storage Account Customer Managed Keys used for encryption.
---

# List resource: azurerm_storage_account_customer_managed_key

Lists Storage Account Customer Managed Keys used for encryption.

## Example Usage

### List all Storage Account Customer Managed Keys in the subscription

```hcl
list "azurerm_storage_account_customer_managed_key" "example" {
  provider = azurerm
  config {}
}
```

### List all Storage Account Customer Managed Keys in a specific resource group

```hcl
list "azurerm_storage_account_customer_managed_key" "example" {
  provider = azurerm
  config {
    resource_group_name = "example-rg"
  }
}
```

## Argument Reference

This list resource supports the following attributes:

* `resource_group_name` - (Optional) The name of the resource group to query.

* `subscription_id` - (Optional) The Subscription ID to query. Defaults to the value specified in the Provider Configuration.
