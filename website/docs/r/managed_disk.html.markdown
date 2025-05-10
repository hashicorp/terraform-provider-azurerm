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

* `resource_group_name` - (Required) The name of the Resource Group where the Managed Disk should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specified the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `storage_account_type` - (Required) The type of storage to use for the managed disk. Possible values are `Standard_LRS`, `StandardSSD_ZRS`, `Premium_LRS`, `PremiumV2_LRS`, `Premium_ZRS`, `StandardSSD_LRS` or `UltraSSD_LRS`.

-> **Note:** Azure Ultra Disk Storage is only available in a region that support availability zones and can only enabled on the following VM series: `ESv3`, `DSv3`, `FSv3`, `LSv2`, `M` and `Mv2`. For more information see the `Azure Ultra Disk Storage` [product documentation](https://docs.microsoft.com/azure/virtual-machines/windows/disks-enable-ultra-ssd).

* `create_option` - (Required) The method to use when creating the managed disk. Changing this forces a new resource to be created. Possible values include: * `Import` - Import a VHD file in to the managed disk (VHD specified with `source_uri`). * `ImportSecure` - Securely import a VHD file in to the managed disk (VHD specified with `source_uri`). * `Empty` - Create an empty managed disk. * `Copy` - Copy an existing managed disk or snapshot (specified with `source_resource_id`). * `FromImage` - Copy a Platform Image (specified with `image_reference_id`) * `Restore` - Set by Azure Backup or Site Recovery on a restored disk (specified with `source_resource_id`). * `Upload` - Upload a VHD disk with the help of SAS URL (to be used with `upload_size_bytes`).

---

* `disk_encryption_set_id` - (Optional) The ID of a Disk Encryption Set which should be used to encrypt this Managed Disk. Conflicts with `secure_vm_disk_encryption_set_id`.

~> **Note:** The Disk Encryption Set must have the `Reader` Role Assignment scoped on the Key Vault - in addition to an Access Policy to the Key Vault

~> **Note:** Disk Encryption Sets are in Public Preview in a limited set of regions

* `disk_iops_read_write` - (Optional) The number of IOPS allowed for this disk; only settable for UltraSSD disks and PremiumV2 disks. One operation can transfer between 4k and 256k bytes.

* `disk_mbps_read_write` - (Optional) The bandwidth allowed for this disk; only settable for UltraSSD disks and PremiumV2 disks. MBps means millions of bytes per second.

* `disk_iops_read_only` - (Optional) The number of IOPS allowed across all VMs mounting the shared disk as read-only; only settable for UltraSSD disks and PremiumV2 disks with shared disk enabled. One operation can transfer between 4k and 256k bytes.

* `disk_mbps_read_only` - (Optional) The bandwidth allowed across all VMs mounting the shared disk as read-only; only settable for UltraSSD disks and PremiumV2 disks with shared disk enabled. MBps means millions of bytes per second.

* `upload_size_bytes` - (Optional) Specifies the size of the managed disk to create in bytes. Required when `create_option` is `Upload`. The value must be equal to the source disk to be copied in bytes. Source disk size could be calculated with `ls -l` or `wc -c`. More information can be found at [Copy a managed disk](https://learn.microsoft.com/en-us/azure/virtual-machines/linux/disks-upload-vhd-to-managed-disk-cli#copy-a-managed-disk). Changing this forces a new resource to be created.

* `disk_size_gb` - (Optional) (Optional, Required for a new managed disk) Specifies the size of the managed disk to create in gigabytes. If `create_option` is `Copy` or `FromImage`, then the value must be equal to or greater than the source's size. The size can only be increased.

-> **Note:** In certain conditions the Data Disk size can be updated without shutting down the Virtual Machine, however only a subset of Virtual Machine SKUs/Disk combinations support this. More information can be found [for Linux Virtual Machines](https://learn.microsoft.com/en-us/azure/virtual-machines/linux/expand-disks?tabs=azure-cli%2Cubuntu#expand-without-downtime) and [Windows Virtual Machines](https://learn.microsoft.com/azure/virtual-machines/windows/expand-os-disk#expand-without-downtime) respectively.

~> **Note:** If No Downtime Resizing is not available, be aware that changing this value is disruptive if the disk is attached to a Virtual Machine. The VM will be shut down and de-allocated as required by Azure to action the change. Terraform will attempt to start the machine again after the update if it was in a `running` state when the apply was started.

~> **Note:** When upgrading `disk_size_gb` from a value less than 4095 to one greater than 4095, and if `storage_account_type` is not set to `PremiumV2_LRS` or `UltraSSD_LRS`, the disk will be detached from its associated Virtual Machine as required by Azure to action the change. Terraform will attempt to reattach the disk again after the update.

* `edge_zone` - (Optional) Specifies the Edge Zone within the Azure Region where this Managed Disk should exist. Changing this forces a new Managed Disk to be created.

* `encryption_settings` - (Optional) A `encryption_settings` block as defined below.

~> **Note:** Removing `encryption_settings` forces a new resource to be created.

* `hyper_v_generation` - (Optional) The HyperV Generation of the Disk when the source of an `Import` or `Copy` operation targets a source that contains an operating system. Possible values are `V1` and `V2`. For `ImportSecure` it must be set to `V2`. Changing this forces a new resource to be created.

* `image_reference_id` - (Optional) ID of an existing platform/marketplace disk image to copy when `create_option` is `FromImage`. This field cannot be specified if gallery_image_reference_id is specified. Changing this forces a new resource to be created.

* `gallery_image_reference_id` - (Optional) ID of a Gallery Image Version to copy when `create_option` is `FromImage`. This field cannot be specified if image_reference_id is specified. Changing this forces a new resource to be created.

* `logical_sector_size` - (Optional) Logical Sector Size. Possible values are: `512` and `4096`. Defaults to `4096`. Changing this forces a new resource to be created.

~> **Note:** Setting logical sector size is supported only with `UltraSSD_LRS` disks and `PremiumV2_LRS` disks.

* `optimized_frequent_attach_enabled` - (Optional) Specifies whether this Managed Disk should be optimized for frequent disk attachments (where a disk is attached/detached more than 5 times in a day). Defaults to `false`.

-> **Note:** Setting `optimized_frequent_attach_enabled` to `true` causes the disks to not align with the fault domain of the Virtual Machine, which can have operational implications.

* `performance_plus_enabled` - (Optional) Specifies whether Performance Plus is enabled for this Managed Disk. Defaults to `false`. Changing this forces a new resource to be created.

* `os_type` - (Optional) Specify a value when the source of an `Import`, `ImportSecure` or `Copy` operation targets a source that contains an operating system. Valid values are `Linux` or `Windows`.

* `source_resource_id` - (Optional) The ID of an existing Managed Disk or Snapshot to copy when `create_option` is `Copy` or the recovery point to restore when `create_option` is `Restore`. Changing this forces a new resource to be created.

* `source_uri` - (Optional) URI to a valid VHD file to be used when `create_option` is `Import` or `ImportSecure`. Changing this forces a new resource to be created.

* `storage_account_id` - (Optional) The ID of the Storage Account where the `source_uri` is located. Required when `create_option` is set to `Import` or `ImportSecure`. Changing this forces a new resource to be created.

* `tier` - (Optional) The disk performance tier to use. Possible values are documented [here](https://docs.microsoft.com/azure/virtual-machines/disks-change-performance). This feature is currently supported only for premium SSDs.

~> **Note:** Changing this value is disruptive if the disk is attached to a Virtual Machine. The VM will be shut down and de-allocated as required by Azure to action the change. Terraform will attempt to start the machine again after the update if it was in a `running` state when the apply was started.

* `max_shares` - (Optional) The maximum number of VMs that can attach to the disk at the same time. Value greater than one indicates a disk that can be mounted on multiple VMs at the same time.

-> **Note:** Premium SSD maxShares limit: `P15` and `P20` disks: 2. `P30`,`P40`,`P50` disks: 5. `P60`,`P70`,`P80` disks: 10. For ultra disks the `max_shares` minimum value is 1 and the maximum is 5.

* `trusted_launch_enabled` - (Optional) Specifies if Trusted Launch is enabled for the Managed Disk. Changing this forces a new resource to be created.

-> **Note:** Trusted Launch can only be enabled when `create_option` is `FromImage` or `Import`.

* `security_type` - (Optional) Security Type of the Managed Disk when it is used for a Confidential VM. Possible values are `ConfidentialVM_VMGuestStateOnlyEncryptedWithPlatformKey`, `ConfidentialVM_DiskEncryptedWithPlatformKey` and `ConfidentialVM_DiskEncryptedWithCustomerKey`. Changing this forces a new resource to be created.

~> **Note:** When `security_type` is set to `ConfidentialVM_DiskEncryptedWithCustomerKey` the value of `create_option` must be one of `FromImage` or `ImportSecure`.


~> **Note:** `security_type` cannot be specified when `trusted_launch_enabled` is set to true.

~> **Note:** `secure_vm_disk_encryption_set_id` must be specified when `security_type` is set to `ConfidentialVM_DiskEncryptedWithCustomerKey`.

* `secure_vm_disk_encryption_set_id` - (Optional) The ID of the Disk Encryption Set which should be used to Encrypt this OS Disk when the Virtual Machine is a Confidential VM. Conflicts with `disk_encryption_set_id`. Changing this forces a new resource to be created.

~> **Note:** `secure_vm_disk_encryption_set_id` can only be specified when `security_type` is set to `ConfidentialVM_DiskEncryptedWithCustomerKey`.

* `on_demand_bursting_enabled` - (Optional) Specifies if On-Demand Bursting is enabled for the Managed Disk.

-> **Note:** Credit-Based Bursting is enabled by default on all eligible disks. More information on [Credit-Based and On-Demand Bursting can be found in the documentation](https://docs.microsoft.com/azure/virtual-machines/disk-bursting#disk-level-bursting).

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `zone` - (Optional) Specifies the Availability Zone in which this Managed Disk should be located. Changing this property forces a new resource to be created.

~> **Note:** Availability Zones are [only supported in select regions at this time](https://docs.microsoft.com/azure/availability-zones/az-overview).

* `network_access_policy` - (Optional) Policy for accessing the disk via network. Allowed values are `AllowAll`, `AllowPrivate`, and `DenyAll`.

* `disk_access_id` - (Optional) The ID of the disk access resource for using private endpoints on disks.

~> **Note:** `disk_access_id` is only supported when `network_access_policy` is set to `AllowPrivate`.

* `public_network_access_enabled` - (Optional) Whether it is allowed to access the disk via public network. Defaults to `true`.

For more information on managed disks, such as sizing options and pricing, please check out the [Azure Documentation](https://docs.microsoft.com/azure/storage/storage-managed-disks-overview).

---

The `disk_encryption_key` block supports:

* `secret_url` - (Required) The URL to the Key Vault Secret used as the Disk Encryption Key. This can be found as `id` on the `azurerm_key_vault_secret` resource.

* `source_vault_id` - (Required) The ID of the source Key Vault. This can be found as `id` on the `azurerm_key_vault` resource.

---

The `encryption_settings` block supports:

* `disk_encryption_key` - (Optional) A `disk_encryption_key` block as defined above.

* `key_encryption_key` - (Optional) A `key_encryption_key` block as defined below.

---

The `key_encryption_key` block supports:

* `key_url` - (Required) The URL to the Key Vault Key used as the Key Encryption Key. This can be found as `id` on the `azurerm_key_vault_key` resource.

* `source_vault_id` - (Required) The ID of the source Key Vault. This can be found as `id` on the `azurerm_key_vault` resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Managed Disk.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Managed Disk.
* `read` - (Defaults to 5 minutes) Used when retrieving the Managed Disk.
* `update` - (Defaults to 30 minutes) Used when updating the Managed Disk.
* `delete` - (Defaults to 30 minutes) Used when deleting the Managed Disk.

## Import

Managed Disks can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_managed_disk.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/disks/manageddisk1
```
