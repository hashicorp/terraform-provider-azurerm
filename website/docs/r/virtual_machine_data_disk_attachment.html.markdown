---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_machine_data_disk_attachment"
sidebar_current: "docs-azurerm-resource-compute-virtual-machine-data-disk-attachment"
description: |-
  Manages attaching a Disk to a Virtual Machine.
---

# azurerm_virtual_machine_data_disk_attachment

Manages attaching a Disk to a Virtual Machine.

~> **NOTE:** Data Disks can be attached either directly on the `azurerm_virtual_machine` resource, or using the `azurerm_virtual_machine_data_disk_attachment` resource - but the two cannot be used together. If both are used against the same Virtual Machine, spurious changes will occur.

-> **Please Note:** only Managed Disks are supported via this separate resource, Unmanaged Disks can be attached using the `storage_data_disk` block in the `azurerm_virtual_machine` resource.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  # ...
}

resource "azurerm_virtual_network" "example" {
  # ...
}

resource "azurerm_subnet" "example" {
  # ...
}

resource "azurerm_network_interface" "example" {
  # ...
}

resource "azurerm_virtual_machine" "example" {
  # ...
}

resource "azurerm_managed_disk" "example" {
  # ...
}

resource "azurerm_virtual_machine_data_disk_attachment" "example" {
  managed_disk_id    = "${azurerm_managed_disk.example.id}"
  virtual_machine_id = "${azurerm_virtual_machine.example.id}"
  lun                = "10"
  caching            = "ReadWrite"
}
```

## Argument Reference

The following arguments are supported:

* `virtual_machine_id` - (Required) The ID of the Virtual Machine to which the Data Disk should be attached. Changing this forces a new resource to be created.

* `managed_disk_id` - (Required) The ID of an existing Managed Disk which should be attached. Changing this forces a new resource to be created.

* `lun` - (Required) The Logical Unit Number of the Data Disk, which needs to be unique within the Virtual Machine. Changing this forces a new resource to be created.

* `caching` - (Required) Specifies the caching requirements for this Data Disk. Possible values include `None`, `ReadOnly` and `ReadWrite`.

* `create_option` - (Optional) The Create Option of the Data Disk, such as `Empty` or `Attach`. Defaults to `Attach`. Changing this forces a new resource to be created.

* `write_accelerator_enabled` - (Optional) Specifies if Write Accelerator is enabled on the disk. This can only be enabled on `Premium_LRS` managed disks with no caching and [M-Series VMs](https://docs.microsoft.com/en-us/azure/virtual-machines/workloads/sap/how-to-enable-write-accelerator). Defaults to `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Virtual Machine Data Disk attachment.

## Import

Virtual Machines Data Disk Attachments can be imported using the `resource id`, e.g.

```hcl
terraform import azurerm_virtual_machine_data_disk_attachment.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/microsoft.compute/virtualMachines/machine1/dataDisks/disk1
```

-> **Please Note:** This is a Terraform Unique ID matching the format: `{virtualMachineID}/dataDisks/{diskName}`
