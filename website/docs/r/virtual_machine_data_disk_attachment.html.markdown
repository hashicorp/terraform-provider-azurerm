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

## Example Usage

```hcl
resource "azurerm_virtual_machine_data_disk_attachment" "test" {
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of this Disk Attachment, which needs to be unique within the Virtual Machine. Changing this forces a new resource to be created.

* `virtual_machine_id` - (Required) The ID of the Virtual Machine to which the Data Disk should be attached. Changing this forces a new resource to be created.

* `create_option` - (Required) The Create Option of the Data Disk, such as `Empty` or `Attach`. Changing this forces a new resource to be created.

* `lun` - (Required) The Logical Unit Number of the Data Disk, which needs to be unique within the Virtual Machine.

* `vhd_uri` - (Optional) The URI of a Blob in a Storage Account where the VHD for this Disk should be placed. Cannot be specified when `managed_disk_id` or `managed_disk_type` is specified.

* `managed_disk_id` - (Optional) The ID of an existing Managed Disk which should be attached. When set, `create_option` should be set to `Attach`.

* `managed_disk_type` - (Optional) Specifies the type of managed disk to create. Value you must be either `Standard_LRS` or `Premium_LRS`. Cannot be used when `vhd_uri` is specified.

* `caching` - (Optional) Specifies the caching requirements for this Data Disk, such as `ReadWrite`.

* `disk_size_gb` - (Optional) Specifies the size of the Data Disk in GB.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Virtual Machine Data Disk attachment.

## Import

Virtual Machines Data Disk Attachments can be imported using the `resource id`, e.g.

```hcl
terraform import azurerm_virtual_machine_data_disk_attachment.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/microsoft.compute/virtualMachines/machine1/dataDisks/disk1
```
