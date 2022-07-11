---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_disk"
description: |-
  Manages a Managed Disk.
---

# azurerm_managed_disk

Manages a managed disk.

## Example Usage with Create Empty

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_managed_disk" "example" {
  name                 = "acctestmd"
  location             = azurerm_resource_group.example.location
  resource_group_name  = azurerm_resource_group.example.name
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
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_managed_disk" "source" {
  name                 = "acctestmd1"
  location             = azurerm_resource_group.example.location
  resource_group_name  = azurerm_resource_group.example.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_managed_disk" "copy" {
  name                 = "acctestmd2"
  location             = azurerm_resource_group.example.location
  resource_group_name  = azurerm_resource_group.example.name
  storage_account_type = "Standard_LRS"
  create_option        = "Copy"
  source_resource_id   = azurerm_managed_disk.source.id
  disk_size_gb         = "1"

  tags = {
    environment = "staging"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Managed Disk. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Managed Disk should exist.

* `location` - (Required) Specified the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `storage_account_type` - (Required) The type of storage to use for the managed disk. Possible values are `Standard_LRS`, `StandardSSD_ZRS`, `Premium_LRS`, `Premium_ZRS`, `StandardSSD_LRS` or `UltraSSD_LRS`.

-> **Note:** Azure Ultra Disk Storage is only available in a region that support availability zones and can only enabled on the following VM series: `ESv3`, `DSv3`, `FSv3`, `LSv2`, `M` and `Mv2`. For more information see the `Azure Ultra Disk Storage` [product documentation](https://docs.microsoft.com/azure/virtual-machines/windows/disks-enable-ultra-ssd).

* `create_option` - (Required) The method to use when creating the managed disk. Changing this forces a new resource to be created. Possible values include:
 * `Import` - Import a VHD file in to the managed disk (VHD specified with `source_uri`).
 * `Empty` - Create an empty managed disk.
 * `Copy` - Copy an existing managed disk or snapshot (specified with `source_resource_id`).
 * `FromImage` - Copy a Platform Image (specified with `image_reference_id`)
 * `Restore` - Set by Azure Backup or Site Recovery on a restored disk (specified with `source_resource_id`).

---

* `disk_encryption_set_id` - (Optional) The ID of a Disk Encryption Set which should be used to encrypt this Managed Disk. Conflicts with `secure_vm_disk_encryption_set_id`.

~> **NOTE:** The Disk Encryption Set must have the `Reader` Role Assignment scoped on the Key Vault - in addition to an Access Policy to the Key Vault

~> **NOTE:** Disk Encryption Sets are in Public Preview in a limited set of regions

* `disk_iops_read_write` - (Optional) The number of IOPS allowed for this disk; only settable for UltraSSD disks. One operation can transfer between 4k and 256k bytes.

* `disk_mbps_read_write` - (Optional) The bandwidth allowed for this disk; only settable for UltraSSD disks. MBps means millions of bytes per second.

* `disk_iops_read_only` - (Optional) The number of IOPS allowed across all VMs mounting the shared disk as read-only; only settable for UltraSSD disks with shared disk enabled. One operation can transfer between 4k and 256k bytes.

* `disk_mbps_read_only` - (Optional) The bandwidth allowed across all VMs mounting the shared disk as read-only; only settable for UltraSSD disks with shared disk enabled. MBps means millions of bytes per second.

* `disk_size_gb` - (Optional, Required for a new managed disk) Specifies the size of the managed disk to create in gigabytes. If `create_option` is `Copy` or `FromImage`, then the value must be equal to or greater than the source's size. The size can only be increased.

~> **NOTE:** Changing this value is disruptive if the disk is attached to a Virtual Machine. The VM will be shut down and de-allocated as required by Azure to action the change. Terraform will attempt to start the machine again after the update if it was in a `running` state when the apply was started.

* `edge_zone` - (Optional) Specifies the Edge Zone within the Azure Region where this Managed Disk should exist. Changing this forces a new Managed Disk to be created.

* `encryption_settings` - (Optional) A `encryption_settings` block as defined below.

* `hyper_v_generation` - (Optional) The HyperV Generation of the Disk when the source of an `Import` or `Copy` operation targets a source that contains an operating system. Possible values are `V1` and `V2`. Changing this forces a new resource to be created.

* `image_reference_id` - (Optional) ID of an existing platform/marketplace disk image to copy when `create_option` is `FromImage`. This field cannot be specified if gallery_image_reference_id is specified.

* `gallery_image_reference_id` - (Optional) ID of a Gallery Image Version to copy when `create_option` is `FromImage`. This field cannot be specified if image_reference_id is specified.

* `logical_sector_size` - (Optional) Logical Sector Size. Possible values are: `512` and `4096`. Defaults to `4096`. Changing this forces a new resource to be created.

~> **NOTE:** Setting logical sector size is supported only with `UltraSSD_LRS` disks.

* `os_type` - (Optional) Specify a value when the source of an `Import` or `Copy` operation targets a source that contains an operating system. Valid values are `Linux` or `Windows`.

* `source_resource_id` - (Optional) The ID of an existing Managed Disk to copy `create_option` is `Copy` or the recovery point to restore when `create_option` is `Restore`

* `source_uri` - (Optional) URI to a valid VHD file to be used when `create_option` is `Import`.

* `storage_account_id` - (Optional) The ID of the Storage Account where the `source_uri` is located. Required when `create_option` is set to `Import`.  Changing this forces a new resource to be created.

* `tier` - (Optional) The disk performance tier to use. Possible values are documented [here](https://docs.microsoft.com/azure/virtual-machines/disks-change-performance). This feature is currently supported only for premium SSDs.

~> **NOTE:** Changing this value is disruptive if the disk is attached to a Virtual Machine. The VM will be shut down and de-allocated as required by Azure to action the change. Terraform will attempt to start the machine again after the update if it was in a `running` state when the apply was started.

* `max_shares` - (Optional) The maximum number of VMs that can attach to the disk at the same time. Value greater than one indicates a disk that can be mounted on multiple VMs at the same time.

-> **Note:** Premium SSD maxShares limit: `P15` and `P20` disks: 2. `P30`,`P40`,`P50` disks: 5. `P60`,`P70`,`P80` disks: 10. For ultra disks the `max_shares` minimum value is 1 and the maximum is 5.

* `trusted_launch_enabled` (Optional) Specifies if Trusted Launch is enabled for the Managed Disk. Defaults to `false`.

-> **Note:** Trusted Launch can only be enabled when `create_option` is `FromImage` or `Import`.

* `security_type` - (Optional) Security Type of the Managed Disk when it is used for a Confidential VM. Possible values are `ConfidentialVM_VMGuestStateOnlyEncryptedWithPlatformKey`, `ConfidentialVM_DiskEncryptedWithPlatformKey` and `ConfidentialVM_DiskEncryptedWithCustomerKey`. Changing this forces a new resource to be created.

~> **NOTE:** `security_type` cannot be specified when `trusted_launch_enabled` is set to true.

~> **NOTE:** `secure_vm_disk_encryption_set_id` must be specified when `security_type` is set to `ConfidentialVM_DiskEncryptedWithCustomerKey`.

* `secure_vm_disk_encryption_set_id` - (Optional) The ID of the Disk Encryption Set which should be used to Encrypt this OS Disk when the Virtual Machine is a Confidential VM. Conflicts with `disk_encryption_set_id`. Changing this forces a new resource to be created.

~> **NOTE:** `secure_vm_disk_encryption_set_id` can only be specified when `security_type` is set to `ConfidentialVM_DiskEncryptedWithCustomerKey`.

* `on_demand_bursting_enabled` (Optional) Specifies if On-Demand Bursting is enabled for the Managed Disk. Defaults to `false`.

-> **Note:** Credit-Based Bursting is enabled by default on all eligible disks. More information on [Credit-Based and On-Demand Bursting can be found in the documentation](https://docs.microsoft.com/azure/virtual-machines/disk-bursting#disk-level-bursting).

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `zone` - (Optional) Specifies the Availability Zone in which this Managed Disk should be located. Changing this property forces a new resource to be created.

~> **Note:** Availability Zones are [only supported in select regions at this time](https://docs.microsoft.com/azure/availability-zones/az-overview).

* `network_access_policy` - Policy for accessing the disk via network. Allowed values are `AllowAll`, `AllowPrivate`, and `DenyAll`.

* `disk_access_id` - The ID of the disk access resource for using private endpoints on disks.

~> **Note:** `disk_access_id` is only supported when `network_access_policy` is set to `AllowPrivate`.

* `public_network_access_enabled` - (Optional) Whether it is allowed to access the disk via public network. Defaults to `true`.

For more information on managed disks, such as sizing options and pricing, please check out the [Azure Documentation](https://docs.microsoft.com/azure/storage/storage-managed-disks-overview).

---

The `disk_encryption_key` block supports:

* `secret_url` - (Required) The URL to the Key Vault Secret used as the Disk Encryption Key. This can be found as `id` on the `azurerm_key_vault_secret` resource.

* `source_vault_id` - (Required) The URL of the Key Vault. This can be found as `vault_uri` on the `azurerm_key_vault` resource.

---

The `encryption_settings` block supports:

* `enabled` - (Required) Is Encryption enabled on this Managed Disk? Changing this forces a new resource to be created.

* `disk_encryption_key` - (Optional) A `disk_encryption_key` block as defined above.

* `key_encryption_key` - (Optional) A `key_encryption_key` block as defined below.

---

The `key_encryption_key` block supports:

* `key_url` - (Required) The URL to the Key Vault Key used as the Key Encryption Key. This can be found as `id` on the `azurerm_key_vault_key` resource.

* `source_vault_id` - (Required) The ID of the source Key Vault.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Managed Disk.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Managed Disk.
* `update` - (Defaults to 30 minutes) Used when updating the Managed Disk.
* `read` - (Defaults to 5 minutes) Used when retrieving the Managed Disk.
* `delete` - (Defaults to 30 minutes) Used when deleting the Managed Disk.

## Import

Managed Disks can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_managed_disk.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/microsoft.compute/disks/manageddisk1
```
