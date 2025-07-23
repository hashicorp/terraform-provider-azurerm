---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_scheduled_query_rules_log"
description: |-
  Manages a LogToMetricAction Scheduled Query Rules resources within Azure Monitor
---

# azurerm_monitor_scheduled_query_rules_log

Manages a LogToMetricAction Scheduled Query Rules resource within Azure Monitor.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "monitoring-resources"
  location = "West Europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "loganalytics"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_monitor_action_group" "example" {
  name                = "example-actiongroup"
  resource_group_name = azurerm_resource_group.example.name
  short_name          = "exampleact"

  webhook_receiver {
    name        = "callmyapi"
    service_uri = "http://example.com/alert"
  }
}

# Example: Creates alert using the new Scheduled Query Rules metric
resource "azurerm_monitor_metric_alert" "example" {
  name                = "example-metricalert"
  resource_group_name = azurerm_resource_group.example.name
  scopes              = [azurerm_log_analytics_workspace.example.id]
  description         = "Action will be triggered when Average_% Idle Time metric is less than 10."
  frequency           = "PT1M"
  window_size         = "PT5M"

  criteria {
    metric_namespace = "Microsoft.OperationalInsights/workspaces"
    metric_name      = "UsedCapacity"
    aggregation      = "Average"
    operator         = "LessThan"
    threshold        = 10
  }

  action {
    action_group_id = azurerm_monitor_action_group.example.id
  }
}

# Example: LogToMetric Action for the named Computer
resource "azurerm_monitor_scheduled_query_rules_log" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  criteria {
    metric_name = "Average_% Idle Time"
    dimension {
      name     = "Computer"
      operator = "Include"
      values   = ["targetVM"]
    }
  }
  data_source_id = azurerm_log_analytics_workspace.example.id
  description    = "Scheduled query rule LogToMetric example"
  enabled        = true
  tags = {
    foo = "bar"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the scheduled query rule. Changing this forces a new resource to be created.
* `resource_group_name` - (Required) The name of the resource group in which to create the scheduled query rule instance. Changing this forces a new resource to be created.
* `location` - (Required) Specifies the Azure Region where the resource should exist. Changing this forces a new resource to be created.
* `criteria` - (Required) A `criteria` block as defined below.
* `data_source_id` - (Required) The resource URI over which log search query is to be run. Changing this forces a new resource to be created.
* `authorized_resource_ids` - (Optional) A list of IDs of Resources referred into query.
* `description` - (Optional) The description of the scheduled query rule.
* `enabled` - (Optional) Whether this scheduled query rule is enabled. Default is `true`.
* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `criteria` block supports the following:

* `dimension` - (Required) A `dimension` block as defined below.
* `metric_name` - (Required) Name of the metric. Supported metrics are listed in the Azure Monitor [Microsoft.OperationalInsights/workspaces](https://docs.microsoft.com/azure/azure-monitor/platform/metrics-supported#microsoftoperationalinsightsworkspaces) metrics namespace.

---

The `dimension` block supports the following:

* `name` - (Required) Name of the dimension.
* `operator` - (Optional) Operator for dimension values, - 'Include'. Defaults to `Include`.
* `values` - (Required) List of dimension values.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the scheduled query rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Scheduled Query Rule Log.
* `read` - (Defaults to 5 minutes) Used when retrieving the Scheduled Query Rule Log.
* `update` - (Defaults to 30 minutes) Used when updating the Scheduled Query Rule Log.
* `delete` - (Defaults to 30 minutes) Used when deleting the Scheduled Query Rule Log.

## Import

Scheduled Query Rule Log can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_scheduled_query_rules_log.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Insights/scheduledQueryRules/myrulename
```
