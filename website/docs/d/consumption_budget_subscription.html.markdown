---
subcategory: "Consumption"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_consumption_budget_subscription"
description: |-
  Gets information about an existing Consumption Budget in a subscription.
---

# Data Source: azurerm_consumption_budget_subscription

Use this data source to access information about an existing Consumption Budget.

## Example Usage

```hcl
data "azurerm_consumption_budget_subscription" "example" {
  name            = "existing"
  subscription_id = "/subscriptions/00000000-0000-0000-0000-000000000000/"
}

output "id" {
  value = data.azurerm_consumption_budget.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Consumption Budget. Changing this forces a new Consumption Budget to be created.

* `subscription_id` - (Required) The ID of the subscription.

---

* `resource_group_id` - (Optional) The ID of the resource group.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Consumption Budget.

* `amount` - The total amount of cost to track with the budget.

* `filter` - A `filter` block as defined below.

* `notification` - A `notification` block as defined below.

* `time_grain` - The time covered by a budget. Tracking of the amount will be reset based on the time grain. Must be one of `Monthly`, `Quarterly`, `Annually`, `BillingMonth`, `BillingQuarter`, or `BillingYear`. Defaults to `Monthly`.

* `time_period` - A `time_period` block as defined below.

---

A `dimension` block exports the following:

* `name` - The name of the column to use for the filter. The allowed values are

* `operator` -  The operator to use for comparison. The allowed values are `In`.

* `values` - A `values` block as defined below.

---

A `filter` block exports the following:

* `dimension` - A `dimension` block as defined above.

* `not` - A `not` block as defined below.

* `tag` - A `tag` block as defined below.

---

A `not` block exports the following:

* `dimension` - A `dimension` block as defined above.

* `tag` - A `tag` block as defined below.

---

A `notification` block exports the following:

* `contact_emails` - A `contact_emails` block as defined above.

* `contact_groups` - A `contact_groups` block as defined above.

* `contact_roles` - A `contact_roles` block as defined above.

* `enabled` - Should the notification enabled?

* `operator` - The comparison operator for the notification. Must be one of `EqualTo`, `GreaterThan`, or `GreaterThanOrEqualTo`.

* `threshold` - Threshold value associated with a notification. Notification is sent when the cost exceeded the threshold. It is always percent and has to be between 0 and 1000.

---

A `tag` block exports the following:

* `name` - The name of the tag to use for the filter.

* `operator` - The operator to use for comparison. The allowed values are `In`.

* `values` - A `values` block as defined below.

---

A `time_period` block exports the following:

* `end_date` - The end date for the budget. If not set this will be 10 years after the start date.

* `start_date` - The start date for the budget. The start date must be first of the month and should be less than the end date. Budget start date must be on or after June 1, 2017. Future start date should not be more than twelve months. Past start date should be selected within the timegrain period. Changing this forces a new Subscription Consumption Budget to be created.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Consumption Budget.
