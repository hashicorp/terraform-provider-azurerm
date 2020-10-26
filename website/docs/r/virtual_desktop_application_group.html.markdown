---
subcategory: "DesktopVirtualization"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_desktop_application_group"
description: |-
  Manages a Virtual Desktop Application Group.
---

# virtual_desktop_application_group

Manages a Virtual Desktop Application Group.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "rg-example-virtualdesktop"
  location = "eastus"
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

* `resource_group_name` - (Required) The name of the resource group in which to
    create the Virtual Desktop Application Group. Changing the resource group name forces
    a new resource to be created.

* `location` - (Required) The location/region where the Virtual Desktop Application Group is
    located. Changing the location/region forces a new resource to be created.

* `type` - (Required) Type of Virtual Desktop Application Group.
    Valid options are `RemoteApp` or `Desktop` application groups.

* `host_pool_id` - (Required) Resource ID for a Virtual Desktop Host Pool to associate with the
    Virtual Desktop Application Group.

* `friendly_name` - (Optional) Option to set a friendly name for the Virtual Desktop Application Group.

* `description` - (Optional) Option to set a description for the Virtual Desktop Application Group.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Virtual Desktop Application Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Virtual Desktop Application Group.
* `update` - (Defaults to 60 minutes) Used when updating the Virtual Desktop Application Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Desktop Application Group.
* `delete` - (Defaults to 60 minutes) Used when deleting the Virtual Desktop Application Group.

## Import

Virtual Desktop Application Groups can be imported using the `resource id`, e.g.

```
terraform import virtual_desktop_application_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/providers/Microsoft.DesktopVirtualization/applicationGroups/myapplicationgroup
```
