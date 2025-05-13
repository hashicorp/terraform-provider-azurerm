---
subcategory: "Desktop Virtualization"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_desktop_workspace"
description: |-
  Manages a Virtual Desktop Workspace.
---

# azurerm_virtual_desktop_workspace

Manages a Virtual Desktop Workspace.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "rg-example-virtualdesktop"
  location = "West Europe"
}

resource "azurerm_virtual_desktop_workspace" "workspace" {
  name                = "workspace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  friendly_name = "FriendlyName"
  description   = "A description of my workspace"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Virtual Desktop Workspace. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Virtual Desktop Workspace. Changing this forces a new resource to be created.

* `location` - (Required) The location/region where the Virtual Desktop Workspace is located. Changing the location/region forces a new resource to be created.

* `friendly_name` - (Optional) A friendly name for the Virtual Desktop Workspace.

* `description` - (Optional) A description for the Virtual Desktop Workspace.
  
* `public_network_access_enabled` - (Optional) Whether public network access is allowed for this Virtual Desktop Workspace. Defaults to `true`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Virtual Desktop Workspace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Virtual Desktop Workspace.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Desktop Workspace.
* `update` - (Defaults to 1 hour) Used when updating the Virtual Desktop Workspace.
* `delete` - (Defaults to 1 hour) Used when deleting the Virtual Desktop Workspace.

## Import

Virtual Desktop Workspaces can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_desktop_workspace.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/providers/Microsoft.DesktopVirtualization/workspaces/myworkspace
```
