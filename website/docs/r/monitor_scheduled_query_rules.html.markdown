---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_scheduled_query_rules"
sidebar_current: "docs-azurerm-resource-monitor-scheduled-query-rules"
description: |-
  Manages a Scheduled Query Rule within Azure Monitor
---

# azurerm_monitor_action_group

Manages a Scheduled Query Rule within Azure Monitor.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "monitoring-resources"
  location = "West US"
}

resource "azurerm_application_insights" "example" {
  name                = format("%s-insights", var.prefix)
  location            = var.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}

resource "azurerm_scheduled_query_rule" "example" {
  name                   = format("%s-queryrule", var.prefix)
  location               = azurerm_resource_group.example.location
  resource_group_name    = azurerm_resource_group.example.name

  enabled                = true
  description            = "Scheduled query rule example resource with log query and schedule"
  frequency_in_minutes   = 5
  time_window_in_minutes = 30
  query                  = "requests | where status_code >= 500 | summarize AggregatedValue = count() by bin(TimeGenerated, 5m)"
  data_source_id         = azurerm_application_insights.example.id
  authorized_resources   = [azurerm_application_insights.example.id]
  query_type             = "ResultCount"
  action                 = {}
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Action Group. Changing this forces a new resource to be created.
* `resource_group_name` - (Required) The name of the resource group in which to create the Action Group instance.
* `action_groups` - (Required) List of Action Group resource IDs.
* `authorized_resources` - (Required) List of Resource IDs referred into query.
* `custom_webhook_payload` - (Optional) Custom payload to be sent for all webhook URI in Azure action group
* `data_source_id` - (Required) The resource uri over which log search query is to be run.
* `description` - (Optional) The description of the Scheduled Query Rule.
* `email_subject` - (Optional) Custom subject override for all email ids in Azure action group.
* `enabled` - (Optional) Whether this scheduled query rule is enabled.  Default is `true`.
* `frequency_in_minutes` - (Optional) Frequency (in minutes) at which rule condition should be evaluated.
* `query` - (Required) Log search query. Required for action type - `alerting_action`.
* `query_type` - (Required) Must equal "ResultCount" for now.
* `severity` - (Optional) Severity of the alert. Possible values include: 'Zero', 'One', 'Two', 'Three', 'Four'.
* `throttling` - (Optional) Time (in minutes) for which Alerts should be throttled or suppressed.
* `time_window_in_minutes` - (Optional) Time window for which data needs to be fetched for query (should be greater than or equal to frequency_in_minutes).
* `trigger` - (Optional) The trigger condition that results in the alert rule being run.
# FIXME: https://docs.microsoft.com/en-us/rest/api/monitor/scheduledqueryrules/createorupdate#triggercondition<Paste>

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Scheduled Query Rule.
* `last_updated_time` - Last time the rule was updated in IS08601 format.
* `provisioning_state` - Provisioning state of the scheduled query rule. Possible values include: 'Succeeded', 'Deploying', 'Canceled', 'Failed'

## Import

Scheduled Query Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_scheduled_query_rules.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Insights/scheduledQueryRules/myrulename
```
