---
subcategory: "Log Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_analytics_linked_service"
description: |-
  Manages a Log Analytics Linked Service.
---

# azurerm_log_analytics_linked_service

Manages a Log Analytics Linked Service.

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
  workspace_id        = azurerm_log_analytics_workspace.example.id
  read_access_id      = azurerm_automation_account.example.id
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group in which the Log Analytics Linked Service is created. Changing this forces a new resource to be created.

* `workspace_id` - (Required) The ID of the Log Analytics Workspace that will contain the Log Analytics Linked Service resource. 

* `read_access_id` - (Optional) The ID of the readable Resource that will be linked to the workspace. This should be used for linking to an Automation Account resource.

* `write_access_id` - (Optional) The ID of the writable Resource that will be linked to the workspace. This should be used for linking to a Log Analytics Cluster resource.

~> **Note:** You must define at least one of the above access resource id attributes (e.g. `read_access_id` or `write_access_id`).

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Log Analytics Linked Service ID.

* `name` - The generated name of the Linked Service. The format for this attribute is always `<workspace name>/<linked service type>`(e.g. `workspace1/Automation` or `workspace1/Cluster`)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Log Analytics Workspace.
* `read` - (Defaults to 5 minutes) Used when retrieving the Log Analytics Workspace.
* `update` - (Defaults to 30 minutes) Used when updating the Log Analytics Workspace.
* `delete` - (Defaults to 30 minutes) Used when deleting the Log Analytics Workspace.

## Import

Log Analytics Workspaces can be imported using the `resource id`, e.g.

When `read_access_id` has been specified:
```shell
terraform import azurerm_log_analytics_linked_service.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/linkedServices/Automation
```
When `read_access_id` has been omitted:
```shell
terraform import azurerm_log_analytics_linked_service.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/linkedServices/Cluster
```
