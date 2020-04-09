---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_scheduled_query_rules_alert"
description: |-
  Manages an AlertingAction Scheduled Query Rules resource within Azure Monitor
---

# azurerm_monitor_scheduled_query_rules_alert

Manages an AlertingAction Scheduled Query Rules resource within Azure Monitor.

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
  description    = "Alert when total results cross threshold"
  enabled        = true
  # Count all requests with server error result code grouped into 5-minute bins
  query       = <<-QUERY
  requests
    | where tolong(resultCode) >= 500
    | summarize count() by bin(timestamp, 5m)
  QUERY
  severity    = 1
  frequency   = 5
  time_window = 30
  trigger {
    operator  = "GreaterThan"
    threshold = 3
  }
}

# Example: Alerting Action with metric trigger
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
  description    = "Query results grouped into AggregatedValue; alert when results cross threshold"
  enabled        = true
  # Count all requests with server error result code grouped into 5-minute bins by HTTP operation
  query       = <<-QUERY
  requests
    | where tolong(resultCode) >= 500
    | summarize AggregatedValue = count() by operation_Name, bin(timestamp, 5m)
QUERY
  severity    = 1
  frequency   = 5
  time_window = 30
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
  description    = "Query may access data within multiple resources"
  enabled        = true
  # Count requests in multiple log resources and group into 5-minute bins by HTTP operation
  query = format(<<-QUERY
  let a=requests
    | where toint(resultCode) >= 500
    | extend fail=1; let b=app('%s').requests
    | where toint(resultCode) >= 500 | extend fail=1; a
    | join b on fail
QUERY
  , azurerm_application_insights.example2.id)
  severity    = 1
  frequency   = 5
  time_window = 30
  trigger {
    operator  = "GreaterThan"
    threshold = 3
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the scheduled query rule. Changing this forces a new resource to be created.
* `resource_group_name` - (Required) The name of the resource group in which to create the scheduled query rule instance.
* `data_source_id` - (Required) The resource URI over which log search query is to be run.
* `frequency` - (Required) Frequency (in minutes) at which rule condition should be evaluated.  Values must be between 5 and 1440 (inclusive).
* `query` - (Required) Log search query.
* `time_window` - (Required) Time window for which data needs to be fetched for query (must be greater than or equal to `frequency`).  Values must be between 5 and 2880 (inclusive).
* `trigger` - (Required) The condition that results in the alert rule being run.
* `action` - (Required) An `action` block as defined below.
* `authorized_resource_ids` - (Optional) List of Resource IDs referred into query.
* `description` - (Optional) The description of the scheduled query rule.
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

* `id` - The ID of the scheduled query rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Scheduled Query Rule Alert.
* `update` - (Defaults to 30 minutes) Used when updating the Scheduled Query Rule Alert.
* `read` - (Defaults to 5 minutes) Used when retrieving the Scheduled Query Rule Alert.
* `delete` - (Defaults to 30 minutes) Used when deleting the Scheduled Query Rule Alert.

## Import

Scheduled Query Rule Alerts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_scheduled_query_rules_alert.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Insights/scheduledQueryRules/myrulename
```
