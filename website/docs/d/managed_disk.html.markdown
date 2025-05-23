---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_disk"
description: |-
  Get information about an existing Managed Disk.
---

# Data Source: azurerm_managed_disk

Use this data source to access information about an existing Managed Disk.

## Example Usage

```hcl
data "azurerm_managed_disk" "existing" {
  name                = "example-datadisk"
  resource_group_name = "example-resources"
}

output "id" {
  value = data.azurerm_managed_disk.existing.id
}
```

## Argument Reference

* `name` - Specifies the name of the Managed Disk.

* `resource_group_name` - Specifies the name of the Resource Group where this Managed Disk exists.

## Attributes Reference

* `disk_encryption_set_id` - The ID of the Disk Encryption Set used to encrypt this Managed Disk.

* `disk_iops_read_write` - The number of IOPS allowed for this disk, where one operation can transfer between 4k and 256k bytes.

* `disk_mbps_read_write` - The bandwidth allowed for this disk.

* `disk_size_gb` - The size of the Managed Disk in gigabytes.

* `image_reference_id` - The ID of the source image used for creating this Managed Disk.

* `location` - The Azure location of the Managed Disk.

* `os_type` - The operating system used for this Managed Disk.

* `storage_account_type` - The storage account type for the Managed Disk.

* `source_uri` - The Source URI for this Managed Disk.

* `source_resource_id` - The ID of an existing Managed Disk which this Disk was created from.

* `storage_account_id` - The ID of the Storage Account where the `source_uri` is located.

* `tags` - A mapping of tags assigned to the resource.

* `zones` - A list of Availability Zones where the Managed Disk exists.

* `network_access_policy` - Policy for accessing the disk via network.

* `disk_access_id` - The ID of the disk access resource for using private endpoints on disks.

* `encryption_settings` - A `encryption_settings` block as defined below.

---

The `encryption_settings` block supports:

* `disk_encryption_key` - A `disk_encryption_key` block as defined above.

* `key_encryption_key` - A `key_encryption_key` block as defined below.

---

The `disk_encryption_key` block supports:

* `secret_url` - The URL to the Key Vault Secret used as the Disk Encryption Key.

* `source_vault_id` - The ID of the source Key Vault.

---

The `key_encryption_key` block supports:

* `key_url` - The URL to the Key Vault Key used as the Key Encryption Key.

* `source_vault_id` - The ID of the source Key Vault.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Managed Disk.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Compute`: 2023-04-02
