---
subcategory: "Consumption"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_consumption_budget_resource_group"
description: |-
  Manages a Resource Group Consumption Budget.
---

# azurerm_consumption_budget_resource_group

Manages a Resource Group Consumption Budget.

## Example Usage

```hcl
data "azurerm_subscription" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "eastus"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  short_name          = "example"
}

resource "azurerm_consumption_budget_resource_group" "example" {
  name              = "example"
  subscription_id   = data.azurerm_subscription.current.subscription_id
  resource_group_id = azurerm_resource_group.example.id

  amount     = 1000
  time_grain = "Monthly"

  time_period {
    start_date = "2020-11-01T00:00:00Z"
    end_date   = "2020-12-01T00:00:00Z"
  }

  filter {
    dimension {
      name = "ResourceId"
      values = [
        azurerm_monitor_action_group.example.id,
      ]
    }

    tag {
      name = "foo"
      values = [
        "bar",
        "baz",
      ]
    }
  }

  notification {
    enabled   = true
    threshold = 90.0
    operator  = "EqualTo"

    contact_emails = [
      "foo@example.com",
      "bar@example.com",
    ]

    contact_groups = [
      azurerm_monitor_action_group.example.id,
    ]

    contact_roles = [
      "Owner",
    ]
  }

  notification {
    enabled   = false
    threshold = 100.0
    operator  = "GreaterThan"

    contact_emails = [
      "foo@example.com",
      "bar@example.com",
    ]
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Resource Group Consumption Budget. Changing this forces a new Resource Group Consumption Budget to be created.

* `resource_group_id` - (Required) The ID of the Resource Group to create the consumption budget for in the form of /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1. Changing this forces a new Resource Group Consumption Budget to be created.

* `amount` - (Required) The total amount of cost to track with the budget.

* `time_grain` - (Required) The time covered by a budget. Tracking of the amount will be reset based on the time grain. Must be one of `Monthly`, `Quarterly`, `Annually`, `BillingMonth`, `BillingQuarter`, or `BillingYear`. Defaults to `Monthly`.

* `time_period` - (Required) A `time_period` block as defined below.

* `notification` - (Required) One or more `notification` blocks as defined below.

* `filter` - (Optional) A `filter` block as defined below.

---

A `filter` block supports the following:

* `dimension` - (Optional) One or more `dimension` blocks as defined below to filter the budget on.

* `tag` - (Optional) One or more `tag` blocks as defined below to filter the budget on.

* `not` - (Optional) A `not` block as defined below to filter the budget on.

---

A `not` block supports the following:

* `dimension` - (Optional) One `dimension` block as defined below to filter the budget on. Conflicts with `tag`.

* `tag` - (Optional) One `tag` block as defined below to filter the budget on. Conflicts with `dimension`.

---

A `notification` block supports the following:

* `operator` - (Required) The comparison operator for the notification. Must be one of `EqualTo`, `GreaterThan`, or `GreaterThanOrEqualTo`.

* `threshold` - (Required) Threshold value associated with a notification. Notification is sent when the cost exceeded the threshold. It is always percent and has to be between 0 and 1000.

* `contact_emails` - (Optional) Specifies a list of email addresses to send the budget notification to when the threshold is exceeded.

* `contact_groups` - (Optional) Specifies a list of Action Group IDs to send the budget notification to when the threshold is exceeded.

* `contact_roles` - (Optional) Specifies a list of contact roles to send the budget notification to when the threshold is exceeded.

* `enabled` - (Optional) Should the notification be enabled?

---

A `dimension` block supports the following:

* `name` - (Required) The name of the column to use for the filter. The allowed values are

* `operator` - (Optional) The operator to use for comparison. The allowed values are `In`.

* `values` - (Required) Specifies a list of values for the column. The allowed values are `ChargeType`, `Frequency`, `InvoiceId`, `Meter`, `MeterCategory`, `MeterSubCategory`, `PartNumber`, `PricingModel`, `Product`, `ProductOrderId`, `ProductOrderName`, `PublisherType`, `ReservationId`, `ReservationName`, `ResourceGroupName`, `ResourceGuid`, `ResourceId`, `ResourceLocation`, `ResourceType`, `ServiceFamily`, `ServiceName`, `UnitOfMeasure`.

---

A `tag` block supports the following:

* `name` - (Required) The name of the tag to use for the filter.

* `operator` - (Optional) The operator to use for comparison. The allowed values are `In`.

* `values` - (Required) Specifies a list of values for the tag.

---

A `time_period` block supports the following:

* `start_date` - (Required) The start date for the budget. The start date must be first of the month and should be less than the end date. Budget start date must be on or after June 1, 2017. Future start date should not be more than twelve months. Past start date should be selected within the timegrain period. Changing this forces a new Resource Group Consumption Budget to be created.

* `end_date` - (Optional) The end date for the budget. If not set this will be 10 years after the start date.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Resource Group Consumption Budget.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Resource Group Consumption Budget.
* `read` - (Defaults to 5 minutes) Used when retrieving the Resource Group Consumption Budget.
* `update` - (Defaults to 30 minutes) Used when updating the Resource Group Consumption Budget.
* `delete` - (Defaults to 30 minutes) Used when deleting the Resource Group Consumption Budget.

## Import

Resource Group Consumption Budgets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_consumption_budget_resource_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Consumption/budgets/resourceGroup1
```
