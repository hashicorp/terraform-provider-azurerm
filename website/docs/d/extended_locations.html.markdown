---
subcategory: "Base"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_extended_locations"
description: |-
  Get information about the available extended locations per the Azure location for the Subscription ID of the Provider.
---

# Data Source: azurerm_extended_locations

Use this data source to access information about the available extended locations per the Azure Location for the Subscription ID of the Provider.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

data "azurerm_extended_locations" "example" {
  location = "West Europe"
}
```

## Argument Reference

* `location` - The Azure location that the Extended Location corresponds.

## Attributes Reference

* `id` - The ID of the Extended Locations.

* `extended_locations` - The available extended locations for the Azure Location.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Extended Locations.
