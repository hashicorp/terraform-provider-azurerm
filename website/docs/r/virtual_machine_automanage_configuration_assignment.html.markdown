---
subcategory: "Automanage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_machine_automanage_configuration_assignment"
description: |-
  Manages a Virtual Machine Automanage Configuration Profile Assignment.
---

# azurerm_virtual_machine_automanage_configuration_assignment

Manages a Virtual Machine Automanage Configuration Profile Assignment.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "westus"
}

resource "azurerm_virtual_network" "example" {
  name                = "examplevnet"
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
  name                = "exampleni"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.example.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "example" {
  name                            = "examplevm"
  resource_group_name             = azurerm_resource_group.example.name
  location                        = azurerm_resource_group.example.location
  size                            = "Standard_F2"
  admin_username                  = "adminuser"
  admin_password                  = "P@$$w0rd1234!"
  disable_password_authentication = false
  network_interface_ids = [
    azurerm_network_interface.example.id,
  ]

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

resource "azurerm_automanage_configuration" "example" {
  name                = "exampleconfig"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_virtual_machine_automanage_configuration_assignment" "example" {
  virtual_machine_id = azurerm_linux_virtual_machine.example.id
  configuration_id   = azurerm_automanage_configuration.example.id
}

```

## Arguments Reference

The following arguments are supported:

* `virtual_machine_id` - (Required) The ARM resource ID of the Virtual Machine to assign the Automanage Configuration to. Changing this forces a new resource to be created.

* `configuration_id` - (Required) The ARM resource ID of the Automanage Configuration to assign to the Virtual Machine. Changing this forces a new resource to be created.

---
## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Virtual Machine Automanage Configuration Profile Assignment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Automanage Configuration.
* `read` - (Defaults to 5 minutes) Used when retrieving the Automanage Configuration.
* `delete` - (Defaults to 30 minutes) Used when deleting the Automanage Configuration.

## Import

Virtual Machine Automanage Configuration Profile Assignment can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_machine_automanage_configuration_assignment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/virtualMachines/vm1/providers/Microsoft.AutoManage/configurationProfileAssignments/default
```
