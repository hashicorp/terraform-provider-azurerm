---
subcategory: "Data Share"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_data_share_dataset_data_lake_gen1"
description: |-
  Gets information about an existing Data Share Data Lake Gen1 Dataset.
---

# Data Source: azurerm_data_share_dataset_data_lake_gen1

Use this data source to access information about an existing Data Share Data Lake Gen1 Dataset.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_data_share_dataset_data_lake_gen1" "example" {
  name          = "example-dsdlg1ds"
  data_share_id = "example-share-id"
}

output "id" {
  value = data.azurerm_data_share_dataset_data_lake_gen1.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Data Share Data Lake Gen1 Dataset.

* `data_share_id` - (Required) The ID of the Data Share in which this Data Share Data Lake Gen1 Dataset should be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Data Share Data Lake Gen1 Dataset.

* `data_lake_store` - A `data_lake_store` block as defined below.

* `folder_path` - The folder path of the data lake store to be shared with the receiver.

* `file_name` - The file name in the folder path of the data lake store to be shared with the receiver.

* `display_name` - The name of the Data Share Dataset.

---

A `data_lake_store` block supports the following:

* `name` - The name of the data lake store to be shared with the receiver.

* `resource_group_name` - The resource group name of the data lake store to be shared with the receiver.

* `subscription_id` - The subscription id of the data lake store to be shared with the receiver.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Data Share Data Lake Gen1 Dataset.
