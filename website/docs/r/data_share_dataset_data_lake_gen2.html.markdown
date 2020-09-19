---
subcategory: "Data Share"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_share_dataset_data_lake_gen2"
description: |-
  Manages a Data Share Data Lake Gen2 Dataset.
---

# azurerm_data_share_dataset_data_lake_gen2

Manages a Data Share Data Lake Gen2 Dataset.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_data_share_account" "example" {
  name                = "example-dsa"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_data_share" "example" {
  name       = "example_ds"
  account_id = azurerm_data_share_account.example.id
  kind       = "CopyBased"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestr"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "example" {
  name               = "example-dlg2fs"
  storage_account_id = azurerm_storage_account.example.id
}

data "azuread_service_principal" "example" {
  display_name = azurerm_data_share_account.example.name
}

resource "azurerm_role_assignment" "example" {
  scope                = azurerm_storage_account.example.id
  role_definition_name = "Storage Blob Data Reader"
  principal_id         = data.azuread_service_principal.example.object_id
}

resource "azurerm_data_share_dataset_data_lake_gen2" "example" {
  name               = "accexample-dlg2ds"
  share_id           = azurerm_data_share.example.id
  storage_account_id = azurerm_storage_account.example.id
  file_system_name   = azurerm_storage_data_lake_gen2_filesystem.example.name
  file_path          = "myfile.txt"
  depends_on = [
    azurerm_role_assignment.example,
  ]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Data Share Data Lake Gen2 Dataset. Changing this forces a new Data Share Data Lake Gen2 Dataset to be created.

* `share_id` - (Required) The resource ID of the Data Share where this Data Share Data Lake Gen2 Dataset should be created. Changing this forces a new Data Share Data Lake Gen2 Dataset to be created.

* `file_system_name` - (Required) The name of the data lake file system to be shared with the receiver. Changing this forces a new Data Share Data Lake Gen2 Dataset to be created.

* `storage_account_id` - (Required) The resource id of the storage account of the data lake file system to be shared with the receiver. Changing this forces a new Data Share Data Lake Gen2 Dataset to be created.

---

* `file_path` - (Optional) The path of the file in the data lake file system to be shared with the receiver. Conflicts with `folder_path` Changing this forces a new Data Share Data Lake Gen2 Dataset to be created.

* `folder_path` - (Optional) The folder path in the data lake file system to be shared with the receiver. Conflicts with `file_path` Changing this forces a new Data Share Data Lake Gen2 Dataset to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The resource ID of the Data Share Data Lake Gen2 Dataset.

* `display_name` - The name of the Data Share Dataset.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Share Data Lake Gen2 Dataset.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Share Data Lake Gen2 Dataset.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Share Data Lake Gen2 Dataset.

## Import

Data Share Data Lake Gen2 Datasets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_share_dataset_data_lake_gen2.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DataShare/accounts/account1/shares/share1/dataSets/dataSet1
```
