---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_scheduled_query_rules_log"
description: |-
  Manages a LogToMetricAction Scheduled Query Rule within Azure Monitor
---

# azurerm_monitor_scheduled_query_rules_log

Manages a LogToMetricAction Scheduled Query Rule within Azure Monitor.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "monitoring-resources"
  location = "West US"
}

resource "azurerm_application_insights" "example" {
  name                = "appinsights"
  location            = var.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}

resource "azurerm_scheduled_query_rule_log" "example3" {
  name                = format("%s-queryrule3", var.prefix)
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  criteria {
    metric_name = "Average_% Idle Time"
    dimensions {
      name     = "InstanceName"
      operator = "Include"
      values   = [""]
    }
  }
  data_source_id = azurerm_application_insights.example.id
  description    = "Scheduled query rule LogToMetric example"
  enabled        = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Scheduled Query Rule. Changing this forces a new resource to be created.
* `resource_group_name` - (Required) The name of the resource group in which to create the Scheduled Query Rule instance.
* `criteria` - (Required) A `criteria` block as defined below.
* `data_source_id` - (Required) The resource uri over which log search query is to be run.
* `description` - (Optional) The description of the Scheduled Query Rule.
* `enabled` - (Optional) Whether this scheduled query rule is enabled.  Default is `true`.
* `throttling` - (Optional) Time (in minutes) for which Alerts should be throttled or suppressed.  Values must be between 0 and 10000 (
inclusive).

---

`criteria` supports the following:

* `dimension` - (Required) A `dimension` block as defined below.
* `metric_name` - (Required) Name of the metric.  Supported metrics are listed in the Azure Monitor [Microsoft.OperationalInsights/workspaces](https://docs.microsoft.com/en-us/azure/azure-monitor/platform/metrics-supported#microsoftoperationalinsightsworkspaces) metrics namespace.

---

`dimension` supports the following:

* `name` - (Required) Name of the dimension.
* `operator` - (Required) Operator for dimension values, - 'Include'.
* `values` - (Required) List of dimension values.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Scheduled Query Rule.

## Import

Scheduled Query Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_scheduled_query_rules_log.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Insights/scheduledQueryRules/myrulename
```
