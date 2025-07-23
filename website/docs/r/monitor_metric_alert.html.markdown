---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_metric_alert"
description: |-
  Manages a Metric Alert within Azure Monitor
---

# azurerm_monitor_metric_alert

Manages a Metric Alert within Azure Monitor.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "to_monitor" {
  name                     = "examplestorageaccount"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_monitor_action_group" "main" {
  name                = "example-actiongroup"
  resource_group_name = azurerm_resource_group.example.name
  short_name          = "exampleact"

  webhook_receiver {
    name        = "callmyapi"
    service_uri = "http://example.com/alert"
  }
}

resource "azurerm_monitor_metric_alert" "example" {
  name                = "example-metricalert"
  resource_group_name = azurerm_resource_group.example.name
  scopes              = [azurerm_storage_account.to_monitor.id]
  description         = "Action will be triggered when Transactions count is greater than 50."

  criteria {
    metric_namespace = "Microsoft.Storage/storageAccounts"
    metric_name      = "Transactions"
    aggregation      = "Total"
    operator         = "GreaterThan"
    threshold        = 50

    dimension {
      name     = "ApiName"
      operator = "Include"
      values   = ["*"]
    }
  }

  action {
    action_group_id = azurerm_monitor_action_group.main.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Metric Alert. Changing this forces a new resource to be created.
* `resource_group_name` - (Required) The name of the resource group in which to create the Metric Alert instance. Changing this forces a new resource to be created.
* `scopes` - (Required) A set of strings of resource IDs at which the metric criteria should be applied.
* `criteria` - (Optional) One or more (static) `criteria` blocks as defined below.

-> **Note:** One of either `criteria`, `dynamic_criteria` or `application_insights_web_test_location_availability_criteria` must be specified.

* `dynamic_criteria` - (Optional) A `dynamic_criteria` block as defined below.

-> **Note:** One of either `criteria`, `dynamic_criteria` or `application_insights_web_test_location_availability_criteria` must be specified.

* `application_insights_web_test_location_availability_criteria` - (Optional) A `application_insights_web_test_location_availability_criteria` block as defined below.

-> **Note:** One of either `criteria`, `dynamic_criteria` or `application_insights_web_test_location_availability_criteria` must be specified.

* `action` - (Optional) One or more `action` blocks as defined below.
* `enabled` - (Optional) Should this Metric Alert be enabled? Defaults to `true`.
* `auto_mitigate` - (Optional) Should the alerts in this Metric Alert be auto resolved? Defaults to `true`.
* `description` - (Optional) The description of this Metric Alert.
* `frequency` - (Optional) The evaluation frequency of this Metric Alert, represented in ISO 8601 duration format. Possible values are `PT1M`, `PT5M`, `PT15M`, `PT30M` and `PT1H`. Defaults to `PT1M`.
* `severity` - (Optional) The severity of this Metric Alert. Possible values are `0`, `1`, `2`, `3` and `4`. Defaults to `3`.
* `target_resource_type` - (Optional) The resource type (e.g. `Microsoft.Compute/virtualMachines`) of the target resource.

-> **Note:** This is Required when using a Subscription as scope, a Resource Group as scope or Multiple Scopes.

* `target_resource_location` - (Optional) The location of the target resource.

-> **Note:** This is Required when using a Subscription as scope, a Resource Group as scope or Multiple Scopes.

* `window_size` - (Optional) The period of time that is used to monitor alert activity, represented in ISO 8601 duration format. This value must be greater than `frequency`. Possible values are `PT1M`, `PT5M`, `PT15M`, `PT30M`, `PT1H`, `PT6H`, `PT12H` and `P1D`. Defaults to `PT5M`.
* `tags` - (Optional) A mapping of tags to assign to the resource.

---

An `action` block supports the following:

* `action_group_id` - (Required) The ID of the Action Group can be sourced from [the `azurerm_monitor_action_group` resource](./monitor_action_group.html)
* `webhook_properties` - (Optional) The map of custom string properties to include with the post operation. These data are appended to the webhook payload.

---

A `criteria` block supports the following:

* `metric_namespace` - (Required) One of the metric namespaces to be monitored.
* `metric_name` - (Required) One of the metric names to be monitored.
* `aggregation` - (Required) The statistic that runs over the metric values. Possible values are `Average`, `Count`, `Minimum`, `Maximum` and `Total`.
* `operator` - (Required) The criteria operator. Possible values are `Equals`, `GreaterThan`, `GreaterThanOrEqual`, `LessThan` and `LessThanOrEqual`.
* `threshold` - (Required) The criteria threshold value that activates the alert.
* `dimension` - (Optional) One or more `dimension` blocks as defined below.
* `skip_metric_validation` - (Optional) Skip the metric validation to allow creating an alert rule on a custom metric that isn't yet emitted? Defaults to `false`.

---

A `dynamic_criteria` block supports the following:

* `metric_namespace` - (Required) One of the metric namespaces to be monitored.
* `metric_name` - (Required) One of the metric names to be monitored.
* `aggregation` - (Required) The statistic that runs over the metric values. Possible values are `Average`, `Count`, `Minimum`, `Maximum` and `Total`.
* `operator` - (Required) The criteria operator. Possible values are `LessThan`, `GreaterThan` and `GreaterOrLessThan`.
* `alert_sensitivity` - (Required) The extent of deviation required to trigger an alert. Possible values are `Low`, `Medium` and `High`.
* `dimension` - (Optional) One or more `dimension` blocks as defined below.
* `evaluation_total_count` - (Optional) The number of aggregated lookback points. The lookback time window is calculated based on the aggregation granularity (`window_size`) and the selected number of aggregated points. Defaults to `4`.
* `evaluation_failure_count` - (Optional) The number of violations to trigger an alert. Should be smaller or equal to `evaluation_total_count`. Defaults to `4`.
* `ignore_data_before` - (Optional) The [ISO8601](https://en.wikipedia.org/wiki/ISO_8601) date from which to start learning the metric historical data and calculate the dynamic thresholds.
* `skip_metric_validation` - (Optional) Skip the metric validation to allow creating an alert rule on a custom metric that isn't yet emitted? 

---

A `application_insights_web_test_location_availability_criteria` block supports the following:

* `web_test_id` - (Required) The ID of the Application Insights Web Test.
* `component_id` - (Required) The ID of the Application Insights Resource.
* `failed_location_count` - (Required) The number of failed locations.

---

A `dimension` block supports the following:

* `name` - (Required) One of the dimension names.
* `operator` - (Required) The dimension operator. Possible values are `Include`, `Exclude` and `StartsWith`.
* `values` - (Required) The list of dimension values.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the metric alert.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Metric Alert.
* `read` - (Defaults to 5 minutes) Used when retrieving the Metric Alert.
* `update` - (Defaults to 30 minutes) Used when updating the Metric Alert.
* `delete` - (Defaults to 30 minutes) Used when deleting the Metric Alert.

## Import

Metric Alerts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_metric_alert.main /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.Insights/metricAlerts/example-metricalert
```
