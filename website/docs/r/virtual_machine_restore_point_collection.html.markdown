---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_machine_restore_point_collection"
description: |-
  Manages a Virtual Machine Restore Point Collection
---

# azurerm_virtual_machine_restore_point_collection

Manages a Virtual Machine Restore Point Collection.

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

resource "azurerm_subnet" "example" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "example" {
  name                = "example-nic"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.example.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "example" {
  name                = "example-machine"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  network_interface_ids = [
    azurerm_network_interface.example.id,
  ]

  admin_ssh_key {
    username   = "adminuser"
    public_key = file("~/.ssh/id_rsa.pub")
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }
}

resource "azurerm_virtual_machine_restore_point_collection" "example" {
  name                      = "example-collection"
  resource_group_name       = azurerm_resource_group.example.name
  location                  = azurerm_linux_virtual_machine.example.location
  source_virtual_machine_id = azurerm_linux_virtual_machine.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Virtual Machine Restore Point Collection. Changing this forces a new resource to be created.

* `location` - (Required) The Azure location where the Virtual Machine Restore Point Collection should exist. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Virtual Machine Restore Point Collection should exist. Changing this forces a new resource to be created.

* `source_virtual_machine_id` - (Required) The ID of the virtual machine that will be associated with this Virtual Machine Restore Point Collection. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to this Virtual Machine Restore Point Collection.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Virtual Machine Restore Point Collection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:


* `create` - (Defaults to 30 minutes) Used when creating the Virtual Machine Restore Point Collection.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Machine Restore Point Collection.
* `update` - (Defaults to 30 minutes) Used when updating the Virtual Machine Restore Point Collection.
* `delete` - (Defaults to 30 minutes) Used when deleting the Virtual Machine Restore Point Collection.

## Import

Virtual Machine Restore Point Collections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_machine_restore_point_collection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/restorePointCollections/collection1
```
