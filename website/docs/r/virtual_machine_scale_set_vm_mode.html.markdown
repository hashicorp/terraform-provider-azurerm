---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_machine_scale_set_vm_mode"
description: |-
  Manages a Virtual Machine Scale Set VMO Mode.
---

# azurerm_windows_virtual_machine_scale_set

Manages a Virtual Machine Scale Set VMO Mode.

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

resource "azurerm_virtual_machine_scale_set_vm_mode" "eample" {
  name                = "example-VMSS"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  platform_fault_domain_count = 5

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

  virtual_machine_scale_set_id = azurerm_virtual_machine_scale_set_vm_mode.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Virtual Machine Scale Set VM Mode. Changing this forces a new resource to be created.

* `location` - (Required) The Azure location where the Virtual Machine Scale Set VM Mode should exist. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Virtual Machine Scale Set VM Mode should be exist. Changing this forces a new resource to be created.

* `platform_fault_domain_count` - (Required) Specifies the number of fault domains that are used by this Virtual Machine Scale Set VM Mode.

~> **NOTE:** The number of Fault Domains varies depending on which Azure Region you're using - [a list can be found here](https://github.com/MicrosoftDocs/azure-docs/blob/master/includes/managed-disks-common-fault-domain-region-list.md).

* `zones` - (Optional) A list of Availability Zones in which the Virtual Machines in this Scale Set should be created in. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to this Virtual Machine Scale Set.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Virtual Machine Scale Set VM Mode.

* `unique_id` - The Unique ID for this Virtual Machine Scale Set VM Mode.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Virtual Machine Scale Set VM Mode.
* `update` - (Defaults to 60 minutes) Used when updating the Virtual Machine Scale Set VM Mode.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Machine Scale Set VM Mode.
* `delete` - (Defaults to 30 minutes) Used when deleting the Virtual Machine Scale Set VM Mode.

## Import

A Virtual Machine Scale Set VM Mode can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_windows_virtual_machine_scale_set.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/Microsoft.Compute/virtualMachineScaleSets/scaleset1
```
