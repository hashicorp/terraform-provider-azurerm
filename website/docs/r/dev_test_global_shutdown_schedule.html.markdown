---
subcategory: "Dev Test"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dev_test_global_shutdown_schedule"
description: |-
    Manages automated shutdown schedules for Azure Resource Manager VMs outside of Dev Test Labs.
---

# azurerm_dev_test_schedule

Manages automated shutdown schedules for Azure Resource Manager VMs outside of Dev Test Labs.

## Example Usage

```hcl
resource "azurerm_resource_group" "sample" {
  name     = "sample-rg"
  location = "eastus"
}

resource "azurerm_virtual_network" "sample" {
  name                = "sample-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.sample.location}"
  resource_group_name = "${azurerm_resource_group.sample.name}"
}
  
resource "azurerm_subnet" "sample" {
  name                 = "sample-subnet"
  resource_group_name  = "${azurerm_resource_group.sample.name}"
  virtual_network_name = "${azurerm_virtual_network.sample.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_network_interface" "sample" {
  name                = "sample-nic"
  location            = "${azurerm_resource_group.sample.location}"
  resource_group_name = "${azurerm_resource_group.sample.name}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.sample.id}"
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurerm_virtual_machine" "sample" {
  name                  = "SampleVM"
  location              = "${azurerm_resource_group.sample.location}"
  resource_group_name   = "${azurerm_resource_group.sample.name}"
  network_interface_ids = ["${azurerm_network_interface.sample.id}"]
  vm_size               = "Standard_B2s"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "18.04-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "myosdisk-%d"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = "hostname"
    admin_username = "testadmin"
    admin_password = "Password1234!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }
}

resource "azurerm_dev_test_global_shutdown_schedule" "sample" {
  target_resource_id = "${azurerm_virtual_machine.sample.id}"
  location           = "${azurerm_resource_group.example.location}"
  status             = "Enabled"

  daily_recurrence {
    time = "1100"
  }

  time_zone_id = "Pacific Standard Time"

  notification_settings {
    status          = "Enabled"
    time_in_minutes = "60"
    webhook_url     = "https://sample-webhook-url.example.com"
  }
}

```

## Argument Reference

The following arguments are supported:

* `location` - (Required) The location where the schedule is created. Changing this forces a new resource to be created.

* `target_resource_id` - (Required) The resource ID of the target ARM-based Virtual Machine. Changing this forces a new resource to be created.

* `status` - (Optional) The status of this schedule. Possible values are `Enabled` and `Disabled`. Defaults to `Enabled`.

* `time_zone_id` - (Required) The time zone ID (e.g. Pacific Standard time). Refer to this guide for a [full list of accepted time zone names](https://jackstromberg.com/2017/01/list-of-time-zones-consumed-by-azure/).

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `daily_recurrence` - (Required) - block supports the following:

* `time` - (Required) The time each day when the schedule takes effect. Must match the format HHmm where HH is 00-23 and mm is 00-59 (e.g. 0930, 2300, etc.)

---

A `notification_settings` - (Required)  - block supports the following:

* `status` - (Optional) The status of the notification. Possible values are `Enabled` and `Disabled`. Defaults to `Disabled`

* `time_in_minutes` - Time in minutes between 15 and 120 before a shutdown event at which a notification will be sent. Required if `status` is `Enabled`. Optional otherwise.

* `webhook_url` - The webhook URL to which the notification will be sent. Required if `status` is `Enabled`. Optional otherwise.

## Attributes Reference

The following additional attributes are exported:

* `id` - The Dev Test Global Schedule ID.

## Import

An existing Dev Test Global Shutdown Schedule can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dev_test_global_shutdown_schedule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/sample-rg/providers/Microsoft.DevTestLab/schedules/compute-vm-SampleVM
```
