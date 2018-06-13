---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_subscriptions"
sidebar_current: "docs-azurerm-datasource-subscriptions"
description: |-
  Get information about the available subscriptions.
---

# Data Source: azurerm_subscriptions

Use this data source to access a list of all Azure subscriptions currently available.

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

There are no arguments available for this data source.

## Attributes Reference

* `subscriptions` - One or more `subscription` blocks as defined below.

The `subscription` block contains:

* `display_name` - The subscription display name.
* `state` - The subscription state. Possible values are Enabled, Warned, PastDue, Disabled, and Deleted.
* `location_placement_id` - The subscription location placement ID.
* `quota_id` - The subscription quota ID.
* `spending_limit` - The subscription spending limit.
