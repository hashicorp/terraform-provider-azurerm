---
subcategory: "Disks"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_disk_pool"
description: |-
  Manages a Disk Pool.
---

# azurerm_disk_pool

Manages a Disk Pool.

!> **Note:** Azure are officially [halting](https://learn.microsoft.com/en-us/azure/azure-vmware/attach-disk-pools-to-azure-vmware-solution-hosts?tabs=azure-cli) the preview of Azure Disk Pools, and it **will not** be made generally available. New customers will not be able to register the Microsoft.StoragePool resource provider on their subscription and deploy new Disk Pools. Existing subscriptions registered with Microsoft.StoragePool may continue to deploy and manage disk pools for the time being.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
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

resource "azurerm_disk_pool" "example" {
  name                = "example-disk-pool"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku_name            = "Basic_B1"
  subnet_id           = azurerm_subnet.example.id
  zones               = ["1"]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Disk Pool. Changing this forces a new Disk Pool to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Disk Pool should exist. Changing this forces a new Disk Pool to be created.

* `location` - (Required) The Azure Region where the Disk Pool should exist. Changing this forces a new Disk Pool to be created.

* `zones` - (Required) Specifies a list of Availability Zones in which this Disk Pool should be located. Changing this forces a new Disk Pool to be created.

* `sku_name` - (Required) The SKU of the Disk Pool. Possible values are `Basic_B1`, `Standard_S1` and `Premium_P1`. Changing this forces a new Disk Pool to be created.

* `subnet_id` - (Required) The ID of the Subnet where the Disk Pool should be created. Changing this forces a new Disk Pool to be created.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the Disk Pool.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Disk Pool.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Disk Pool.
* `read` - (Defaults to 5 minutes) Used when retrieving the Disk Pool.
* `update` - (Defaults to 30 minutes) Used when updating the Disk Pool.
* `delete` - (Defaults to 30 minutes) Used when deleting the Disk Pool.

## Import

Disk Pools can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_disk_pool.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.StoragePool/diskPools/diskPool1
```
