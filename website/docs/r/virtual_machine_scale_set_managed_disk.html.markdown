---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_machine_scale_set_managed_disk"
description: |-
  Manages a Managed Disk for use with a Virtual Machine Scale Set.
---

# azurerm_virtual_machine_scale_set_managed_disk

Manages a Managed Disk for use with a Virtual Machine Scale Set.

-> **Note:** Unlike `azurerm_managed_disk`, this resource is able to update a Managed Disk while it is attached to a Virtual Machine Scale Set instance. When a change requires the disk to be taken offline, the disk is hot-detached from the instance, updated, and re-attached without deallocating the instance. Attaching the disk to a Virtual Machine Scale Set instance is not managed by this resource. To manage a Managed Disk attached to a standalone Virtual Machine, use [`azurerm_managed_disk`](managed_disk.html) instead.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_machine_scale_set_managed_disk" "example" {
  name                 = "example-disk"
  resource_group_name  = azurerm_resource_group.example.name
  location             = azurerm_resource_group.example.location
  storage_account_type = "Premium_LRS"

  creation {
    option = "Empty"
  }

  disk_size_gb = 128

  tags = {
    environment = "production"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Managed Disk. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which to create the Managed Disk. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Managed Disk should exist. Changing this forces a new resource to be created.

* `storage_account_type` - (Required) The type of storage to use for the Managed Disk. Possible values are `Standard_LRS`, `StandardSSD_LRS`, `StandardSSD_ZRS`, `Premium_LRS`, `Premium_ZRS`, `PremiumV2_LRS`, and `UltraSSD_LRS`.

* `creation` - (Required) A `creation` block as defined below. Changing this forces a new resource to be created.

---

* `data_access_auth_mode` - (Optional) The authentication mode used when exporting or uploading to the Managed Disk. Possible values are `AzureActiveDirectory` and `None`. Defaults to `None`.

* `disk_access_id` - (Optional) The ID of the Disk Access resource used to control the export of the Managed Disk.

~> **Note:** `disk_access_id` is only supported when `network_access_policy` is set to `AllowPrivate`.

* `disk_encryption_set_id` - (Optional) The ID of a Disk Encryption Set used to encrypt this Managed Disk. Conflicts with `secure_vm_disk_encryption_set_id`.

* `disk_iops_read_only` - (Optional) The number of IOPS allowed across all Virtual Machine Scale Set instances mounting the Managed Disk as read-only. Only settable when `storage_account_type` is `UltraSSD_LRS` or `PremiumV2_LRS` and `max_shares` is greater than `1`.

* `disk_iops_read_write` - (Optional) The number of IOPS allowed for this Managed Disk. Only settable when `storage_account_type` is `UltraSSD_LRS` or `PremiumV2_LRS`.

* `disk_mbps_read_only` - (Optional) The bandwidth allowed across all Virtual Machine Scale Set instances mounting the Managed Disk as read-only, in MB per second. Only settable when `storage_account_type` is `UltraSSD_LRS` or `PremiumV2_LRS` and `max_shares` is greater than `1`.

* `disk_mbps_read_write` - (Optional) The bandwidth allowed for this Managed Disk in MB per second. Only settable when `storage_account_type` is `UltraSSD_LRS` or `PremiumV2_LRS`.

* `disk_size_gb` - (Optional) The size of the Managed Disk in gigabytes. Required when `creation.option` is `Empty`. The size can only be increased.

* `edge_zone` - (Optional) The Edge Zone within the Azure Region where this Managed Disk should exist. Changing this forces a new resource to be created.

* `encryption_settings` - (Optional) An `encryption_settings` block as defined below.

~> **Note:** Removing `encryption_settings` forces a new resource to be created.

* `hyper_v_generation` - (Optional) The Hyper-V Generation of the Managed Disk. Possible values are `V1` and `V2`. Changing this forces a new resource to be created.

* `max_shares` - (Optional) The maximum number of Virtual Machine Scale Set instances that can attach the Managed Disk at once. Possible values are between `2` and `10`.

* `network_access_policy` - (Optional) The policy for accessing the Managed Disk via the network. Possible values are `AllowAll`, `AllowPrivate`, and `DenyAll`. Defaults to `AllowAll`.

* `on_demand_bursting_enabled` - (Optional) Whether On-Demand Bursting is enabled for the Managed Disk. Defaults to `false`.

~> **Note:** `on_demand_bursting_enabled` can only be enabled when `storage_account_type` is `Premium_LRS` or `Premium_ZRS` and `disk_size_gb` is larger than `512`.

* `optimized_frequent_attach_enabled` - (Optional) Whether this Managed Disk should be optimized for frequent attachments, where the disk is attached and detached more than five times a day. Defaults to `false`.

* `os_type` - (Optional) The Operating System type of the Managed Disk. Possible values are `Linux` and `Windows`.

* `public_network_access_enabled` - (Optional) Whether public network access to the Managed Disk is allowed. Defaults to `true`.

* `secure_vm_disk_encryption_set_id` - (Optional) The ID of the Disk Encryption Set used to encrypt this Managed Disk when it is used for a Confidential VM. Conflicts with `disk_encryption_set_id`. Changing this forces a new resource to be created.

~> **Note:** `secure_vm_disk_encryption_set_id` can only be specified when `security_type` is set to `ConfidentialVM_DiskEncryptedWithCustomerKey`.

* `security_type` - (Optional) The Security Type of the Managed Disk when it is used for a Confidential VM. Possible values are `ConfidentialVM_VMGuestStateOnlyEncryptedWithPlatformKey`, `ConfidentialVM_DiskEncryptedWithPlatformKey`, and `ConfidentialVM_DiskEncryptedWithCustomerKey`. Changing this forces a new resource to be created.

~> **Note:** When `security_type` is set to `ConfidentialVM_DiskEncryptedWithCustomerKey` the `creation.option` must be one of `FromImage`, `Import` or `ImportSecure`.

~> **Note:** `security_type` cannot be specified when `trusted_launch_enabled` is set to `true`.

~> **Note:** `secure_vm_disk_encryption_set_id` must be specified when `security_type` is set to `ConfidentialVM_DiskEncryptedWithCustomerKey`.

* `tier` - (Optional) The performance tier to use for the Managed Disk.

~> **Note:** `tier` can only be set when `storage_account_type` is `Premium_LRS` or `Premium_ZRS`.

* `trusted_launch_enabled` - (Optional) Whether Trusted Launch is enabled for the Managed Disk. Changing this forces a new resource to be created.

~> **Note:** `trusted_launch_enabled` can only be set to `true` when `creation.option` is `FromImage`, `Import` or `ImportSecure`.

* `zone` - (Optional) The Availability Zone where the Managed Disk should exist. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the Managed Disk.

---

The `creation` block supports:

* `option` - (Required) The method used to create the Managed Disk. Possible values are `Empty`, `Copy`, `Restore`, `FromImage`, `Import`, `ImportSecure`, and `Upload`. Changing this forces a new resource to be created.

* `gallery_image_reference_id` - (Optional) The ID of the Shared Image Version used to create the Managed Disk. Required when `option` is `FromImage` and a Shared Image is used. Changing this forces a new resource to be created.

* `image_reference_id` - (Optional) The ID of the Platform Image used to create the Managed Disk. Required when `option` is `FromImage` and a Platform Image is used. Changing this forces a new resource to be created.

* `logical_sector_size` - (Optional) The logical sector size in bytes for the Managed Disk. Possible values are `512` and `4096`. Only settable when `storage_account_type` is `UltraSSD_LRS` or `PremiumV2_LRS`. Changing this forces a new resource to be created.

* `performance_plus_enabled` - (Optional) Whether Performance Plus is enabled for the Managed Disk. Defaults to `false`. Changing this forces a new resource to be created.

* `source_resource_id` - (Optional) The ID of an existing Managed Disk, Snapshot or Restore Point used to create the Managed Disk. Required when `option` is `Copy` or `Restore`. Changing this forces a new resource to be created.

* `source_uri` - (Optional) The URI of the source used to create the Managed Disk. Required when `option` is `Import` or `ImportSecure`. Changing this forces a new resource to be created.

* `storage_account_id` - (Optional) The ID of the Storage Account containing the source used to create the Managed Disk. Required when `option` is `Import` or `ImportSecure`. Changing this forces a new resource to be created.

* `upload_size_bytes` - (Optional) The size of the source in bytes used to create the Managed Disk. Required when `option` is `Upload`. Changing this forces a new resource to be created.

---

The `encryption_settings` block supports:

* `disk_encryption_key` - (Required) A `disk_encryption_key` block as defined below.

* `key_encryption_key` - (Optional) A `key_encryption_key` block as defined below.

---

The `disk_encryption_key` block supports:

* `secret_url` - (Required) The URL to the Key Vault Secret used as the Disk Encryption Key. This can be found as `id` on the `azurerm_key_vault_secret` resource.

* `source_vault_id` - (Required) The ID of the source Key Vault. This can be found as `id` on the `azurerm_key_vault` resource.

---

The `key_encryption_key` block supports:

* `key_url` - (Required) The URL to the Key Vault Key used as the Key Encryption Key. This can be found as `id` on the `azurerm_key_vault_key` resource.

* `source_vault_id` - (Required) The ID of the source Key Vault. This can be found as `id` on the `azurerm_key_vault` resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Managed Disk.

* `disk_size_bytes` - The size of the Managed Disk in bytes.

* `unique_id` - The unique ID of the Managed Disk.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Managed Disk.
* `read` - (Defaults to 5 minutes) Used when retrieving the Managed Disk.
* `update` - (Defaults to 30 minutes) Used when updating the Managed Disk.
* `delete` - (Defaults to 30 minutes) Used when deleting the Managed Disk.

## Import

Managed Disks can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_virtual_machine_scale_set_managed_disk.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/disks/manageddisk1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Compute` - 2023-04-02, 2024-03-01
