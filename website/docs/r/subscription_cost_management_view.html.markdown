---
subcategory: "Cost Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_subscription_cost_management_view"
description: |-
  Manages an Azure Cost Management View for a Subscription.
---

# azurerm_subscription_cost_management_view

Manages an Azure Cost Management View for a Subscription.

## Example Usage

```hcl
resource "azurerm_subscription_cost_management_view" "example" {
  name         = "example"
  display_name = "Cost View per Month"
  chart_type   = "StackedColumn"
  accumulated  = false

  subscription_id = "/subscription/00000000-0000-0000-0000-000000000000"

  report_type = "Usage"
  timeframe   = "MonthToDate"

  dataset {
    granularity = "Monthly"

    aggregation {
      name        = "totalCost"
      column_name = "Cost"
    }
  }
  pivot {
    type = "Dimension"
    name = "ServiceName"
  }
  pivot {
    type = "Dimension"
    name = "ResourceLocation"
  }
  pivot {
    type = "Dimension"
    name = "ResourceGroupName"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `accumulated` - (Required) Whether the costs data in the Cost Management View are accumulated over time. Changing this forces a new Cost Management View for a Subscription to be created.

* `chart_type` - (Required) Chart type of the main view in Cost Analysis. Possible values are `Area`, `GroupedColumn`, `Line`, `StackedColumn` and `Table`.

* `dataset` - (Required) A `dataset` block as defined below.

* `display_name` - (Required) User visible input name of the Cost Management View.

* `name` - (Required) The name which should be used for this Cost Management View for a Subscription. Changing this forces a new Cost Management View for a Subscription to be created.

* `report_type` - (Required) The type of the report. The only possible value is `Usage`.

* `subscription_id` - (Required) The ID of the Subscription this View is scoped to. Changing this forces a new Cost Management View for a Subscription to be created.

* `timeframe` - (Required) The time frame for pulling data for the report. Possible values are `Custom`, `MonthToDate`, `WeekToDate` and `YearToDate`.

---

* `kpi` - (Optional) One or more `kpi` blocks as defined below, to show in Cost Analysis UI.

* `pivot` - (Optional) One or more `pivot` blocks as defined below, containing the configuration of 3 sub-views in the Cost Analysis UI. Non table views should have three pivots.

---

A `aggregation` block supports the following:

* `name` - (Required) The name which should be used for this aggregation. Changing this forces a new Cost Management View for a Subscription to be created.

* `column_name` - (Required) The name of the column to aggregate. Changing this forces a new Cost Management View for a Subscription to be created.

---

A `dataset` block supports the following:

* `aggregation` - (Required) One or more `aggregation` blocks as defined above.

* `granularity` - (Required) The granularity of rows in the report. Possible values are `Daily` and `Monthly`.

* `grouping` - (Optional) One or more `grouping` blocks as defined below.

* `sorting` - (Optional) One or more `sorting` blocks as defined below, containing the order by expression to be used in the report

---

A `grouping` block supports the following:

* `name` - (Required) The name of the column to group.

* `type` - (Required) The type of the column. Possible values are `Dimension` and `TagKey`.

---

A `kpi` block supports the following:

* `type` - (Required) KPI type. Possible values are `Budget` and `Forecast`.

---

A `pivot` block supports the following:

* `name` - (Required) The name of the column which should be used for this sub-view in the Cost Analysis UI.

* `type` - (Required) The data type to show in this sub-view. Possible values are `Dimension` and `TagKey`.

---

A `sorting` block supports the following:

* `direction` - (Required) Direction of sort. Possible values are `Ascending` and `Descending`.

* `name` - (Required) The name of the column to sort.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Cost Management View for a Subscription.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cost Management View for a Subscription.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cost Management View for a Subscription.
* `update` - (Defaults to 30 minutes) Used when updating the Cost Management View for a Subscription.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cost Management View for a Subscription.

## Import

Cost Management View for a Subscriptions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_subscription_cost_management_view.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.CostManagement/views/costmanagementview
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.CostManagement` - 2023-08-01
