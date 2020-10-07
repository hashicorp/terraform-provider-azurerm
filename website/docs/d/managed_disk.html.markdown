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

* `os_type` - The operating system used for this Managed Disk.

* `storage_account_type` - The storage account type for the Managed Disk.

* `source_uri` - The Source URI for this Managed Disk.

* `source_resource_id` - The ID of an existing Managed Disk which this Disk was created from.

* `storage_account_id` - The ID of the Storage Account where the `source_uri` is located.

* `tags` - A mapping of tags assigned to the resource.

* `zones` - A list of Availability Zones where the Managed Disk exists.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Managed Disk.
