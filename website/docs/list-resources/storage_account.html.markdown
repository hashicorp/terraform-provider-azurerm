---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_account"
description: |-
  Lists Storage Account resources.
---

# List resource: azurerm_storage_account

~> **Note:** The `azurerm_storage_account` List Resource is in beta. Its interface and behaviour may change as the feature evolves, and breaking changes are possible. It is offered as a technical preview without compatibility guarantees until Terraform 1.14 is generally available.

Lists Storage Account resources.

## Example Usage

### List all Storage Accounts in the subscription

```hcl
list "azurerm_storage_account" "example" {
  provider = azurerm
  config {}
}
```

### List all Storage Accounts in a specific resource group

```hcl
list "azurerm_storage_account" "example" {
  provider = azurerm
  config {
    resource_group_name = "example-rg"
  }
}
```

## Arguments Reference

This list resource supports the following attributes:

* `resource_group_name` - (Optional) The name of the resource group to query.

* `subscription_id` - (Optional) The Subscription ID to query. Defaults to the value specified in the Provider Configuration.