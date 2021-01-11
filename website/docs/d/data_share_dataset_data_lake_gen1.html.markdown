---
subcategory: "Data Share"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_data_share_dataset_data_lake_gen1"
description: |-
  Gets information about an existing DataShareDataLakeGen1Dataset.
---

# Data Source: azurerm_data_share_dataset_data_lake_gen1

Use this data source to access information about an existing DataShareDataLakeGen1Dataset.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_data_share_dataset_data_lake_gen1" "example" {
  name          = "example-dsdsdlg1"
  data_share_id = "example-share-id"
}

output "id" {
  value = data.azurerm_data_share_dataset_data_lake_gen1.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Data Share Data Lake Gen1 Dataset.

* `data_share_id` - (Required) The resource ID of the Data Share where this Data Share Data Lake Gen1 Dataset should be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The resource ID of the Data Share Data Lake Gen1 Dataset.

* `data_lake_store_id` - The resource ID of the Data Lake Store to be shared with the receiver.

* `display_name` - The displayed name of the Data Share Dataset.

* `file_name` - The file name of the data lake store to be shared with the receiver.

* `folder_path` - The folder path of the data lake store to be shared with the receiver.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the DataShareDataLakeGen1Dataset.
