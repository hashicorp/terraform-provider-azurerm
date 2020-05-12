---
subcategory: "Sentinel"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sentinel_alert_rule_action"
description: |-
  Manages a Sentinel Alert Rule Action.
---

# azurerm_sentinel_alert_rule_action

Manages a Sentinel Alert Rule Action.

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
  name                       = "example-alert-rule"
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

resource "azurerm_logic_app_workflow" "example" {
  name                = "example-workflow"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_logic_app_trigger_custom" "example" {
  name         = "When_a_response_to_an_Azure_Sentinel_alert_is_triggered"
  logic_app_id = azurerm_logic_app_workflow.example.id

  body = <<BODY
{
    "type": "ApiConnectionWebhook",
    "inputs": {
        "body": {
            "callback_url": "@{listCallbackUrl()}"
        },
        "host": {
            "connection": {
                "name": "@parameters('$connections')['azuresentinel']['connectionId']"
            }
        },
        "path": "/subscribe"
    }
}
BODY
}

resource "azurerm_sentinel_alert_rule_action" "example" {
  name                   = "example-alert-rule-action"
  rule_id                = azurerm_sentinel_alert_rule_scheduled.example.id
  logic_app_id           = azurerm_logic_app_trigger_custom.example.logic_app_id
  logic_app_trigger_name = azurerm_logic_app_trigger_custom.example.name
  depends_on             = [azurerm_logic_app_trigger_custom.example]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Sentinel Alert Rule Action. Changing this forces a new Sentinel Alert Rule Action to be created.

* `logic_app_id` - (Required) The ID of the Logic App Workflow. Changing this forces a new Sentinel Alert Rule Action to be created.

* `logic_app_trigger_name` - (Required) The Name of the Logic App Workflow Trigger. Changing this forces a new Sentinel Alert Rule Action to be created.

* `rule_id` - (Required) The ID of the Sentinel Alert Rule where the Sentinel Alert Rule Action should exist. Changing this forces a new Sentinel Alert Rule Action to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Sentinel Alert Rule Action.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Sentinel Alert Rule Action.
* `read` - (Defaults to 5 minutes) Used when retrieving the Sentinel Alert Rule Action.
* `delete` - (Defaults to 30 minutes) Used when deleting the Sentinel Alert Rule Action.

## Import

Sentinel Alert Rule Actions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sentinel_alert_rule_action.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/alertRules/rule1/actions/action1
```
