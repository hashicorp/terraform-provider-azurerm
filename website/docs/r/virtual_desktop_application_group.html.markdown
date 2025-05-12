---
subcategory: "Desktop Virtualization"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_desktop_application_group"
description: |-
  Manages a Virtual Desktop Application Group.
---

# azurerm_virtual_desktop_application_group

Manages a Virtual Desktop Application Group.

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

  type               = "Pooled"
  load_balancer_type = "BreadthFirst"
}

resource "azurerm_virtual_desktop_host_pool" "personalautomatic" {
  name                = "personalautomatic"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  type                             = "Personal"
  personal_desktop_assignment_type = "Automatic"
  load_balancer_type               = "BreadthFirst"
}

resource "azurerm_virtual_desktop_application_group" "remoteapp" {
  name                = "acctag"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  type          = "RemoteApp"
  host_pool_id  = azurerm_virtual_desktop_host_pool.pooledbreadthfirst.id
  friendly_name = "TestAppGroup"
  description   = "Acceptance Test: An application group"
}

resource "azurerm_virtual_desktop_application_group" "desktopapp" {
  name                = "appgroupdesktop"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  type          = "Desktop"
  host_pool_id  = azurerm_virtual_desktop_host_pool.personalautomatic.id
  friendly_name = "TestAppGroup"
  description   = "Acceptance Test: An application group"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Virtual Desktop Application Group. Changing the name forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Virtual Desktop Application Group. Changing this forces a new resource to be created.

* `location` - (Required) The location/region where the Virtual Desktop Application Group is located. Changing this forces a new resource to be created.

* `type` - (Required) Type of Virtual Desktop Application Group. Valid options are `RemoteApp` or `Desktop` application groups. Changing this forces a new resource to be created.

* `host_pool_id` - (Required) Resource ID for a Virtual Desktop Host Pool to associate with the Virtual Desktop Application Group. Changing the name forces a new resource to be created.

* `friendly_name` - (Optional) Option to set a friendly name for the Virtual Desktop Application Group.

* `default_desktop_display_name` - (Optional) Option to set the display name for the default sessionDesktop desktop when `type` is set to `Desktop`. A value here is mandatory for connections to the desktop using the Windows 365 portal. Without it the connection will hang at 'Loading Client'.

* `description` - (Optional) Option to set a description for the Virtual Desktop Application Group.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Virtual Desktop Application Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Virtual Desktop Application Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Desktop Application Group.
* `update` - (Defaults to 1 hour) Used when updating the Virtual Desktop Application Group.
* `delete` - (Defaults to 1 hour) Used when deleting the Virtual Desktop Application Group.

## Import

Virtual Desktop Application Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_desktop_application_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/providers/Microsoft.DesktopVirtualization/applicationGroups/myapplicationgroup
```
