---
subcategory: "Log Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_analytics_linked_service"
description: |-
  Manages a Log Analytics (formally Operational Insights) Linked Service.
---

# azurerm_log_analytics_linked_service

Links a Log Analytics (formally Operational Insights) Workspace to another resource. The (currently) only linkable service is an Azure Automation Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "resourcegroup-01"
  location = "West Europe"
}

resource "azurerm_automation_account" "example" {
  name                = "automation-01"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "Basic"

  tags = {
    environment = "development"
  }
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "workspace-01"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_log_analytics_linked_service" "example" {
  resource_group_name = azurerm_resource_group.example.name
  workspace_name      = azurerm_log_analytics_workspace.example.name
  resource_id         = azurerm_automation_account.example.id
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group in which the Log Analytics Linked Service is created. Changing this forces a new resource to be created.

* `workspace_name` - (Required) Name of the Log Analytics Workspace that will contain the linkedServices resource. Changing this forces a new resource to be created.

* `linked_service_name` - (Optional) Name of the type of linkedServices resource to connect to the Log Analytics Workspace specified in `workspace_name`. Currently it defaults to and only supports `automation` as a value. Changing this forces a new resource to be created.

* `resource_id` - (Required) The ID of the Resource that will be linked to the workspace. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The Log Analytics Linked Service ID.

* `name` - The automatically generated name of the Linked Service. This cannot be specified. The format is always `<workspace_name>/<linked_service_name>` e.g. `workspace1/Automation`

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Log Analytics Workspace.
* `update` - (Defaults to 30 minutes) Used when updating the Log Analytics Workspace.
* `read` - (Defaults to 5 minutes) Used when retrieving the Log Analytics Workspace.
* `delete` - (Defaults to 30 minutes) Used when deleting the Log Analytics Workspace.

## Import

Log Analytics Workspaces can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_log_analytics_linked_service.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/linkedservices/automation
```
