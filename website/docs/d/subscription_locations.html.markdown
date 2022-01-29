---
subcategory: "Base"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_subscription_locations"
description: |-
  Get information about the available extended locations of the subscription.
---

# Data Source: azurerm_subscription_locations

Use this data source to access information about the available extended locations of the subscription.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

data "azurerm_subscription_locations" "example" {
  subscription_id = data.azurerm_client_config.current.subscription_id
}
```

## Argument Reference

* `subscription_id` - (Required) The ID of the subscription.

## Attributes Reference

* `id` - The ID of the Subscription Locations.

* `locations` - One or more `locations` blocks as defined below.

---

The `locations` block contains:

* `extended_location` - The available extended location for the subscription.
  
* `location` - The available region for the subscription.
  
* `type` - The type of the Subscription Location.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Subscription Locations.
