---
subcategory: "Import Export"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_export_job"
description: |-
  Manages an Azure Export Job.
---

# azurerm_export_job

Manages an Azure Export Job.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "example-sa"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_export_job" "test" {
  name                = "example-export-job"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  storage_account_id        = azurerm_storage_account.example.id
  export_blob_path_prefixes = ["/"]

  return_shipping {
    carrier_account_number = "123456789"
    carrier_name           = "DHL"
  }

  return_address {
    recipient_name    = "Tets"
    street_address1   = "Street1"
    street_address2   = "street2"
    city              = "Redmond"
    state_or_province = "wa"
    postal_code       = "98007"
    country_or_region = "USA"
    phone             = "4250000000"
    email             = "Test@contoso.com"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Azure Export Job. Changing this forces a new resource to be created.

* `location` - (Required) The Azure location where the Azure Export Job should exist. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Azure Export Job should exist. Changing this forces a new resource to be created.

* `storage_account_id` - (Required) The resource ID of the storage account where data will be exported from. Changing this forces a new resource to be created.

* `export_blob_paths` - (Optional) A collection of blob-path strings. At least one of `export_blob_paths`, `export_blob_path_prefixes` or `export_blob_list_path` must be specified. Changing this forces a new resource to be created.

* `export_blob_path_prefixes` - (Optional) A collection of blob-prefix strings. At least one of `export_blob_paths`, `export_blob_path_prefixes` or `export_blob_list_path` must be specified. Changing this forces a new resource to be created.

* `export_blob_list_path` - (Optional) The relative URI to the block blob that contains the list of blob paths or blob path prefixes as defined above, beginning with the container name. If the blob is in root container, the URI must begin with `$root`. At least one of `export_blob_paths`, `export_blob_path_prefixes` or `export_blob_list_path` must be specified. Changing this forces a new resource to be created.

* `return_address` - (Required) One `return_address` block as defined below.

* `return_shipping` - (Required) One `return_shipping` block as defined below.

* `backup_drive_manifest` - (Optional) Should the manifest files on the drives be copied to block blobs? Defaults to `false`.

* `diagnostics_path` - (Optional) Specifies the virtual blob directory to which the copy logs and backups of drive manifest files (if enabled) will be stored. Defaults to `waimportexport`. Changing this forces a new resource to be created.

* `log_level` - (Optional) Indicates whether error logging or verbose logging will be enabled. Possible values are `Error` or `Verbose`. Defaults to `Error`.

---

A `return_address` block exports the following:

* `city` - (Required) The city name to use when returning the drives.

* `country_or_region` - (Required) The country or region to use when returning the drives.

* `email` - (Required) Email address of the recipient of the returned drives.

* `phone` - (Required) Phone number of the recipient of the returned drives.

* `postal_code` - (Required) The postal code to use when returning the drives.

* `recipient_name` - (Required) The name of the recipient who will receive the hard drives when they are returned.

* `street_address1` - (Required) The first line of the street address to use when returning the drives.

* `state_or_province` - (Optional) The state or province to use when returning the drives.

* `street_address2` - (Optional) The first line of the street address to use when returning the drives.

---

A `return_shipping` block exports the following:

* `carrier_account_number` - (Required) The customer's account number with the carrier.

* `carrier_name` - (Required) The carrier's name.

## Attributes Reference

The following attributes are exported:

* `id` - The Resource ID of the Azure Export Job.

* `shipping_information` - One `shipping_information` block as defined below.

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

* `create` - (Defaults to 30 minutes) Used when creating the Azure Export Job.
* `update` - (Defaults to 30 minutes) Used when updating the Azure Export Job.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Export Job.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Export Job.

## Import

Azure Export Job can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_export_job.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ImportExport/jobs/job1
```
