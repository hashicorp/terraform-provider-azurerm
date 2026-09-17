---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_capacity_reservation_group"
description: |-
    Lists Capacity Reservation Group resources.
---

# List resource: azurerm_capacity_reservation_group

Lists Capacity Reservation Group resources.

## Example Usage

### List all Capacity Reservation Groups in the subscription

```hcl
list "azurerm_capacity_reservation_group" "example" {
  provider = azurerm
  config {}
}
```

### List all Capacity Reservation Groups in a specific resource group

```hcl
list "azurerm_capacity_reservation_group" "example" {
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
