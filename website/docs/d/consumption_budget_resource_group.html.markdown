---
subcategory: "Consumption"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_consumption_budget_resource_group"
description: |-
  Gets information about an existing Consumption Budget for a specific resource group.
---

# Data Source: azurerm_consumption_budget_resource_group

Use this data source to access information about an existing Consumption Budget for a specific resource group.

## Example Usage

```hcl
data "azurerm_consumption_budget_resource_group" "example" {
  name              = "existing"
  resource_group_id = azurerm_resource_group.example.id
}

output "id" {
  value = data.azurerm_consumption_budget_resource_group.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Consumption Budget.

* `resource_group_id` - (Required) The ID of the subscription.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Consumption Budget.

* `amount` - The total amount of cost to track with the budget.

* `filter` - A `filter` block as defined below.

* `notification` - A `notification` block as defined below.

* `time_grain` - The time covered by a budget.

* `time_period` - A `time_period` block as defined below.

---

A `filter` block exports the following:

* `dimension` - A `dimension` block as defined above.

* `not` - A `not` block as defined below.

* `tag` - A `tag` block as defined below.

-> **Note:** The order of multiple filter entries is not guaranteed to be consistent by the API.

---

A `not` block exports the following:

* `dimension` - A `dimension` block as defined below.

* `tag` - A `tag` block as defined below.

---

A `dimension` block exports the following:

* `name` - The name of the column to use for the filter.

* `operator` -  The operator to use for comparison.

* `values` - A `values` block as defined below.

---

A `notification` block exports the following:

* `contact_emails` - A list of email addresses to send the budget notification to when the threshold is exceeded.

* `contact_groups` - A list of Action Group IDs to send the budget notification to when the threshold is exceeded.

* `contact_roles` - A list of contact roles to send the budget notification to when the threshold is exceeded.

* `enabled` - Whether the notification is enabled.

* `operator` - The comparison operator for the notification.

* `threshold` - Threshold value associated with the notification.

-> **Note:** The order of multiple filter entries is not guaranteed to be consistent by the API.

---

A `tag` block exports the following:

* `name` - The name of the tag used for the filter.

* `operator` - The operator used for comparison.

* `values` - A list of values for the tag.

---

A `time_period` block exports the following:

* `end_date` - The end date for the budget.

* `start_date` - The start date for the budget.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Consumption Budget.
