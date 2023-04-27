---
subcategory: "Sentinel"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sentinel_alert_rule_threat_intelligence"
description: |-
  Manages a Sentinel Threat Intelligence Alert Rule.
---

# azurerm_sentinel_alert_rule_threat_intelligence

Manages a Sentinel Threat Intelligence Alert Rule.

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

resource "azurerm_log_analytics_solution" "example" {
  solution_name         = "SecurityInsights"
  location              = azurerm_resource_group.example.location
  resource_group_name   = azurerm_resource_group.example.name
  workspace_resource_id = azurerm_log_analytics_workspace.example.id
  workspace_name        = azurerm_log_analytics_workspace.example.name

  plan {
    publisher = "Microsoft"
    product   = "OMSGallery/SecurityInsights"
  }
}

data "azurerm_sentinel_alert_rule_template" "example" {
  display_name               = "(Preview) Microsoft Defender Threat Intelligence Analytics"
  log_analytics_workspace_id = azurerm_log_analytics_solution.example.workspace_resource_id
}

resource "azurerm_sentinel_alert_rule_threat_intelligence" "example" {
  name                       = "example-rule"
  log_analytics_workspace_id = azurerm_log_analytics_solution.example.workspace_resource_id
  alert_rule_template_guid   = data.azurerm_sentinel_alert_rule_template.example.name
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Sentinel Threat Intelligence Alert Rule. Changing this forces a new Sentinel Threat Intelligence Alert Rule to be created.

* `log_analytics_workspace_id` - (Required) The ID of the Log Analytics Workspace this Sentinel Threat Intelligence Alert Rule belongs to. Changing this forces a new Sentinel Threat Intelligence Alert Rule to be created.

* `alert_rule_template_guid` - (Required) The GUID of the alert rule template which is used for this Sentinel Threat Intelligence Alert Rule. Changing this forces a new Sentinel Threat Intelligence Alert Rule to be created.

* `enabled` - (Optional) Whether the Threat Intelligence Alert rule enabled? Defaults to `true`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Sentinel NRT Alert Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Sentinel NRT Alert Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Sentinel NRT Alert Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Sentinel NRT Alert Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Sentinel NRT Alert Rule.

## Import

Sentinel Threat Intelligence Alert Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sentinel_alert_rule_threat_intelligence.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/alertRules/rule1
```
