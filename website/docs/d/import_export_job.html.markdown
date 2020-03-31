---
subcategory: "Import Export"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_import_export_job"
description: |-
  Gets information about an existing Azure Import/Export Job
---

# Data Source: azurerm_import_export_job

Uses this data source to access information about an existing Azure Import/Export Job


## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_import_export_job" "example" {
  name                = "example-job"
  resource_group_name = "example-resource-group"
}

output "import_export_job_id" {
  value = data.azurerm_import_export_job.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Azure Import/Export Job.

* `resource_group_name` - (Required) Specifies the name of the resource group the Job is located in.

## Attributes Reference

The following attributes are exported:

* `location` - The Azure Location where the Azure Import/Export Job exists.

* `storage_account_id` - The resource ID of the storage account where data will be imported to or exported from .

* `export_blob_paths` - The collection files of blob-path that need to be exported.

* `export_blob_path_prefixes` - The collection files of blob-prefix that need to be exported.

* `export_blob_list_path` - The relative URI to the block blob that contains the list of blob paths or blob path prefixes

* `backup_drive_manifest` - Indicates the manifest files on the drives should be copied to block blobs

* `diagnostics_path` - The virtual blob directory to which the copy logs and backups of drive manifest files (if enabled) will be stored.

* `log_level` - The log level of Azure Import/Export job.

* `drives` - One or more `drive` block defined below.

* `return_address` - One `return_address` block defined below.

* `return_shipping` - One `return_shipping` block defined below.

* `shipping_information` - One `shipping_information` block defined below.

---

A `drive` block exports the following:

* `bit_locker_key` - The BitLocker key used to encrypt the drive.

* `drive_id` - The drive's hardware serial number(without spaces).

* `manifest_file` - The relative path of the manifest file on the drive.

* `manifest_hash` - The Base16-encoded MD5 hash of the manifest file on the drive.

---

A `return_address` block exports the following:

* `city` - The city name to use when returning the drives.

* `country_or_region` - The country or region to use when returning the drives.

* `email` - Email address of the recipient of the returned drives.

* `phone` - Phone number of the recipient of the returned drives.

* `postal_code` - The postal code to use when returning the drives.

* `recipient_name` - The name of the recipient who will receive the hard drives when they are returned.

* `street_address1` - The first line of the street address to use when returning the drives.

* `state_or_province` - The state or province to use when returning the drives.

* `street_address2` - The first line of the street address to use when returning the drives.

---

A `return_shipping` block exports the following:

* `carrier_account_number` - The customer's account number with the carrier.

* `carrier_name` - The carrier's name.

---

A `shipping_information` block exports the following:

* `city` - The city name to use when returning the drives.

* `country_or_region` - The country or region to use when returning the drives.

* `phone` - Phone number of the recipient of the returned drives.

* `postal_code` - The postal code to use when returning the drives.

* `recipient_name` - The name of the recipient who will receive the hard drives when they are returned.

* `state_or_province` - The state or province to use when returning the drives.

* `street_address1` - The first line of the street address to use when returning the drives.

* `street_address2` - The first line of the street address to use when returning the drives.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Import/Export Job.
