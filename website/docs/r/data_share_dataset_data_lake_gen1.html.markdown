---
subcategory: "Data Share"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_share_dataset_data_lake_gen1"
description: |-
  Manages a Data Share Data Lake Gen1 Dataset.
---

# azurerm_data_share_dataset_data_lake_gen1

Manages a Data Share Data Lake Gen1 Dataset.

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

resource "azurerm_data_lake_store" "example" {
  name                = "exampledls"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  firewall_state      = "Disabled"
}

resource "azurerm_data_lake_store_file" "example" {
  account_name     = azurerm_data_lake_store.example.name
  local_file_path  = "./example/myfile.txt"
  remote_file_path = "/example/myfile.txt"
}

data "azuread_service_principal" "example" {
  display_name = azurerm_data_share_account.example.name
}

resource "azurerm_role_assignment" "example" {
  scope                = azurerm_data_lake_store.example.id
  role_definition_name = "Owner"
  principal_id         = data.azuread_service_principal.example.object_id
}

resource "azurerm_data_share_dataset_data_lake_gen1" "example" {
  name               = "example-dlg1ds"
  data_share_id      = azurerm_data_share.example.id
  data_lake_store_id = azurerm_data_lake_store.example.id
  file_name          = "myfile.txt"
  folder_path        = "example"
  depends_on = [
    azurerm_role_assignment.example,
  ]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Data Share Data Lake Gen1 Dataset. Changing this forces a new Data Share Data Lake Gen1 Dataset to be created.

* `data_share_id` - (Required) The resource ID of the Data Share where this Data Share Data Lake Gen1 Dataset should be created. Changing this forces a new Data Share Data Lake Gen1 Dataset to be created.

* `data_lake_store_id` - (Required) The resource ID of the Data Lake Store to be shared with the receiver.

* `folder_path` - (Required) The folder path of the data lake store to be shared with the receiver. Changing this forces a new Data Share Data Lake Gen1 Dataset to be created.

* `file_name` - (Optional) The file name of the data lake store to be shared with the receiver. Changing this forces a new Data Share Data Lake Gen1 Dataset to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The resource ID of the Data Share Data Lake Gen1 Dataset.

* `display_name` - The displayed name of the Data Share Dataset.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Share Data Lake Gen1 Dataset.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Share Data Lake Gen1 Dataset.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Share Data Lake Gen1 Dataset.

## Import

Data Share Data Lake Gen1 Datasets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_share_dataset_data_lake_gen1.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DataShare/accounts/account1/shares/share1/dataSets/dataSet1
```
