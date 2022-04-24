---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_capacity_reservation_group"
description: |-
  Manages a Capacity Reservation Group.
---

# azurerm_capacity_reservation_group

Manages a Capacity Reservation Group.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}

resource "azurerm_capacity_reservation_group" "example" {
  name                = "example-capacity-reservation-group"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of this Capacity Reservation Group. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the resource group the Capacity Reservation Group is located in. Changing this forces a new resource to be created.

* `location` - (Required) The Azure location where the Capacity Reservation Group exists. Changing this forces a new resource to be created.

* `zones` - (Optional) Specifies a list of Availability Zones for this Capacity Reservation Group. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following Attributes are exported:

* `id` - The ID of the Capacity Reservation Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Capacity Reservation Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Capacity Reservation Group.
* `update` - (Defaults to 30 minutes) Used when updating the Capacity Reservation Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Capacity Reservation Group.

## Import

Capacity Reservation Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_capacity_reservation_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/capacityReservationGroups/capacityReservationGroup1
```
