---
subcategory: "Data Share"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_data_share_dataset_data_lake_gen2"
description: |-
  Gets information about an existing Data Share Data Lake Gen2 Dataset.
---

# Data Source: azurerm_data_share_dataset_data_lake_gen2

Use this data source to access information about an existing Data Share Data Lake Gen2 Dataset.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_data_share_dataset_data_lake_gen2" "example" {
  name     = "example-dsdlg2ds"
  share_id = "example-share-id"
}

output "id" {
  value = data.azurerm_data_share_dataset_data_lake_gen2.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Data Share Data Lake Gen2 Dataset.

* `share_id` - (Required) The resource ID of the Data Share where this Data Share Data Lake Gen2 Dataset should be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The resource ID of the Data Share Data Lake Gen2 Dataset.

* `display_name` - The name of the Data Share Dataset.

* `file_path` - The path of the file in the data lake file system to be shared with the receiver.

* `file_system_name` - The name of the data lake file system to be shared with the receiver.

* `folder_path` - The folder path in the data lake file system to be shared with the receiver.

* `storage_account_id` - The resource ID of the storage account of the data lake file system to be shared with the receiver.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Data Share Data Lake Gen2 Dataset.
