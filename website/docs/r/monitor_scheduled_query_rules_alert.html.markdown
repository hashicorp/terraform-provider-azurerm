---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_scheduled_query_rules_alert"
description: |-
  Manages an AlertingAction Scheduled Query Rule within Azure Monitor
---

# azurerm_monitor_scheduled_query_rules_alert

Manages an AlertingAction Scheduled Query Rule within Azure Monitor.

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

resource "azurerm_log_analytics_workspace" "example" {
  name                = "loganalytics"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

# Example: Alerting Action
resource "azurerm_scheduled_query_rule_alert" "example" {
  name                = format("%s-queryrule", var.prefix)
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  action {
    action_group           = []
    email_subject          = "Email Header"
    custom_webhook_payload = "{}"
  }
  data_source_id = azurerm_application_insights.example.id
  description    = "Scheduled query rule Alerting Action example"
  enabled        = true
  frequency      = 5
  query          = "requests | where status_code >= 500 | summarize AggregatedValue = count() by bin(timestamp, 5m)"
  severity       = 1
  time_window    = 30
  trigger {
    operator  = "GreaterThan"
    threshold = 3
    metric_trigger {
      operator            = "GreaterThan"
      threshold           = 1
      metric_trigger_type = "Total"
      metric_column       = "timestamp"
    }
  }
}

# Example: Alerting Action Cross-Resource
resource "azurerm_scheduled_query_rule_alert" "example2" {
  name                = format("%s-queryrule2", var.prefix)
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  authorized_resource_ids = [azurerm_application_insights.example.id,
  azurerm_log_analytics_workspace.example.id]
  action {
    action_group           = []
    email_subject          = "Email Header"
    custom_webhook_payload = "{}"
  }
  data_source_id = azurerm_application_insights.example.id
  description    = "Scheduled query rule Alerting Action cross-resource example"
  enabled        = true
  frequency      = 5
  query          = format("let a=workspace('%s').Perf | where Computer='dependency' and TimeGenerated > ago(1h) | where ObjectName == 'Processor' and CounterName == '%% Processor Time' | summarize cpu=avg(CounterValue) by bin(TimeGenerated, 1m) | extend ts=tostring(TimeGenerated); let b=requests | where resultCode == '200' and timestamp > ago(1h) | summarize reqs=count() by bin(timestamp, 1m) | extend ts = tostring(timestamp); a | join b on $left.ts == $right.ts | where cpu > 50 and reqs > 5", azurerm_log_analytics_workspace.test.id)
  severity       = "1"
  time_window    = 30
  trigger {
    operator  = "GreaterThan"
    threshold = 3
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Scheduled Query Rule. Changing this forces a new resource to be created.
* `resource_group_name` - (Required) The name of the resource group in which to create the Scheduled Query Rule instance.
* `data_source_id` - (Required) The resource URI over which log search query is to be run.
* `frequency` - (Required) Frequency (in minutes) at which rule condition should be evaluated.  Values must be between 5 and 1440 (inclusive).
* `time_window` - (Required) Time window for which data needs to be fetched for query (must be greater than or equal to `frequency`).  Values must be between 5 and 2880 (inclusive).
* `trigger` - (Required) The condition that results in the alert rule being run.
* `authorized_resource_ids` - (Optional) List of Resource IDs referred into query.
* `action` - (Optional) An `action` block as defined below.
* `description` - (Optional) The description of the Scheduled Query Rule.
* `enabled` - (Optional) Whether this scheduled query rule is enabled.  Default is `true`.
* `severity` - (Optional) Severity of the alert. Possible values include: 0, 1, 2, 3, or 4.
* `throttling` - (Optional) Time (in minutes) for which Alerts should be throttled or suppressed.  Values must be between 0 and 10000 (inclusive).

---

* `action` supports the following:

* `action_group` - (Required) List of action group reference resource IDs.
* `custom_webhook_payload` - (Optional) Custom payload to be sent for all webhook payloads in alerting action.
* `email_subject` - (Optional) Custom subject override for all email ids in Azure action group.

---

`metricTrigger` supports the following:

* `metricColumn` - (Required) Evaluation of metric on a particular column.
* `metricTriggerType` - (Required) Metric Trigger Type - 'Consecutive' or 'Total'.
* `operator` - (Required) Evaluation operation for rule - 'Equal', 'GreaterThan' or 'LessThan'.
* `threshold` - (Required) The threshold of the metric trigger.    Values must be between 0 and 10000 inclusive.

---

`trigger` supports the following:

* `metricTrigger` - (Optional) A `metricTrigger` block as defined above. Trigger condition for metric query rule.
* `operator` - (Required) Evaluation operation for rule - 'Equal', 'GreaterThan' or 'LessThan'.
* `threshold` - (Required) Result or count threshold based on which rule should be triggered.  Values must be between 0 and 10000 inclusive.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Scheduled Query Rule.

## Import

Scheduled Query Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_scheduled_query_rules_alert.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Insights/scheduledQueryRules/myrulename
```
