---
subcategory: "Dynatrace"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dynatrace_monitor"
description: |-
  Gets information about an existing Dynatrace monitor.
---

# Data Source: azurerm_dynatrace_monitor

Use this data source to access information about an existing Dynatrace Monitor.

## Example Usage

```hcl

data "azurerm_dynatrace_monitor" "example" {
  name                = "example-dynatracemonitor"
  resource_group_name = "example-resources"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Name of the Dynatrace monitor.

* `resource_group_name` - (Required) The name of the Resource Group where the Dynatrace monitor should exist.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Dynatrace monitor.

* `location` - The Azure Region where the Dynatrace monitor should exist.

* `identity` - The kind of managed identity assigned to this resource.  A `identity` block as defined below.

* `marketplace_subscription` - Flag specifying the Marketplace Subscription Status of the resource. If payment is not made in time, the resource will go in Suspended state.

* `plan` - Billing plan information. A `plan` block as defined below.

* `user` - User's information. A `user` block as defined below.

* `monitoring_enabled` - Flag specifying if the resource monitoring is enabled or disabled.

* `tags` - A mapping of tags to assign to the resource.

---

An `identity` block exports the following:

* `type` - The type of identity used for the resource.

---

A `plan` block exports the following:

* `billing_cycle` - Different billing cycles.

* `effective_date` - Date when plan was applied.

* `plan` - Plan id as published by Dynatrace.

* `usage_type` - Different usage type.

---

A `user` block exports the following:

* `country` - Country of the user.

* `email` - Email of the user used by Dynatrace for contacting them if needed.

* `first_name` - First name of the user.

* `last_name` - Last name of the user.

* `phone_number` - phone number of the user by Dynatrace for contacting them if needed.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Dynatrace monitor.
