---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_capacity_reservation"
description: |-
  Manages a Capacity Reservation within a Capacity Reservation Group.
---

# azurerm_capacity_reservation

Manages a Capacity Reservation within a Capacity Reservation Group.

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

resource "azurerm_capacity_reservation" "example" {
  name                          = "example-capacity-reservation"
  capacity_reservation_group_id = azurerm_capacity_reservation_group.example.id
  sku {
    name     = "Standard_D2s_v3"
    capacity = 1
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of this Capacity Reservation. Changing this forces a new resource to be created.

* `capacity_reservation_group_id` - (Required) The ID of the Capacity Reservation Group where the Capacity Reservation exists. Changing this forces a new resource to be created.

* `sku` - (Required) A `sku` block as defined below.

* `zone` - (Optional) Specifies the Availability Zone for this Capacity Reservation. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

An `sku` block exports the following:

* `name` - (Required) Name of the sku, such as `Standard_F2`. Changing this forces a new resource to be created.

* `capacity` - (Required) Specifies the number of instances to be reserved. It must be greater than or equal to `0` and not exceed the quota in the subscription.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Capacity Reservation.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Capacity Reservation.
* `read` - (Defaults to 5 minutes) Used when retrieving the Capacity Reservation.
* `update` - (Defaults to 30 minutes) Used when updating the Capacity Reservation.
* `delete` - (Defaults to 30 minutes) Used when deleting the Capacity Reservation.

## Import

Capacity Reservations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_capacity_reservation.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/capacityReservationGroups/capacityReservationGroup1/capacityReservations/capacityReservation1
```
