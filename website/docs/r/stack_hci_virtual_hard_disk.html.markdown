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
resource "azurerm_stack_hci_virtual_hard_disk" "example" {
  name = "example"
  resource_group_name = "example"
  location = "West Europe"
  custom_location_id = "TODO"
}
```

## Arguments Reference

The following arguments are supported:

* `custom_location_id` - (Required) The ID of the TODO. Changing this forces a new Azure Stack HCI Virtual Hard Disk to be created.

* `location` - (Required) The Azure Region where the Azure Stack HCI Virtual Hard Disk should exist. Changing this forces a new Azure Stack HCI Virtual Hard Disk to be created.

* `name` - (Required) The name which should be used for this Azure Stack HCI Virtual Hard Disk. Changing this forces a new Azure Stack HCI Virtual Hard Disk to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Azure Stack HCI Virtual Hard Disk should exist. Changing this forces a new Azure Stack HCI Virtual Hard Disk to be created.

---

* `block_size_in_bytes` - (Optional) TODO. Changing this forces a new Azure Stack HCI Virtual Hard Disk to be created.

* `disk_file_format` - (Optional) TODO. Changing this forces a new Azure Stack HCI Virtual Hard Disk to be created.

* `disk_size_in_gb` - (Optional) TODO. Changing this forces a new Azure Stack HCI Virtual Hard Disk to be created.

* `dynamic_enabled` - (Optional) Should the TODO be enabled? Defaults to `false`. Changing this forces a new Azure Stack HCI Virtual Hard Disk to be created.

* `hyperv_generation` - (Optional) TODO. Changing this forces a new Azure Stack HCI Virtual Hard Disk to be created.

* `logical_sector_in_bytes` - (Optional) TODO. Changing this forces a new Azure Stack HCI Virtual Hard Disk to be created.

* `physical_sector_in_bytes` - (Optional) TODO. Changing this forces a new Azure Stack HCI Virtual Hard Disk to be created.

* `storage_path_id` - (Optional) The ID of the TODO. Changing this forces a new Azure Stack HCI Virtual Hard Disk to be created.

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
