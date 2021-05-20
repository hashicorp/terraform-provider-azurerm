---
subcategory: "Desktop Virtualization"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_desktop_workspace_application_group_association"
description: |-
  Manages a Virtual Desktop Workspace Application Group Association.
---

# azurerm_virtual_desktop_workspace_application_group_association

Manages a Virtual Desktop Workspace Application Group Association.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "rg-example-virtualdesktop"
  location = "West Europe"
}

resource "azurerm_virtual_desktop_host_pool" "pooledbreadthfirst" {
  name                = "pooledbreadthfirst"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  type                = "Pooled"
  load_balancer_type  = "BreadthFirst"
}

resource "azurerm_virtual_desktop_application_group" "remoteapp" {
  name                = "remoteapp"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  type                = "RemoteApp"
  host_pool_id        = azurerm_virtual_desktop_host_pool.pooledbreadthfirst.id
}

resource "azurerm_virtual_desktop_workspace" "workspace" {
  name                = "workspace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_virtual_desktop_workspace_application_group_association" "workspaceremoteapp" {
  workspace_id         = azurerm_virtual_desktop_workspace.workspace.id
  application_group_id = azurerm_virtual_desktop_application_group.remoteapp.id
}
```

## Argument Reference

The following arguments are supported:

* `workspace_id` - (Required) The resource ID for the Virtual Desktop Workspace.

* `application_group_id` - (Required) The resource ID for the Virtual Desktop Application Group.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Virtual Desktop Workspace Application Group association.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Virtual Desktop Workspace.
* `update` - (Defaults to 60 minutes) Used when updating the Virtual Desktop Workspace.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Desktop Workspace.
* `delete` - (Defaults to 60 minutes) Used when deleting the Virtual Desktop Workspace.

## Import

Associations between Virtual Desktop Workspaces and Virtual Desktop Application Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_desktop_workspace_application_group_association.association1 "/subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/myGroup1/providers/Microsoft.DesktopVirtualization/workspaces/myworkspace|/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/providers/Microsoft.DesktopVirtualization/applicationGroups/myapplicationgroup"
```

-> **NOTE:** This ID is specific to Terraform - and is of the format `{virtualDesktopWorkspaceID}|{virtualDesktopApplicationGroupID}`.
