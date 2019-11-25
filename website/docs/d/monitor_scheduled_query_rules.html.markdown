---
subcategory: ""
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_scheduled_query_rules"
sidebar_current: "docs-azurerm-datasource-monitor-scheduled-query-rules"
description: |-
  Get information about the specified Scheduled Query Rule.
---

# Data Source: azurerm_monitor_scheduled_query_rules

Use this data source to access the properties of a Scheduled Query Rule.

## Example Usage

```hcl
data "azurerm_monitor_scheduled_query_rules" "example" {
  resource_group_name = "terraform-example-rg"
  name                = "tfex-queryrule"
}

output "query_rule_id" {
  value = "${data.azurerm_monitor_scheduled_query_rules.example.id}"
}
```

## Argument Reference

* `name` - (Required) Specifies the name of the Scheduled Query Rule.
* `resource_group_name` - (Required) Specifies the name of the resource group the Scheduled Query Rule is located in.

## Attributes Reference

* `id` - The ID of the Scheduled Query Rule.
* `action_groups` - List of Action Group resource IDs.
* `authorized_resources` - List of Resource IDs referred into query.
* `custom_webhook_payload` - Custom payload to be sent for all webhook URI in Azure action group
* `data_source_id` - The resource uri over which log search query is to be run.
* `description` - The description of the Scheduled Query Rule.
* `email_subject` - Custom subject override for all email ids in Azure action group.
* `enabled` - Whether this scheduled query rule is enabled.
* `frequency_in_minutes` - Frequency (in minutes) at which rule condition should be evaluated.
* `query` - Log search query. Required for action type - `alerting_action`.
* `query_type` - Must equal "ResultCount".
* `severity` - Severity of the alert. Possible values include: 'Zero', 'One', 'Two', 'Three', 'Four'.
* `throttling` - Time (in minutes) for which Alerts should be throttled or suppressed.
* `time_window_in_minutes` - Time window for which data needs to be fetched for query (should be greater than or equal to frequency_in_minutes).
* `trigger` - The trigger condition that results in the alert rule being run
# FIXME: https://docs.microsoft.com/en-us/rest/api/monitor/scheduledqueryrules/createorupdate#triggercondition<Paste>
