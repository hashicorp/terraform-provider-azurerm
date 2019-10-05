---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_metric_alert"
sidebar_current: "docs-azurerm-resource-monitor-metric-alert-x"
description: |-
  Manages a Metric Alert within Azure Monitor
---

# azurerm_monitor_metric_alert

Manages a Metric Alert within Azure Monitor.

## Example Usage

```hcl
resource "azurerm_resource_group" "main" {
  name     = "example-resources"
  location = "West US"
}

resource "azurerm_storage_account" "to_monitor" {
  name                     = "examplestorageaccount"
  resource_group_name      = "${azurerm_resource_group.main.name}"
  location                 = "${azurerm_resource_group.main.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_monitor_action_group" "main" {
  name                = "example-actiongroup"
  resource_group_name = "${azurerm_resource_group.main.name}"
  short_name          = "exampleact"

  webhook_receiver {
    name        = "callmyapi"
    service_uri = "http://example.com/alert"
  }
}

resource "azurerm_monitor_metric_alert" "test" {
  name                = "example-metricalert"
  resource_group_name = "${azurerm_resource_group.main.name}"
  scopes              = ["${azurerm_storage_account.to_monitor.id}"]
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
    action_group_id = "${azurerm_monitor_action_group.main.id}"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Metric Alert. Changing this forces a new resource to be created.
* `resource_group_name` - (Required) The name of the resource group in which to create the Metric Alert instance.
* `scopes` - (Required) A set of strings of resource IDs at which the metric criteria should be applied.
* `criteria` - (Required) One or more `criteria` blocks as defined below.
* `action` - (Optional) One or more `action` blocks as defined below.
* `enabled` - (Optional) Should this Metric Alert be enabled? Defaults to `true`.
* `auto_mitigate` - (Optional) Should the alerts in this Metric Alert be auto resolved? Defaults to `false`.
* `description` - (Optional) The description of this Metric Alert.
* `frequency` - (Optional) The evaluation frequency of this Metric Alert, represented in ISO 8601 duration format. Possible values are `PT1M`, `PT5M`, `PT15M`, `PT30M` and `PT1H`. Defaults to `PT1M`.
* `severity` - (Optional) The severity of this Metric Alert. Possible values are `0`, `1`, `2`, `3` and `4`. Defaults to `3`.
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
* `operator` - (Required) The criteria operator. Possible values are `Equals`, `NotEquals`, `GreaterThan`, `GreaterThanOrEqual`, `LessThan` and `LessThanOrEqual`.
* `threshold` - (Required) The criteria threshold value that activates the alert.
* `dimension` - (Optional) One or more `dimension` blocks as defined below.

---

A `dimension` block supports the following:

* `name` - (Required) One of the dimension names.
* `operator` - (Required) The dimension operator. Possible values are `Include` and `Exclude`.
* `values` - (Required) The list of dimension values.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the metric alert.

## Import

Metric Alerts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_metric_alert.main /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/microsoft.insights/metricalerts/example-metricalert
```
