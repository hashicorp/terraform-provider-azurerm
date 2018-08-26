---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_subscription"
sidebar_current: "docs-azurerm-datasource-subscription-x"
description: |-
  Get information about the specified subscription.
---

# Data Source: azurerm_subscription

Use this data source to access the properties of an Azure subscription.

## Example Usage

```hcl
data "azurerm_subscription" "current" {}

output "current_subscription_display_name" {
  value = "${data.azurerm_subscription.current.display_name}"
}
```

## Argument Reference

* `subscription_id` - (Optional) Specifies the ID of the subscription. If this argument is omitted, the subscription ID of the current Azure Resource Manager provider is used.

## Attributes Reference

* `display_name` - The subscription display name.
* `state` - The subscription state. Possible values are Enabled, Warned, PastDue, Disabled, and Deleted.
* `location_placement_id` - The subscription location placement ID.
* `quota_id` - The subscription quota ID.
* `spending_limit` - The subscription spending limit.
