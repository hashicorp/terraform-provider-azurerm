---
subcategory: "DataBox"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_databox_job"
description: |-
  Get information about an existing DataBox.
---

# Data Source: azurerm_databox_job

Use this data source to access information about an existing DataBox.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_databox_job" "existing" {
  name                = "example-databoxjob"
  resource_group_name = "example-resources"
}

output "id" {
  value = data.azurerm_databox_job.existing.id
}
```

## Argument Reference

* `name` - Specifies the name of the DataBox.

* `resource_group_name` - Specifies the name of the Resource Group where this DataBox exists.

## Attributes Reference

* `location` - The Azure location where the resource exists.

* `destination_account` - One or more `destination_account` block defined below.

* `device_password` - The device password for unlocking DataBox Heavy.

* `sku_name` - The sku name.

* `tags` - A mapping of tags to assign to the resource.

---

A `destination_account` block exports the following:

* `type` -The destination account type.

* `managed_disk_resource_group_id` - The destination Resource Group Id where the Compute disks should be created.

* `managed_disk_staging_storage_account_id` - The Id of the storage account that can be used to copy the vhd for staging.

* `share_password` - The share password to be shared by all shares in SA.

* `storage_account_id` - The Id of the destination where the data has to be moved.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the DataBox.
