---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_scheduled_query_rules_action"
sidebar_current: "docs-azurerm-datasource-monitor-scheduled-query-rules-action"
description: |-
  Get information about the specified AlertingAction Scheduled Query Rule.
---

# Data Source: azurerm_monitor_scheduled_query_rules

Use this data source to access the properties of an AlertingAction Scheduled Query Rule.

## Example Usage

```hcl
data "azurerm_monitor_scheduled_query_rules_action" "example" {
  resource_group_name = "terraform-example-rg"
  name                = "tfex-queryrule"
}

output "query_rule_id" {
  value = "${data.azurerm_monitor_scheduled_query_rules_action.example.id}"
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the Scheduled Query Rule.
* `resource_group_name` - (Required) Specifies the name of the resource group the Scheduled Query Rule is located in.

## Attributes Reference

* `id` - The ID of the Scheduled Query Rule.
* `azns_action` - An `azns_action` block as defined below.
* `authorized_resources` - List of Resource IDs referred into query.
* `data_source_id` - The resource uri over which log search query is to be run.
* `description` - The description of the Scheduled Query Rule.
* `enabled` - Whether this scheduled query rule is enabled.
* `frequency` - Frequency (in minutes) at which rule condition should be evaluated.
* `query` - Log search query. Required for action type - `alerting_action`.
* `query_type` - Must equal "ResultCount".
* `time_window` - Time window for which data needs to be fetched for query (should be greater than or equal to frequency_in_minutes).
* `severity` - Severity of the alert. Possible values include: 0, 1, 2, 3, or 4.
* `throttling` - Time (in minutes) for which Alerts should be throttled or suppressed.
* `trigger` - A `trigger` block as defined below. The condition that results in the alert rule being run.

---

* `azns_action` supports the following:

* `action_group` - (Optional) List of action group reference resource IDs.
* `custom_webhook_payload` - Custom payload to be sent for all webhook URI in Azure action group.
* `email_subject` - Custom subject override for all email ids in Azure action group.

---

`metricTrigger` supports the following:

* `metricColumn` - (Required) Evaluation of metric on a particular column.
* `metricTriggerType` - (Required) Metric Trigger Type - 'Consecutive' or 'Total'.
* `operator` - (Required) Evaluation operation for rule - 'Equal', 'GreaterThan' or 'LessThan'.
* `threshold` - (Required) The threshold of the metric trigger.

---

`trigger` supports the following:

* `metricTrigger` - (Optional) A `metricTrigger` block as defined above. Trigger condition for metric query rule.
* `operator` - (Required) Evaluation operation for rule - 'Equal', 'GreaterThan' or 'LessThan'.
* `threshold` - (Required) Result or count threshold based on which rule should be triggered.
