---
subcategory: "Base"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_subscription"
description: |-
  Gets information about an existing Subscription.
---

# Data Source: azurerm_subscription

Use this data source to access information about an existing Subscription.

## Example Usage

```hcl
data "azurerm_subscription" "current" {
}

output "current_subscription_display_name" {
  value = data.azurerm_subscription.current.display_name
}
```

## Argument Reference

* `subscription_id` - (Optional) Specifies the ID of the subscription. If this argument is omitted, the subscription ID of the current Azure Resource Manager provider is used.

## Attributes Reference

* `id` - The ID of the subscription.
* `subscription_id` - The subscription GUID.
* `display_name` - The subscription display name.
* `tenant_id` - The subscription tenant ID.
* `state` - The subscription state. Possible values are Enabled, Warned, PastDue, Disabled, and Deleted.
* `location_placement_id` - The subscription location placement ID.
* `quota_id` - The subscription quota ID.
* `spending_limit` - The subscription spending limit.
* `tags` - A mapping of tags assigned to the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Subscription.
