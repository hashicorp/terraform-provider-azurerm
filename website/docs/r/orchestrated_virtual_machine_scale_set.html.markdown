---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_orchestrated_virtual_machine_scale_set"
description: |-
  Manages an Orchestrated Virtual Machine Scale Set.
---

# azurerm_orchestrated_virtual_machine_scale_set

Manages an Orchestrated Virtual Machine Scale Set.

~> **NOTE:** This resource is part of the public preview feature of virtual machine scale set orchestration mode, which manages a virtual machine scale set in VM orchestration mode. You can find more information about the orchestration mode [here](https://docs.microsoft.com/en-us/azure/virtual-machine-scale-sets/orchestration-modes).

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-network"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "internal" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "example" {
  name                = "example-nic"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.internal.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_orchestrated_virtual_machine_scale_set" "eample" {
  name                = "example-VMSS"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  platform_fault_domain_count = 5
  single_placement_group      = true

  zones = ["1"]
}

resource "azurerm_linux_virtual_machine" "linux" {
  name                            = "example-linuxVM"
  resource_group_name             = azurerm_resource_group.main.name
  location                        = azurerm_resource_group.main.location
  size                            = "Standard_F2"
  admin_username                  = "adminuser"
  admin_password                  = "P@ssw0rd1234!"
  disable_password_authentication = false
  network_interface_ids = [
    azurerm_network_interface.example.id,
  ]

  source_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  virtual_machine_scale_set_id = azurerm_orchestrated_virtual_machine_scale_set.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Orchestrated Virtual Machine Scale Set. Changing this forces a new resource to be created.

* `location` - (Required) The Azure location where the Orchestrated Virtual Machine Scale Set should exist. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Orchestrated Virtual Machine Scale Set should exist. Changing this forces a new resource to be created.

* `platform_fault_domain_count` - (Required) Specifies the number of fault domains that are used by this Orchestrated Virtual Machine Scale Set. Changing this forces a new resource to be created.

~> **NOTE:** The number of Fault Domains varies depending on which Azure Region you're using - a list can be found [here](https://github.com/MicrosoftDocs/azure-docs/blob/master/includes/managed-disks-common-fault-domain-region-list.md).

* `single_placement_group` - (Required) Should the Orchestrated Virtual Machine Scale Set use single placement group? Changing this forces a new resource to be created.

~> **NOTE:** You cannot assign `single_placement_group` to `true` unless you have opted-in the private preview program of the orchestration mode of virtual machine scale sets.

* `zones` - (Optional) A list of Availability Zones in which the Virtual Machines in this Scale Set should be created in, currently the maximum count of availability zones is 1. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to this Orchestrated Virtual Machine Scale Set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Orchestrated Virtual Machine Scale Set.

* `unique_id` - The Unique ID for the Orchestrated Virtual Machine Scale Set.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Orchestrated Virtual Machine Scale Set.
* `update` - (Defaults to 30 minutes) Used when updating the Orchestrated Virtual Machine Scale Set.
* `read` - (Defaults to 5 minutes) Used when retrieving the Orchestrated Virtual Machine Scale Set.
* `delete` - (Defaults to 30 minutes) Used when deleting the Orchestrated Virtual Machine Scale Set.

## Import

An Orchestrated Virtual Machine Scale Set can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_orchestrated_virtual_machine_scale_set.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/Microsoft.Compute/virtualMachineScaleSets/scaleset1
```
