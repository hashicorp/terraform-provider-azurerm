---
subcategory: "Sentinel"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sentinel_alert_rule_scheduled"
description: |-
  Manages a Sentinel Scheduled Alert Rule.
---

# azurerm_sentinel_alert_rule_scheduled

Manages a Sentinel Scheduled Alert Rule.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-workspace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "pergb2018"
}

resource "azurerm_sentinel_alert_rule_scheduled" "example" {
  name                       = "example"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.example.id
  display_name               = "example"
  severity                   = "High"
  query                      = <<QUERY
AzureActivity |
  where OperationName == "Create or Update Virtual Machine" or OperationName =="Create Deployment" |
  where ActivityStatus == "Succeeded" |
  make-series dcount(ResourceId) default=0 on EventSubmissionTimestamp in range(ago(7d), now(), 1d) by Caller
QUERY
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Sentinel Scheduled Alert Rule. Changing this forces a new Sentinel Scheduled Alert Rule to be created.

* `log_analytics_workspace_id` - (Required) The ID of the Log Analytics Workspace this Sentinel Scheduled Alert Rule belongs to. Changing this forces a new Sentinel Scheduled Alert Rule to be created.

* `display_name` - (Required) The friendly name of this Sentinel Scheduled Alert Rule.

* `severity` - (Required) The alert severity of this Sentinel Scheduled Alert Rule. Possible values are `High`, `Medium`, `Low` and `Informational`.

* `query` - (Required) The query of this Sentinel Scheduled Alert Rule.

---

* `alert_rule_template_guid` - (Optional) The GUID of the alert rule template which is used for this Sentinel Scheduled Alert Rule. Changing this forces a new Sentinel Scheduled Alert Rule to be created.

* `description` - (Optional) The description of this Sentinel Scheduled Alert Rule.

* `enabled` - (Optional) Should the Sentinel Scheduled Alert Rule be enabled? Defaults to `true`.

* `query_frequency` - (Optional) The ISO 8601 timespan duration between two consecutive queries. Defaults to `PT5H`.

* `query_period` - (Optional) The ISO 8601 timespan duration, which determine the time period of the data covered by the query. For example, it can query the past 10 minutes of data, or the past 6 hours of data. Defaults to `PT5H`.

-> **NOTE** `query_period` must larger than or equal to `query_frequency`, which ensures there is no gaps in the overall query coverage.

* `suppression_duration` - (Optional) If `suppression_enabled` is `true`, this is ISO 8601 timespan duration, which specifies the amount of time the query should stop running after alert is generated. Defaults to `PT5H`.

-> **NOTE** `suppression_duration` must larger than or equal to `query_frequency`, otherwise the suppression has no actual effect since no query will happen during the suppression duration.

* `suppression_enabled` - (Optional) Should the Sentinel Scheduled Alert Rulea stop running query after alert is generated? Defaults to `false`.

* `tactics` - (Optional) A list of categories of attacks by which to classify the rule. Possible values are `Collection`, `CommandAndControl`, `CredentialAccess`, `DefenseEvasion`, `Discovery`, `Execution`, `Exfiltration`, `Impact`, `InitialAccess`, `LateralMovement`, `Persistence` and `PrivilegeEscalation`.

* `trigger_operator` - (Optional) The alert trigger operator, combined with `trigger_threshold`, setting alert threshold of this Sentinel Scheduled Alert Rule. Possible values are `Equal`, `GreaterThan`, `LessThan`, `NotEqual`.

* `trigger_threshold` - (Optional) The baseline number of query results generated, combined with `trigger_operator`, setting alert threshold of this Sentinel Scheduled Alert Rule.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Sentinel Scheduled Alert Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Sentinel Scheduled Alert Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Sentinel Scheduled Alert Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Sentinel Scheduled Alert Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Sentinel Scheduled Alert Rule.

## Import

Sentinel Scheduled Alert Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sentinel_alert_rule_scheduled.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/alertRules/rule1
```
