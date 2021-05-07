---
subcategory: "Sentinel"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sentinel_alert_rule_machine_learning_behavior_analytics"
description: |-
  Manages a Sentinel Machine Learning Behavior Analytics Rule.
---

# azurerm_sentinel_alert_rule_machine_learning_behavior_analytics

Manages a Sentinel Machine Learning Behavior Analytics Alert Rule.

## Example Usage

```hcl
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

resource "azurerm_sentinel_alert_rule_machine_learning_behavior_analytics" "example" {
  name                       = "example-ml-alert-rule"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.example.id
  alert_rule_template_guid   = "737a2ce1-70a3-4968-9e90-3e6aca836abf"
}


```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this SentinelMachine Learning Behavior Analytics Alert Rule. Changing this forces a new Sentinel Machine Learning Behavior Analytics Alert Rule to be created.

* `log_analytics_workspace_id` - (Required) The ID of the Log Analytics Workspace this SentinelMachine Learning Behavior Analytics Alert Rule belongs to. Changing this forces a new Sentinel Machine Learning Behavior Analytics Alert Rule to be created.

* `alert_rule_template_guid` - (Required) The GUID of the alert rule template which is used for this Sentinel Machine Learning Behavior Analytics Alert Rule. Changing this forces a new Sentinel Machine Learning Behavior Analytics Alert Rule to be created.

* `enabled` - (Optional) Should this Sentinel Machine Learning Behavior Analytics Alert Rule be enabled? Defaults to `true`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Sentinel Machine Learning Behavior Analytics Alert Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Sentinel Machine Learning Behavior Analytics Alert Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Sentinel Machine Learning Behavior Analytics Alert Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Sentinel Machine Learning Behavior Analytics Alert Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Sentinel Machine Learning Behavior Analytics Alert Rule.

## Import

Sentinel Machine Learning Behavior Analytics Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sentinel_alert_rule_machine_learning_behavior_analytics.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/alertRules/rule1
```
