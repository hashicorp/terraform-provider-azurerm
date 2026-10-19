---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_availability_set"
description: |-
    Lists Availability Set resources.
---

# List resource: azurerm_availability_set

Lists Availability Set resources.

## Example Usage

### List all Availability Sets in the subscription

```hcl
list "azurerm_availability_set" "example" {
  provider = azurerm
  config {}
}
```

### List all Availability Sets in a specific resource group

```hcl
list "azurerm_availability_set" "example" {
  provider = azurerm
  config {
    resource_group_name = "example-rg"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `resource_group_name` - (Optional) The name of the resource group to query.

* `subscription_id` - (Optional) The Subscription ID to query. Defaults to the value specified in the Provider Configuration.
