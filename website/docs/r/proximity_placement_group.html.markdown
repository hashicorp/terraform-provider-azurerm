---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_proximity_placement_group"
description: |-
  Manages a proximity placement group for virtual machines, virtual machine scale sets and availability sets.

---

# azurerm_proximity_placement_group

Manages a proximity placement group for virtual machines, virtual machine scale sets and availability sets.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_proximity_placement_group" "example" {
  name                = "exampleProximityPlacementGroup"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  tags = {
    environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the proximity placement group. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the availability set. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `allowed_vm_sizes` - (Optional) Specifies the supported sizes of Virtual Machines that can be created in the Proximity Placement Group.

~> **Note:** Removing `allowed_vm_sizes` after it is set forces a new resource to be created.

* `zone` - (Optional) Specifies the supported zone of the Proximity Placement Group. Changing this forces a new resource to be created.

~> **Note:** `allowed_vm_sizes` must be set when `zone` is specified.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Proximity Placement Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Proximity Placement Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Proximity Placement Group.
* `update` - (Defaults to 30 minutes) Used when updating the Proximity Placement Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Proximity Placement Group.

## Import

Proximity Placement Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_proximity_placement_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Compute/proximityPlacementGroups/example-ppg
```
