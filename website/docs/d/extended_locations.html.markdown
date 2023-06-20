---
subcategory: "Base"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_extended_locations"
description: |-
  This data source return the available Extended Locations for a specific Azure Region.
---

# Data Source: azurerm_extended_locations

This data source return the available Extended Locations for a specific Azure Region.

## Example Usage

```hcl
data "azurerm_extended_locations" "example" {
  location = "West Europe"
}
```

## Argument Reference

* `location` - The Azure location to retrieve the Extended Locations for.

## Attributes Reference

* `id` - The ID of Location within this Subscription.

* `extended_locations` - The available extended locations for the Azure Location.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Extended Locations.
