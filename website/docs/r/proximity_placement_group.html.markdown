---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_proximity_placement_group"
sidebar_current: "docs-azurerm-resource-compute-proximity-placement-group"
description: |-
  Manages a proximity placement group for virtual machines, virtual machine scale sets and availability sets.

---

# azurerm_proximity_placement_group

Manages a proximity placement group for virtual machines, virtual machine scale sets and availability sets.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "resourceGroup1"
  location = "West US"
}

resource "azurerm_proximity_placement_group" "example" {
  name                = "exampleProximityPlacementGroup"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"

  tags = {
    environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the availability set. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the availability set. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The Proximity Placement Group ID.

* `name` - The name of the Proximity Placement Group.

* `location` - The location of the Proximity Placement Group.

* `resource_group_name` - The name of the resource group in which the Proximity Placement Group exists.

* `tags` - The tags attached to the Proximity Placement Group.

## Import

Proximity Placement Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_proximity_placement_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Compute/proximityPlacementGroups/example-ppg
```
