---
subcategory: "Authorization"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_user_assigned_identity"
description: |-
    Lists User Assigned Identity resources.
---

# List resource: azurerm_user_assigned_identity

Lists User Assigned Identity resources.

## Example Usage

### List all User Assigned Identities in the subscription

```hcl
list "azurerm_user_assigned_identity" "example" {
  provider = azurerm
  config {}
}
```

### List all User Assigned Identities in a Resource Group

```hcl
list "azurerm_user_assigned_identity" "example" {
  provider = azurerm
  config {
    resource_group_name = "example-rg"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `subscription_id` - (Optional) The ID of the Subscription to query. Defaults to the value specified in the Provider Configuration.

* `resource_group_name` - (Optional) The name of the Resource Group to query.
