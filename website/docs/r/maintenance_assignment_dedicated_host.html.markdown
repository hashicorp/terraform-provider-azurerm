---
subcategory: "Maintenance"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_maintenance_assignment_dedicated_host"
description: |-
  Manages a Maintenance Assignment to Dedicated Host.
---

# azurerm_maintenance_assignment_dedicated_host

Manages a maintenance assignment to Dedicated Host.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_dedicated_host_group" "example" {
  name                        = "example-host-group"
  resource_group_name         = azurerm_resource_group.example.name
  location                    = azurerm_resource_group.example.location
  platform_fault_domain_count = 2
}

resource "azurerm_dedicated_host" "example" {
  name                    = "example-host"
  location                = azurerm_resource_group.example.location
  dedicated_host_group_id = azurerm_dedicated_host_group.example.id
  sku_name                = "DSv3-Type1"
  platform_fault_domain   = 1
}

resource "azurerm_maintenance_configuration" "example" {
  name                = "example-mc"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  scope               = "All"
}

resource "azurerm_maintenance_assignment_dedicated_host" "example" {
  location                     = azurerm_resource_group.example.location
  maintenance_configuration_id = azurerm_maintenance_configuration.example.id
  dedicated_host_id            = azurerm_dedicated_host.example.id
}
```

## Argument Reference

The following arguments are supported:

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `maintenance_configuration_id` - (Required) Specifies the ID of the Maintenance Configuration Resource. Changing this forces a new resource to be created.

* `dedicated_host_id` - (Required) Specifies the Dedicated Host ID to which the Maintenance Configuration will be assigned. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Maintenance Assignment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Maintenance Assignment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Maintenance Assignment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Maintenance Assignment.

## Import

Maintenance Assignment can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_maintenance_assignment_dedicated_host.example /subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/resGroup1/providers/microsoft.compute/hostGroups/group1/hosts/host1/providers/Microsoft.Maintenance/configurationAssignments/assign1
```
