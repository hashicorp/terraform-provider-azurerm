---
subcategory: "Dev Test"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dev_test_global_vm_shutdown_schedule"
description: |-
    Manages automated shutdown schedules for Azure Resource Manager VMs outside of Dev Test Labs.
---

# azurerm_dev_test_global_vm_shutdown_schedule

Manages automated shutdown schedules for Azure VMs that are not within an Azure DevTest Lab. While this is part of the DevTest Labs service in Azure,
this resource applies only to standard VMs, not DevTest Lab VMs. To manage automated shutdown schedules for DevTest Lab VMs, reference the
[`azurerm_dev_test_schedule` resource](dev_test_schedule.html)

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "sample-rg"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "sample-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "sample-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_network_interface" "example" {
  name                = "sample-nic"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = azurerm_subnet.example.id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_linux_virtual_machine" "example" {
  name                  = "SampleVM"
  location              = azurerm_resource_group.example.location
  resource_group_name   = azurerm_resource_group.example.name
  network_interface_ids = [azurerm_network_interface.example.id]
  size                  = "Standard_B2s"

  source_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-jammy"
    sku       = "22_04-lts"
    version   = "latest"
  }

  os_disk {
    name                 = "myosdisk-example"
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  admin_username                  = "testadmin"
  admin_password                  = "Password1234!"
  disable_password_authentication = false
}

resource "azurerm_dev_test_global_vm_shutdown_schedule" "example" {
  virtual_machine_id = azurerm_linux_virtual_machine.example.id
  location           = azurerm_resource_group.example.location
  enabled            = true

  daily_recurrence_time = "1100"
  timezone              = "Pacific Standard Time"

  notification_settings {
    enabled         = true
    time_in_minutes = "60"
    webhook_url     = "https://sample-webhook-url.example.com"
  }
}
```

## Argument Reference

The following arguments are supported:

* `location` - (Required) The location where the schedule is created. Changing this forces a new resource to be created.

* `virtual_machine_id` - (Required) The resource ID of the target ARM-based Virtual Machine. Changing this forces a new resource to be created.

* `enabled` - (Optional) Whether to enable the schedule. Possible values are `true` and `false`. Defaults to `true`.

* `timezone` - (Required) The time zone ID (e.g. Pacific Standard time). Refer to this guide for a [full list of accepted time zone names](https://jackstromberg.com/2017/01/list-of-time-zones-consumed-by-azure/).

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `daily_recurrence_time` - (Required) The time each day when the schedule takes effect. Must match the format HHmm where HH is 00-23 and mm is 00-59 (e.g. 0930, 2300, etc.)

* `notification_settings` - (Required) The notification setting of a schedule. A `notification_settings` block as defined below.

---

A `notification_settings` block supports the following:

* `enabled` - (Required) Whether to enable pre-shutdown notifications. Possible values are `true` and `false`.

* `email` - (Optional) E-mail address to which the notification will be sent.

* `time_in_minutes` - (Optional) Time in minutes between 15 and 120 before a shutdown event at which a notification will be sent. Defaults to `30`.

* `webhook_url` - (Optional) The webhook URL to which the notification will be sent.

## Attributes Reference

The following additional attributes are exported:

* `id` - The Dev Test Global Schedule ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Dev Test Global Vm Shutdown Schedule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Dev Test Global Vm Shutdown Schedule.
* `update` - (Defaults to 30 minutes) Used when updating the Dev Test Global Vm Shutdown Schedule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Dev Test Global Vm Shutdown Schedule.

## Import

An existing Dev Test Global Shutdown Schedule can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dev_test_global_vm_shutdown_schedule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/sample-rg/providers/Microsoft.DevTestLab/schedules/shutdown-computevm-SampleVM
```

The name of the resource within the `resource id` will always follow the format `shutdown-computevm-<VM Name>` where `<VM Name>` is replaced by the name of the target Virtual Machine
