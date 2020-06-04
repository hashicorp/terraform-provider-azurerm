---
subcategory: "Data Share"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_data_share_dataset_blob_storage"
description: |-
  Gets information about an existing Data Share Blob Storage Dataset.
---

# Data Source: azurerm_data_share_dataset_blob_storage

Use this data source to access information about an existing Data Share Blob Storage Dataset.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_data_share_dataset_blob_storage" "example" {
  name          = "example-dsbsds"
  data_share_id = "example-share-id"
}

output "id" {
  value = data.azurerm_data_share_dataset_blob_storage.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Data Share Blob Storage Dataset.

* `data_share_id` - (Required) The ID of the Data Share in which this Data Share Blob Storage Dataset should be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Data Share Blob Storage Dataset.

* `container_name` - The name of the storage account container to be shared with the receiver.

* `storage_account` - A `storage_account` block as defined below.

* `file_path` - The path of the file in the storage container to be shared with the receiver.

* `folder_path` - The folder path of the file in the storage container to be shared with the receiver.

* `display_name` - The name of the Data Share Dataset.

---

A `storage_account` block supports the following:

* `name` - The name of the storage account to be shared with the receiver. 

* `resource_group_name` - The resource group name of the storage account to be shared with the receiver.

* `subscription_id` - The subscription id of the storage account to be shared with the receiver.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Data Share Blob Storage Dataset.
