---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_capacity_reservation"
description: |-
    Lists Capacity Reservation resources.
---

# List resource: azurerm_capacity_reservation

Lists Capacity Reservation resources.

## Example Usage

### List Capacity Reservations in a Capacity Reservation Group

```hcl
list "azurerm_capacity_reservation" "example" {
  provider = azurerm
  config {
    capacity_reservation_group_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/capacityReservationGroups/capacityReservationGroup1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `capacity_reservation_group_id` - (Required) The ID of the Capacity Reservation Group to query.
