---
subcategory: "Compute Fleet"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_compute_fleet"
description: |-
  Lists Compute Fleet resources.
---

# List resource: azurerm_compute_fleet

Lists Compute Fleet resources.

## Example Usage

### List all Compute Fleets in the subscription

```hcl
list "azurerm_compute_fleet" "example" {
  provider = azurerm
  config {}
}
```

### List all Compute Fleets in a specific resource group

```hcl
list "azurerm_compute_fleet" "example" {
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
