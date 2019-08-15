---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_disk"
sidebar_current: "docs-azurerm-resource-compute-managed-disk"
description: |-
  Manages a Managed Disk.
---

# azurerm_managed_disk

Manage a managed disk.

## Example Usage with Create Empty

```hcl
resource "azurerm_resource_group" "test" {
  name     = "acctestRG"
  location = "West US 2"
}

resource "azurerm_managed_disk" "test" {
  name                 = "acctestmd"
  location             = "West US 2"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1"

  tags = {
    environment = "staging"
  }
}
```

## Example Usage with Create Copy

```hcl
resource "azurerm_resource_group" "test" {
  name     = "acctestRG"
  location = "West US 2"
}

resource "azurerm_managed_disk" "source" {
  name                 = "acctestmd1"
  location             = "West US 2"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_managed_disk" "copy" {
  name                 = "acctestmd2"
  location             = "West US 2"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_type = "Standard_LRS"
  create_option        = "Copy"
  source_resource_id   = "${azurerm_managed_disk.source.id}"
  disk_size_gb         = "1"

  tags = {
    environment = "staging"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the managed disk. Changing this forces a
    new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create
    the managed disk.

* `location` - (Required) Specified the supported Azure location where the resource exists.
    Changing this forces a new resource to be created.

* `storage_account_type` - (Required) The type of storage to use for the managed disk.
    Allowable values are `Standard_LRS`, `Premium_LRS`, `StandardSSD_LRS` or `UltraSSD_LRS`.

* `create_option` - (Required) The method to use when creating the managed disk. Possible values include:
 * `Import` - Import a VHD file in to the managed disk (VHD specified with `source_uri`).
 * `Empty` - Create an empty managed disk.
 * `Copy` - Copy an existing managed disk or snapshot (specified with `source_resource_id`).
 * `FromImage` - Copy a Platform Image (specified with `image_reference_id`)
 * `Restore` - Set by Azure Backup or Site Recovery on a restored disk (specified with `source_resource_id`).

* `source_uri` - (Optional) URI to a valid VHD file to be used when `create_option` is `Import`.

* `source_resource_id` - (Optional) ID of an existing managed disk to copy `create_option` is `Copy`
    or the recovery point to restore when `create_option` is `Restore`

* `image_reference_id` - (Optional) ID of an existing platform/marketplace disk image to copy when `create_option` is `FromImage`.

* `os_type` - (Optional) Specify a value when the source of an `Import` or `Copy`
    operation targets a source that contains an operating system. Valid values are `Linux` or `Windows`

* `disk_size_gb` - (Optional, Required for a new managed disk) Specifies the size of the managed disk to create in gigabytes.
    If `create_option` is `Copy` or `FromImage`, then the value must be equal to or greater than the source's size.

* `disk_iops_read_write` - (Optional) The number of IOPS allowed for this disk; only settable for UltraSSD disks.

* `disk_mbps_read_write` - (Optional) The bandwidth allowed for this disk; only settable for UltraSSD disks.

* `encryption_settings` - (Optional) an `encryption_settings` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `zones` - (Optional) A collection containing the availability zone to allocate the Managed Disk in.

-> **Please Note**: Availability Zones are [only supported in several regions at this time](https://docs.microsoft.com/en-us/azure/availability-zones/az-overview).

For more information on managed disks, such as sizing options and pricing, please check out the
[azure documentation](https://docs.microsoft.com/en-us/azure/storage/storage-managed-disks-overview).

---

`encryption_settings` supports:

* `enabled` - (Required) Is Encryption enabled on this Managed Disk? Changing this forces a new resource to be created.
* `disk_encryption_key` - (Optional) A `disk_encryption_key` block as defined below.
* `key_encryption_key` - (Optional) A `key_encryption_key` block as defined below.

`disk_encryption_key` supports:

* `secret_url` - (Required) The URL to the Key Vault Secret used as the Disk Encryption Key. This can be found as `id` on the `azurerm_key_vault_secret` resource.

* `source_vault_id` - (Required) The URL of the Key Vault. This can be found as `vault_uri` on the `azurerm_key_vault` resource.

`key_encryption_key` supports:

* `key_url` - (Required) The URL to the Key Vault Key used as the Key Encryption Key. This can be found as `id` on the `azurerm_key_vault_secret` resource.

* `source_vault_id` - (Required) The URL of the Key Vault. This can be found as `vault_uri` on the `azurerm_key_vault` resource.


## Attributes Reference

The following attributes are exported:

* `id` - The managed disk ID.

## Import

Managed Disks can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_managed_disk.test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/microsoft.compute/disks/manageddisk1
```
