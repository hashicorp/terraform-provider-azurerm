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

  time_period {
    start_date = "TODO"    
  }
  amount = 1.23456

  notification {
    threshold = 42
    operator = "TODO"    
  }
  resource_group_id = "TODO"
}

output "id" {
  value = data.azurerm_consumption_budget.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `amount` - (Required) TODO.

* `name` - (Required) The name of this Consumption Budget. Changing this forces a new Consumption Budget to be created.

* `notification` - (Required) One or more `notification` blocks as defined below.

* `resource_group_id` - (Required) The ID of the TODO. Changing this forces a new Consumption Budget to be created.

* `subscription_id` - (Required) The ID of the TODO. Changing this forces a new Consumption Budget to be created.

* `time_period` - (Required) A `time_period` block as defined below.

---

* `filter` - (Optional) A `filter` block as defined below.

* `time_grain` - (Optional) TODO. Changing this forces a new Consumption Budget to be created.

---

A `dimension` block supports the following:

* `name` - (Required) The name which should be used for this TODO.

* `values` - (Required) Specifies a list of TODO.

* `operator` - (Optional) TODO.

---

A `filter` block supports the following:

* `dimension` - (Optional) One or more `dimension` blocks as defined above.

* `not` - (Optional) A `not` block as defined below.

* `tag` - (Optional) One or more `tag` blocks as defined below.

---

A `not` block supports the following:

* `dimension` - (Optional) A `dimension` block as defined above.

* `tag` - (Optional) A `tag` block as defined below.

---

A `notification` block supports the following:

* `operator` - (Required) TODO.

* `threshold` - (Required) TODO.

* `contact_emails` - (Optional) Specifies a list of TODO.

* `contact_groups` - (Optional) Specifies a list of TODO.

* `contact_roles` - (Optional) Specifies a list of TODO.

* `enabled` - (Optional) Should the TODO be enabled?

---

A `tag` block supports the following:

* `name` - (Required) The name which should be used for this TODO.

* `values` - (Required) Specifies a list of TODO.

* `operator` - (Optional) TODO.

---

A `time_period` block supports the following:

* `start_date` - (Required) TODO. Changing this forces a new Consumption Budget to be created.

* `end_date` - (Optional) TODO.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Consumption Budget.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Consumption Budget.