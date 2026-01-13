---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_machine_power"
description: |-
  Sets the power state of an Azure Virtual Machine.
---

# Action: azurerm_virtual_machine_power

Changes the Power state of a Virtual Machine to the specified value, or restarts the Virtual Machine.

## Example Usage

```terraform
resource "azurerm_linux_virtual_machine" "example" {
  # ... Virtual Machine configuration
}

resource "azurerm_network_interface" "example" {
  name                = "example"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.test.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "terraform_data" "example" {
  input = azurerm_network_interface.example.private_ip_address

  lifecycle {
    action_trigger {
      events  = [after_update]
      actions = [action.azurerm_virtual_machine_power.example]
    }
  }
}

action "azurerm_virtual_machine_power" "example" {
  config {
    virtual_machine_id = azurerm_linux_virtual_machine.test.id
    power_action       = "restart"
  }
}
```

## Argument Reference

This action supports the following arguments:

* `virtual_machine_id` - (Required) The ID of the virtual machine on which to perform the action.

* `power_action` - (Required) The power state action to take on this virtual machine. Possible values include `restart`, `power_on`, and `power_off`.

* `timeout` - (Optional) Timeout duration to wait for the Virtual Machine Power action to complete. Defaults to `30m`.