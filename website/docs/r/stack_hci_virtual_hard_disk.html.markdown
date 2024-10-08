---
subcategory: "Azure Stack HCI"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stack_hci_virtual_hard_disk"
description: |-
  Manages an Azure Stack HCI Virtual Hard Disk.
---

# azurerm_stack_hci_virtual_hard_disk

Manages an Azure Stack HCI Virtual Hard Disk.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}

resource "azurerm_stack_hci_storage_path" "example" {
  name                = "example-sp"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  custom_location_id  = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ExtendedLocation/customLocations/cl1"
  path                = "C:\\ClusterStorage\\UserStorage_2\\sp-example"
  tags = {
    foo = "bar"
  }
}

resource "azurerm_stack_hci_virtual_hard_disk" "example" {
  name                = "example-vhd"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  custom_location_id  = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ExtendedLocation/customLocations/cl1"
  disk_size_in_gb     = 2
  storage_path_id     = azurerm_stack_hci_storage_path.example.id
  tags = {
    foo = "bar"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Azure Stack HCI Virtual Hard Disk. Changing this forces a new Azure Stack HCI Virtual Hard Disk to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Azure Stack HCI Virtual Hard Disk should exist. Changing this forces a new Azure Stack HCI Virtual Hard Disk to be created.

* `location` - (Required) The Azure Region where the Azure Stack HCI Virtual Hard Disk should exist. Changing this forces a new Azure Stack HCI Virtual Hard Disk to be created.

* `custom_location_id` - (Required) The ID of the Custom Location where the Azure Stack HCI Virtual Hard Disk should exist. Changing this forces a new Azure Stack HCI Virtual Hard Disk to be created.

* `disk_size_in_gb` - (Required) The size of the disk in GB. Changing this forces a new Azure Stack HCI Virtual Hard Disk to be created.

---

* `block_size_in_bytes` - (Optional) The block size of the disk in bytes. Changing this forces a new Azure Stack HCI Virtual Hard Disk to be created.

* `disk_file_format` - (Optional) The format of the disk file. Possible values are `vhdx` and `vhd`. Changing this forces a new Azure Stack HCI Virtual Hard Disk to be created.

* `dynamic_enabled` - (Optional) Whether to enable dynamic sizing for the Azure Stack HCI Virtual Hard Disk. Defaults to `false`. Changing this forces a new Azure Stack HCI Virtual Hard Disk to be created.

* `hyperv_generation` - (Optional) The hypervisor generation of the Azure Stack HCI Virtual Hard Disk. Possible values are `V1` and `V2`. Changing this forces a new Azure Stack HCI Virtual Hard Disk to be created.

* `logical_sector_in_bytes` - (Optional) The logical sector size of the disk in bytes. Changing this forces a new Azure Stack HCI Virtual Hard Disk to be created.

* `physical_sector_in_bytes` - (Optional) The physical sector size of the disk in bytes. Changing this forces a new Azure Stack HCI Virtual Hard Disk to be created.

* `storage_path_id` - (Optional) The ID of the Azure Stack HCI Storage Path used for this Virtual Hard Disk. Changing this forces a new Azure Stack HCI Virtual Hard Disk to be created.

-> **Note:** If `storage_path_id` is not specified, the Virtual Hard Disk will be placed in a high availability Storage Path. If you experience a diff you may need to add this to `ignore_changes`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Azure Stack HCI Virtual Hard Disk.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Azure Stack HCI Virtual Hard Disk.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure Stack HCI Virtual Hard Disk.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Stack HCI Virtual Hard Disk.
* `update` - (Defaults to 30 minutes) Used when updating the Azure Stack HCI Virtual Hard Disk.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Stack HCI Virtual Hard Disk.

## Import

Azure Stack HCI Virtual Hard Disks can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stack_hci_virtual_hard_disk.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.AzureStackHCI/virtualHardDisks/disk1
```
