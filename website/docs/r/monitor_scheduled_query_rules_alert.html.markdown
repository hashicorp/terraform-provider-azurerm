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

resource "azurerm_application_insights" "example2" {
  name                = "appinsights2"
  location            = var.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}

# Example: Alerting Action with result count trigger
# Alert if more than three HTTP requests returned a >= 500 result code
# in the past 30 minutes
resource "azurerm_monitor_scheduled_query_rules_alert" "example" {
  name                = format("%s-queryrule", var.prefix)
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  action {
    action_group           = []
    email_subject          = "Email Header"
    custom_webhook_payload = "{}"
  }
  data_source_id = azurerm_application_insights.example.id
  description    = "Result count trigger example - alert when total results cross threshold"
  enabled        = true
  frequency      = 5
  query          = "requests | where tolong(resultCode) >= 500 | summarize AggregatedValue = count() by bin(timestamp, 5m)"
  severity       = 1
  time_window    = 30
  trigger {
    operator  = "GreaterThan"
    threshold = 3
  }
}

# Example: Alerting Action with metric trigger
# Alert if more than three HTTP requests returned a >= 500 result code
# in the past 30 minutes that have the same operation (ie: GET /)
resource "azurerm_monitor_scheduled_query_rules_alert" "example" {
  name                = format("%s-queryrule", var.prefix)
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  action {
    action_group           = []
    email_subject          = "Email Header"
    custom_webhook_payload = "{}"
  }
  data_source_id = azurerm_application_insights.example.id
  description    = "Metric trigger example - query results grouped by metric; alert when results per metric_column cross threshold"
  enabled        = true
  frequency      = 5
  query          = "requests | where tolong(resultCode) >= 500 | summarize AggregatedValue = count() by operation_Name, bin(timestamp, 5m)"
  severity       = 1
  time_window    = 30
  trigger {
    operator  = "GreaterThan"
    threshold = 3
    metric_trigger {
      operator            = "GreaterThan"
      threshold           = 1
      metric_trigger_type = "Total"
      metric_column       = "operation_Name"
    }
  }
}

# Example: Alerting Action Cross-Resource
# Enables use of cross-resource queries to analyze query results across 
# multiple Application Insights or Log Analytics resources.
# Alert if more than three HTTP requests returned a >= 500 result code
# in either of the Insights resources in the past 30 minutes.
resource "azurerm_monitor_scheduled_query_rules_alert" "example2" {
  name                = format("%s-queryrule2", var.prefix)
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  authorized_resource_ids = [azurerm_application_insights.example2.id]
  action {
    action_group           = []
    email_subject          = "Email Header"
    custom_webhook_payload = "{}"
  }
  data_source_id = azurerm_application_insights.example.id
  description    = "Scheduled query rule Alerting Action cross-resource example"
  enabled        = true
  frequency      = 5
  query          = format("let a=requests | where toint(resultCode) >= 500 | extend fail=1; let b=app('%s').requests | where toint(resultCode) >= 500 | extend fail=1; a | join b on fail", azurerm_application_insights.example2.id)
  severity       = 1
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
* `query` - (Required) Log search query.
* `time_window` - (Required) Time window for which data needs to be fetched for query (must be greater than or equal to `frequency`).  Values must be between 5 and 2880 (inclusive).
* `trigger` - (Required) The condition that results in the alert rule being run.
* `action` - (Required) An `action` block as defined below.
* `authorized_resource_ids` - (Optional) List of Resource IDs referred into query.
* `description` - (Optional) The description of the Scheduled Query Rule.
* `enabled` - (Optional) Whether this Scheduled Query Rule is enabled.  Default is `true`.
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
