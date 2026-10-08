---
subcategory: "Maintenance"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_maintenance_assignment_arc_machine"
description: |-
  Manages a Maintenance Assignment for an Arc Machine.
---

# azurerm_maintenance_assignment_arc_machine

Manages a Maintenance Assignment for an Arc Machine.

The assignment's location is determined automatically from the Arc Machine.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_arc_machine" "example" {
  name                = "example-am"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  kind                = "SCVMM"
}

resource "azurerm_maintenance_configuration" "example" {
  name                     = "example-mc"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  scope                    = "InGuestPatch"
  in_guest_user_patch_mode = "User"

  window {
    start_date_time = "2026-10-10 00:00"
    duration        = "02:00"
    time_zone       = "UTC"
    recur_every     = "1Day"
  }

  install_patches {
    reboot = "IfRequired"

    windows {
      classifications_to_include = ["Critical", "Security"]
    }
  }
}

resource "azurerm_maintenance_assignment_arc_machine" "example" {
  arc_machine_id               = azurerm_arc_machine.example.id
  maintenance_configuration_id = azurerm_maintenance_configuration.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `arc_machine_id` - (Required) Specifies the Arc Machine ID to which the Maintenance Configuration will be assigned. Changing this forces a new resource to be created.

* `maintenance_configuration_id` - (Required) Specifies the ID of the Maintenance Configuration Resource. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Maintenance Assignment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Maintenance Assignment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Maintenance Assignment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Maintenance Assignment.

## Import

Maintenance Assignments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_maintenance_assignment_arc_machine.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.HybridCompute/machines/machine1/providers/Microsoft.Maintenance/configurationAssignments/assign1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.HybridCompute` - 2024-07-10

* `Microsoft.Maintenance` - 2023-04-01
