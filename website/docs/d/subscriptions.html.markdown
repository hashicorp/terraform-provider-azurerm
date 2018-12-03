---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_subscriptions"
sidebar_current: "docs-azurerm-datasource-subscriptions"
description: |-
  Get information about the available subscriptions.
---

# Data Source: azurerm_subscriptions

Use this data source to access information about all the Subscriptions currently available.

## Example Usage

```hcl
data "azurerm_subscriptions" "available" {}

output "available_subscriptions" {
  value = "${data.azurerm_subscriptions.current.subscriptions}"
}

output "first_available_subscription_display_name" {
  value = "${data.azurerm_subscriptions.current.subscriptions.0.display_name}"
}
```

## Argument Reference

-> **NOTE** The comparisons are case-insensitive.

* `display_name_prefix` - (Optional) This argument can be used to only return subscriptions that start with this given string.
* `display_name_contains` - (Optional) This argument can be used to only return subscriptions that contain the given string.

## Attributes Reference

* `subscriptions` - One or more `subscription` blocks as defined below.

The `subscription` block contains:

* `display_name` - The subscription display name.
* `state` - The subscription state. Possible values are Enabled, Warned, PastDue, Disabled, and Deleted.
* `location_placement_id` - The subscription location placement ID.
* `quota_id` - The subscription quota ID.
* `spending_limit` - The subscription spending limit.
