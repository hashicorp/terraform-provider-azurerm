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
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

resource "azurerm_capacity_reservation_group" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  zones               = ["1", "2", "3"]
  tags = {
    key = "value1"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `name` - (Required) The name which should be used for this Capacity Reservation Group. Changing this forces a new Capacity Reservation Group to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Capacity Reservation Group exists.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the Capacity Reservation Group.

* `zones` - (Optional) Specifies a list of Availability Zones in which this Capacity Reservation Group should be located. Changing this forces a new Capacity Reservation Group to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Capacity Reservation Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Capacity Reservation Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Capacity Reservation Group.
* `update` - (Defaults to 30 minutes) Used when updating the Capacity Reservation Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Capacity Reservation Group.

## Import

Capacity Reservation Group can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_capacity_reservation_group.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Compute/capacityReservationGroups/capacityReservationGroup1
```
