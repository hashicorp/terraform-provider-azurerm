---
subcategory: "Policy"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_machine_guest_configuration_assignment"
description: |-
  Manages a Virtual Machine Guest Configuration Assignment.
---

# azurerm_virtual_machine_guest_configuration_assignment

Manages a Virtual Machine Guest Configuration Assignment.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-gca"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "example" {
  name                = "example-nic"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.example.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "example" {
  name                = "example-vm"
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
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }
}

resource "azurerm_virtual_machine_guest_configuration_assignment" "example" {
  name               = "example-gca"
  location           = azurerm_linux_virtual_machine.example.id
  virtual_machine_id = azurerm_linux_virtual_machine.example.id
  guest_configuration {
    name    = "example-assignment"
    version = "1.0"

    parameter {
      name  = "[InstalledApplication]bwhitelistedapp;Name"
      value = "NotePad,sql"
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Virtual Machine Guest Configuration Assignment. Changing this forces a new resource to be created.

* `location` - (Required) The Azure location where the Virtual Machine Guest Configuration Assignment should exist. Changing this forces a new resource to be created.

* `virtual_machine_id` - (Required) The resource ID of the Virtual Machine which this Guest Configuration Assignment should apply to. Changing this forces a new resource to be created.

---

* `guest_configuration` - (Optional)  A `guest_configuration` block as defined below.

---

An `guest_configuration` block supports the following:

* `name` - (Optional) The name which should be used for this Guest Configuration.

* `parameter` - (Optional)  A `parameter` block as defined below.

* `version` - (Optional) The version of this Guest Configuration.

---

An `parameter` block supports the following:

* `name` - (Required) The name which should be used for this configuration_parameter.

* `value` - (Required) The value of the configuration parameter.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Virtual Machine Guest Configuration Assignment.

* `compliance_status` - A value indicating compliance status of the machine for the assigned Guest Configuration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Virtual Machine Guest Configuration Assignment.
* `update` - (Defaults to 30 minutes) Used when updating the Virtual Machine Guest Configuration Assignment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Virtual Machine Guest Configuration Assignment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Virtual Machine Guest Configuration Assignment.

## Import

Virtual Machine Guest Configuration Assignments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_machine_guest_configuration_assignment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/virtualMachines/vm1/providers/Microsoft.GuestConfiguration/guestConfigurationAssignments/assignment1
```
