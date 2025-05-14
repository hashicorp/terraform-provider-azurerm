---
subcategory: "Policy"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_policy_virtual_machine_configuration_assignment"
description: |-
  Applies a Guest Configuration Policy to a Virtual Machine.
---

# azurerm_policy_virtual_machine_configuration_assignment

Applies a Guest Configuration Policy to a Virtual Machine.

~> **Note:** You can create Guest Configuration Policies without defining a `azurerm_virtual_machine_extension` resource, however the policies will not be executed until a `azurerm_virtual_machine_extension` has been provisioned to the virtual machine.

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

resource "azurerm_windows_virtual_machine" "example" {
  name                = "examplevm"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  size                = "Standard_F2"
  admin_username      = "adminuser"
  admin_password      = "P@$$w0rd1234!"
  network_interface_ids = [
    azurerm_network_interface.example.id,
  ]

  identity {
    type = "SystemAssigned"
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2019-Datacenter"
    version   = "latest"
  }
}

resource "azurerm_virtual_machine_extension" "example" {
  name                       = "AzurePolicyforWindows"
  virtual_machine_id         = azurerm_windows_virtual_machine.example.id
  publisher                  = "Microsoft.GuestConfiguration"
  type                       = "ConfigurationforWindows"
  type_handler_version       = "1.29"
  auto_upgrade_minor_version = "true"
}

resource "azurerm_policy_virtual_machine_configuration_assignment" "example" {
  name               = "AzureWindowsBaseline"
  location           = azurerm_windows_virtual_machine.example.location
  virtual_machine_id = azurerm_windows_virtual_machine.example.id

  configuration {
    assignment_type = "ApplyAndMonitor"
    version         = "1.*"

    parameter {
      name  = "Minimum Password Length;ExpectedValue"
      value = "16"
    }
    parameter {
      name  = "Minimum Password Age;ExpectedValue"
      value = "0"
    }
    parameter {
      name  = "Maximum Password Age;ExpectedValue"
      value = "30,45"
    }
    parameter {
      name  = "Enforce Password History;ExpectedValue"
      value = "10"
    }
    parameter {
      name  = "Password Must Meet Complexity Requirements;ExpectedValue"
      value = "1"
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Guest Configuration that will be assigned in this Guest Configuration Assignment. Changing this forces a new resource to be created.

* `location` - (Required) The Azure location where the Policy Virtual Machine Configuration Assignment should exist. Changing this forces a new resource to be created.

* `virtual_machine_id` - (Required) The resource ID of the Policy Virtual Machine which this Guest Configuration Assignment should apply to. Changing this forces a new resource to be created.

* `configuration` - (Required) A `configuration` block as defined below.

---

A `configuration` block supports the following:

* `assignment_type` - (Optional) The assignment type for the Guest Configuration Assignment. Possible values are `Audit`, `ApplyAndAutoCorrect`, `ApplyAndMonitor` and `DeployAndAutoCorrect`.

* `content_hash` - (Optional) The content hash for the Guest Configuration package.

* `content_uri` - (Optional) The content URI where the Guest Configuration package is stored.

~> **Note:** When deploying a Custom Guest Configuration package the `content_hash` and `content_uri` fields must be defined. For Built-in Guest Configuration packages, such as the `AzureWindowsBaseline` package, the `content_hash` and `content_uri` should not be defined, rather these fields will be returned after the Built-in Guest Configuration package has been provisioned. For more information on guest configuration assignments please see the [product documentation](https://docs.microsoft.com/azure/governance/policy/concepts/guest-configuration-assignments).

* `parameter` - (Optional) One or more `parameter` blocks as defined below which define what configuration parameters and values against.

* `version` - (Optional) The version of the Guest Configuration that will be assigned in this Guest Configuration Assignment.

---

A `parameter` block supports the following:

* `name` - (Required) The name of the configuration parameter to check.

* `value` - (Required) The value to check the configuration parameter with.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Policy Virtual Machine Configuration Assignment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Policy Virtual Machine Configuration Assignment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Policy Virtual Machine Configuration Assignment.
* `update` - (Defaults to 30 minutes) Used when updating the Policy Virtual Machine Configuration Assignment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Policy Virtual Machine Configuration Assignment.

## Import

Policy Virtual Machine Configuration Assignments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_policy_virtual_machine_configuration_assignment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/virtualMachines/vm1/providers/Microsoft.GuestConfiguration/guestConfigurationAssignments/assignment1
```
