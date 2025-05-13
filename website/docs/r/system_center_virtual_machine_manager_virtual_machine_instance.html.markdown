---
subcategory: "System Center Virtual Machine Manager"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_system_center_virtual_machine_manager_virtual_machine_instance"
description: |-
  Manages a System Center Virtual Machine Manager Virtual Machine Instance.
---

# azurerm_system_center_virtual_machine_manager_virtual_machine_instance

Manages a System Center Virtual Machine Manager Virtual Machine Instance.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_arc_machine" "example" {
  name                = "example-arcmachine"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  kind                = "SCVMM"
}

resource "azurerm_system_center_virtual_machine_manager_server" "example" {
  name                = "example-scvmmms"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  custom_location_id  = "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ExtendedLocation/customLocations/customLocation1"
  fqdn                = "example.labtest"
  username            = "testUser"
  password            = "H@Sh1CoR3!"
}

data "azurerm_system_center_virtual_machine_manager_inventory_items" "example" {
  inventory_type                                  = "Cloud"
  system_center_virtual_machine_manager_server_id = azurerm_system_center_virtual_machine_manager_server.example.id
}

resource "azurerm_system_center_virtual_machine_manager_cloud" "example" {
  name                                                           = "example-scvmmc"
  location                                                       = azurerm_resource_group.example.location
  resource_group_name                                            = azurerm_resource_group.example.name
  custom_location_id                                             = azurerm_system_center_virtual_machine_manager_server.example.custom_location_id
  system_center_virtual_machine_manager_server_inventory_item_id = data.azurerm_system_center_virtual_machine_manager_inventory_items.example.inventory_items[0].id
}

data "azurerm_system_center_virtual_machine_manager_inventory_items" "example2" {
  inventory_type                                  = "VirtualMachineTemplate"
  system_center_virtual_machine_manager_server_id = azurerm_system_center_virtual_machine_manager_server.example.id
}

resource "azurerm_system_center_virtual_machine_manager_virtual_machine_template" "example" {
  name                                                           = "example-scvmmvmt"
  location                                                       = azurerm_resource_group.example.location
  resource_group_name                                            = azurerm_resource_group.example.name
  custom_location_id                                             = azurerm_system_center_virtual_machine_manager_server.example.custom_location_id
  system_center_virtual_machine_manager_server_inventory_item_id = data.azurerm_system_center_virtual_machine_manager_inventory_items.example2.inventory_items[0].id
}

resource "azurerm_system_center_virtual_machine_manager_virtual_machine_instance" "example" {
  scoped_resource_id = azurerm_arc_machine.example.id
  custom_location_id = azurerm_system_center_virtual_machine_manager_server.example.custom_location_id

  infrastructure {
    system_center_virtual_machine_manager_cloud_id                  = azurerm_system_center_virtual_machine_manager_cloud.example.id
    system_center_virtual_machine_manager_template_id               = azurerm_system_center_virtual_machine_manager_virtual_machine_template.example.id
    system_center_virtual_machine_manager_virtual_machine_server_id = azurerm_system_center_virtual_machine_manager_server.example.id
  }

  operating_system {
    computer_name = "testComputer"
  }

  hardware {
    cpu_count    = 1
    memory_in_mb = 1024
  }

  lifecycle {
    // Service API always provisions a virtual disk with bus type IDE per Virtual Machine Template by default, so it has to be ignored
    ignore_changes = [storage_disk]
  }
}
```

## Arguments Reference

The following arguments are supported:

* `scoped_resource_id` - (Required) The ID of the Hybrid Compute Machine where this System Center Virtual Machine Manager Virtual Machine Instance is stored. Changing this forces a new resource to be created.

* `custom_location_id` - (Required) The ID of the Custom Location for the System Center Virtual Machine Manager Virtual Machine Instance. Changing this forces a new resource to be created.

* `infrastructure` - (Required) An `infrastructure` block as defined below.

* `hardware` - (Optional) A `hardware` block as defined below. Changing this forces a new resource to be created.

~> **Note:** This resource will be restarted while updating `hardware`.

* `network_interface` - (Optional) A `network_interface` block as defined below.

~> **Note:** This resource will be restarted while updating `network_interface`.

* `operating_system` - (Optional) An `operating_system` block as defined below. Changing this forces a new resource to be created.

* `storage_disk` - (Optional) A `storage_disk` block as defined below.

~> **Note:** This resource will be restarted while updating `storage_disk`.

* `system_center_virtual_machine_manager_availability_set_ids` - (Optional) A list of IDs of System Center Virtual Machine Manager Availability Set.

---

An `infrastructure` block supports the following:

* `checkpoint_type` - (Optional) The type of checkpoint supported for the Virtual Machine. Possible values are `Disabled`, `Production`, `ProductionOnly` and `Standard`.

* `system_center_virtual_machine_manager_cloud_id` - (Optional) The ID of the System Center Virtual Machine Manager Cloud resource to use for deploying the Virtual Machine. Changing this forces a new resource to be created.

* `system_center_virtual_machine_manager_inventory_item_id` - (Optional) The ID of the System Center Virtual Machine Manager Inventory Item for System Center Virtual Machine Manager Virtual Machine Instance. Changing this forces a new resource to be created.

* `system_center_virtual_machine_manager_template_id` - (Optional) The ID of the System Center Virtual Machine Manager Virtual Machine Template to use for deploying the Virtual Machine. Changing this forces a new resource to be created.

* `system_center_virtual_machine_manager_virtual_machine_server_id` - (Optional) The ID of the System Center Virtual Machine Manager Virtual Machine. Changing this forces a new resource to be created.

---

A `hardware` block supports the following:

* `cpu_count` - (Optional) The number of vCPUs for the Virtual Machine. Possible values are between `1` and `64`.

* `dynamic_memory_max_in_mb` - (Optional) The max dynamic memory for the Virtual Machine. Possible values are between `32` and `1048576`.

* `dynamic_memory_min_in_mb` - (Optional) The min dynamic memory for the Virtual Machine. Possible values are between `32` and `1048576`.

* `limit_cpu_for_migration_enabled` - (Optional) Whether processor compatibility mode for live migration of Virtual Machines is enabled.

* `memory_in_mb` - (Optional) The size of a Virtual Machine's memory. Possible values are between `32` and `1048576`.

---

A `network_interface` block supports the following:

* `name` - (Required) The name of the Virtual Network in System Center Virtual Machine Manager Server that the Network Interface is connected to.

* `virtual_network_id` - (Optional) The ID of the System Center Virtual Machine Manager Virtual Network to connect the Network Interface.

* `ipv4_address_type` - (Optional) The IPv4 address type. Possible values are `Dynamic` and `Static`.

* `ipv6_address_type` - (Optional) The IPv6 address type. Possible values are `Dynamic` and `Static`.

* `mac_address_type` - (Optional) The MAC address type. Possible values are `Dynamic` and `Static`.

---

An `operating_system` block supports the following:

* `computer_name` - (Optional) The computer name of the Virtual Machine. Changing this forces a new resource to be created.

* `admin_password` - (Optional) The admin password of the Virtual Machine. Changing this forces a new resource to be created.

---

A `storage_disk` block supports the following:

* `bus` - (Optional) The disk bus. Possible values are between `0` and `3`. 

* `bus_type` - (Optional) The disk bus type. Possible values are `IDE` and `SCSI`.

* `disk_size_gb` - (Optional) The disk total size.

* `lun` - (Optional) The disk lun. Possible values are between `0` and `63`.

* `name` - (Optional) The name of the disk.

* `storage_qos_policy_name` - (Optional) The name of the Storage QoS policy.

* `template_disk_id` - (Optional) The disk ID in the System Center Virtual Machine Manager Virtual Machine Template. Changing this forces a new resource to be created.

* `vhd_type` - (Optional) The disk vhd type. Possible values are `Dynamic` and `Fixed`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the System Center Virtual Machine Manager Virtual Machine Instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the System Center Virtual Machine Manager Virtual Machine Instance.
* `read` - (Defaults to 5 minutes) Used when retrieving the System Center Virtual Machine Manager Virtual Machine Instance.
* `update` - (Defaults to 1 hour) Used when updating the System Center Virtual Machine Manager Virtual Machine Instance.
* `delete` - (Defaults to 1 hour) Used when deleting the System Center Virtual Machine Manager Virtual Machine Instance.

## Import

System Center Virtual Machine Manager Virtual Machine Instances can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_system_center_virtual_machine_manager_virtual_machine_instance.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.HybridCompute/machines/machine1/providers/Microsoft.ScVmm/virtualMachineInstances/default
```
