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
  resource_group_name     = azurerm_resource_group.example.name
  workspace_name          = azurerm_log_analytics_workspace.example.name
  read_access_resource_id = azurerm_automation_account.example.id
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group in which the Log Analytics Linked Service is created. Changing this forces a new resource to be created.

* `workspace_name` - (Required) Name of the Log Analytics Workspace that will contain the Log Analytics Linked Service resource. Changing this forces a new resource to be created.

* `linked_service_type` - (Optional) The resource type that is specified by the `workspace_name` attribute. Supported values are `automation` and `cluster`. Defaults to `automation`. Changing this forces a new resource to be created.

* `read_access_resource_id` - (Optional) The ID of the Resource that will be linked to the workspace. This should be used for linking resources which only require read access.

* `write_access_resource_id` - (Optional) The ID of the Resource that will be linked to the workspace. This should be used for linking resources which require write access.

~> **NOTE:** You must define at least one of the above access resource id attributes (e.g. `read_access_resource_id` or `write_access_resource_id`).

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The Log Analytics Linked Service ID.

* `name` - The generated name of the Linked Service. The format for this attribute is always `<workspace_name>/<linked_service_type>`(e.g. `workspace1/Automation` or `workspace1/Cluster`)

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
