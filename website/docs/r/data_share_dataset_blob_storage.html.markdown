---
subcategory: "Data Share"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_share_dataset_blob_storage"
description: |-
  Manages a Data Share Blob Storage Dataset.
---

# azurerm_data_share_dataset_blob_storage

Manages a Data Share Blob Storage Dataset.

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
  account_tier             = "Standard"
  account_replication_type = "RAGRS"
}

resource "azurerm_storage_container" "example" {
  name                  = "example-sc"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "container"
}

data "azuread_service_principal" "example" {
  display_name = azurerm_data_share_account.example.name
}

resource "azurerm_role_assignment" "example" {
  scope                = azurerm_storage_account.example.id
  role_definition_name = "Storage Blob Data Reader"
  principal_id         = data.azuread_service_principal.example.object_id
}

resource "azurerm_data_share_dataset_blob_storage" "example" {
  name           = "example-dsbsds-file"
  data_share_id  = azurerm_data_share.example.id
  container_name = azurerm_storage_container.example.name
  storage_account {
    name                = azurerm_storage_account.example.name
    resource_group_name = azurerm_storage_account.example.resource_group_name
    subscription_id     = "00000000-0000-0000-0000-000000000000"
  }
  file_path = "myfile.txt"
  depends_on = [
    azurerm_role_assignment.example,
  ]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Data Share Blob Storage Dataset. Changing this forces a new Data Share Blob Storage Dataset to be created.

* `data_share_id` - (Required) The ID of the Data Share in which this Data Share Blob Storage Dataset should be created. Changing this forces a new Data Share Blob Storage Dataset to be created.

* `container_name` - (Required) The name of the storage account container to be shared with the receiver. Changing this forces a new Data Share Blob Storage Dataset to be created.

* `storage_account` - (Required) A `storage_account` block as defined below.

* `file_path` - (Optional) The path of the file in the storage container to be shared with the receiver. Changing this forces a new Data Share Blob Storage Dataset to be created.

* `folder_path` - (Optional) The path of the folder in the storage container to be shared with the receiver. Changing this forces a new Data Share Blob Storage Dataset to be created.

---

A `storage_account` block supports the following:

* `name` - (Required)  The name of the storage account to be shared with the receiver. Changing this forces a new Data Share Blob Storage Dataset to be created.

* `resource_group_name` - (Required)  The resource group name of the storage account to be shared with the receiver. Changing this forces a new Data Share Blob Storage Dataset to be created.

* `subscription_id` - (Required) The subscription id of the storage account to be shared with the receiver. Changing this forces a new Data Share Blob Storage Dataset to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Data Share Blob Storage Dataset.

* `display_name` - The name of the Data Share Dataset.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Share Blob Storage Dataset.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Share Blob Storage Dataset.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Share Blob Storage Dataset.

## Import

Data Share Blob Storage Datasets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_share_dataset_blob_storage.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DataShare/accounts/account1/shares/share1/dataSets/dataSet1
```
