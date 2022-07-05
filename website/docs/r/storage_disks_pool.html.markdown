---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_disks_pool"
description: |-
  Manages a Disk Pool.
---

# azurerm_storage_disks_pool

Manages a Disk Pool.

!> **Note:** This resource has been deprecated in favour of `azurerm_disk_pool` and will be removed in version 3.0 of the Azure Provider.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-network"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_virtual_network.example.resource_group_name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.0.0/24"]
  delegation {
    name = "diskspool"
    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/read"]
      name    = "Microsoft.StoragePool/diskPools"
    }
  }
}

resource "azurerm_storage_disks_pool" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  subnet_id           = azurerm_subnet.example.id
  availability_zones  = ["1"]
  sku_name            = "Basic_B1"
  tags = {
    foo = "bar"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Disks Pool. The name must begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens, and length should be in the range [7 - 30]. Changing this forces a new Disk Pool to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Disk Pool should exist. Changing this forces a new Disk Pool to be created.

* `location` - (Required) The Azure Region where the Disks Pool should exist. Changing this forces a new Disk Pool to be created.

* `sku_name` - (Required) The SKU name of the Disk Pool. Possible values are `Basic_B1`, `Standard_S1` and `Premium_P1`. Changing this forces a new Disk Pool to be created.

* `subnet_id` - (Required) The ID of the Subnet for the Disk Pool. Changing this forces a new Disks Pool to be created.

* `zones` - (Required) Specifies a list of Availability Zones where this Disk Pool should exist. Changing this forces a new Disk Pool to be created.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the Disks Pool.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The Resource ID of the Disks Pool.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Disks Pool.
* `read` - (Defaults to 5 minutes) Used when retrieving the Disks Pool.
* `update` - (Defaults to 30 minutes) Used when updating the Disks Pool.
* `delete` - (Defaults to 30 minutes) Used when deleting the Disks Pool.

## Import

Disk Pools can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_disks_pool.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.StoragePool/diskPools/disksPool1
```
