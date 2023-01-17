---
subcategory: "Consumption"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_consumption_budget_management_group"
description: |-
  Manages a Consumption Budget for a Management Group.
---

# azurerm_consumption_budget_management_group

Manages a Consumption Budget for a Management Group.

## Example Usage

```hcl
resource "azurerm_management_group" "example" {
  display_name = "example"
}

resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "eastus"
}

resource "azurerm_consumption_budget_management_group" "example" {
  name                = "example"
  management_group_id = azurerm_management_group.example.id

  amount     = 1000
  time_grain = "Monthly"

  time_period {
    start_date = "2022-06-01T00:00:00Z"
    end_date   = "2022-07-01T00:00:00Z"
  }

  filter {
    dimension {
      name = "ResourceGroupName"
      values = [
        azurerm_resource_group.example.name,
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
  }

  notification {
    enabled        = false
    threshold      = 100.0
    operator       = "GreaterThan"
    threshold_type = "Forecasted"

    contact_emails = [
      "foo@example.com",
      "bar@example.com",
    ]
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Management Group Consumption Budget. Changing this forces a new resource to be created.

* `management_group_id` - (Required) The ID of the Management Group. Changing this forces a new resource to be created.

* `amount` - (Required) The total amount of cost to track with the budget.

* `time_grain` - (Optional) The time covered by a budget. Tracking of the amount will be reset based on the time grain. Must be one of `BillingAnnual`, `BillingMonth`, `BillingQuarter`, `Annually`, `Monthly` and `Quarterly`. Defaults to `Monthly`. Changing this forces a new resource to be created.

* `time_period` - (Required) A `time_period` block as defined below.

* `notification` - (Required) One or more `notification` blocks as defined below.

* `filter` - (Optional) A `filter` block as defined below.

---

A `filter` block supports the following:

* `dimension` - (Optional) One or more `dimension` blocks as defined below to filter the budget on.

* `tag` - (Optional) One or more `tag` blocks as defined below to filter the budget on.

* `not` - (Optional) A `not` block as defined below to filter the budget on. This is deprecated as the API no longer supports it and will be removed in version 4.0 of the provider.

---

A `not` block supports the following:

* `dimension` - (Optional) One `dimension` block as defined below to filter the budget on. Conflicts with `tag`.

* `tag` - (Optional) One `tag` block as defined below to filter the budget on. Conflicts with `dimension`.

---

A `notification` block supports the following:

* `operator` - (Required) The comparison operator for the notification. Must be one of `EqualTo`, `GreaterThan`, or `GreaterThanOrEqualTo`.

* `threshold` - (Required) Threshold value associated with a notification. Notification is sent when the cost exceeded the threshold. It is always percent and has to be between 0 and 1000.

* `contact_emails` - (Required) Specifies a list of email addresses to send the budget notification to when the threshold is exceeded.

* `threshold_type` - (Optional) The type of threshold for the notification. This determines whether the notification is triggered by forecasted costs or actual costs. The allowed values are `Actual` and `Forecasted`. Default is `Actual`. Changing this forces a new resource to be created.

* `enabled` - (Optional) Should the notification be enabled? Defaults to `true`.

---

A `dimension` block supports the following:

* `name` - (Required) The name of the column to use for the filter. The allowed values are `ChargeType`, `Frequency`, `InvoiceId`, `Meter`, `MeterCategory`, `MeterSubCategory`, `PartNumber`, `PricingModel`, `Product`, `ProductOrderId`, `ProductOrderName`, `PublisherType`, `ReservationId`, `ReservationName`, `ResourceGroupName`, `ResourceGuid`, `ResourceId`, `ResourceLocation`, `ResourceType`, `ServiceFamily`, `ServiceName`, `SubscriptionID`, `SubscriptionName`, `UnitOfMeasure`.

* `operator` - (Optional) The operator to use for comparison. The allowed values are `In`.

* `values` - (Required) Specifies a list of values for the column.

---

A `tag` block supports the following:

* `name` - (Required) The name of the tag to use for the filter.

* `operator` - (Optional) The operator to use for comparison. The allowed values are `In`.

* `values` - (Required) Specifies a list of values for the tag.

---

A `time_period` block supports the following:

* `start_date` - (Required) The start date for the budget. The start date must be first of the month and should be less than the end date. Budget start date must be on or after June 1, 2017. Future start date should not be more than twelve months. Past start date should be selected within the timegrain period. Changing this forces a new resource to be created.

* `end_date` - (Optional) The end date for the budget. If not set this will be 10 years after the start date.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Management Group Consumption Budget.

* `etag` - (Optional) The ETag of the Management Group Consumption Budget.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Management Group Consumption Budget.
* `read` - (Defaults to 5 minutes) Used when retrieving the Management Group Consumption Budget.
* `update` - (Defaults to 30 minutes) Used when updating the Management Group Consumption Budget.
* `delete` - (Defaults to 30 minutes) Used when deleting the Management Group Consumption Budget.

## Import

Management Group Consumption Budgets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_consumption_budget_management_group.example /providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Consumption/budgets/budget1
```
