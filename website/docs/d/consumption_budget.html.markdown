---
subcategory: "Consumption"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_consumption_budget"
description: |-
  Gets information about an existing Consumption Budget.
---

# Data Source: azurerm_consumption_budget

Use this data source to access information about an existing Consumption Budget.

## Example Usage

```hcl
data "azurerm_consumption_budget" "example" {
  name = "existing"
  subscription_id = "TODO"
}

output "id" {
  value = data.azurerm_consumption_budget.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Consumption Budget. Changing this forces a new Consumption Budget to be created.

* `subscription_id` - (Required) The ID of the TODO.

---

* `resource_group_id` - (Optional) The ID of the TODO.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Consumption Budget.

* `amount` - TODO.

* `filter` - A `filter` block as defined below.

* `notification` - A `notification` block as defined below.

* `time_grain` - TODO.

* `time_period` - A `time_period` block as defined below.

---

A `dimension` block exports the following:

* `name` - The name of this TODO.

* `operator` - TODO.

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

* `enabled` - Is the TODO enabled?

* `operator` - TODO.

* `threshold` - TODO.

---

A `tag` block exports the following:

* `name` - The name of this TODO.

* `operator` - TODO.

* `values` - A `values` block as defined below.

---

A `time_period` block exports the following:

* `end_date` - TODO.

* `start_date` - TODO.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Consumption Budget.