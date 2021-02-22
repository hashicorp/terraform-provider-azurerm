---
subcategory: "Policy"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_guest_configuration_assignment"
description: |-
  Manages a Guest Configuration Assignment.
---

# azurerm_guest_configuration_assignment

Manages a Guest Configuration Assignment.

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

resource "azurerm_guest_configuration_assignment" "example" {
  name               = "example-gca"
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

* `name` - (Required) The name of the Guest Configuration Assignment. Changing this forces a new resource to be created.

* `location` - (Required) The Azure location where the Guest Configuration Assignment should exist. Changing this forces a new resource to be created.

* `virtual_machine_id` - (Required) The resource ID of the Virtual Machine which this Guest Configuration should apply to. Changing this forces a new resource to be created.

---

* `context` - (Optional) The source which initiated the guest configuration assignment. Ex: Azure Policy.

* `guest_configuration` - (Optional)  A `guest_configuration` block as defined below.

---

An `guest_configuration` block supports the following:

* `name` - (Optional) The name which should be used for this guest_configuration.

* `parameter` - (Optional)  A `parameter` block as defined below.

* `version` - (Optional) The version of this Guest Configuration.

---

An `parameter` block supports the following:

* `name` - (Required) The name which should be used for this configuration_parameter.

* `value` - (Required) The value of the configuration parameter.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Guest Configuration Assignment.

* `assignment_hash` - Combined hash of the configuration package and parameters.

* `compliance_status` - A value indicating compliance status of the machine for the assigned guest configuration.

* `last_compliance_status_checked` - Date and time when last compliance status was checked.

* `latest_assignment_report` - Last reported guest configuration assignment report. A `latest_assignment_report` block as defined below.

* `latest_report_id` - The ID of the latest_report.

* `target_resource_id` - The ID of the target_resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the guestconfiguration Assignment.
* `read` - (Defaults to 5 minutes) Used when retrieving the guestconfiguration Assignment.
* `delete` - (Defaults to 30 minutes) Used when deleting the guestconfiguration Assignment.

## Import

guestconfiguration Assignments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_guest_configuration_assignment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/virtualMachines/vm1/providers/Microsoft.GuestConfiguration/guestConfigurationAssignments/assignment1
```
