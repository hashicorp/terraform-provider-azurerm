---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_table_entity"
description: |-
  Manages an Entity within a Table in an Azure Storage Account.
---

# azurerm_storage_table_entity

Manages an Entity within a Table in an Azure Storage Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "azureexample"
  location = "westus"
}

resource "azurerm_storage_account" "example" {
  name                     = "azureexamplestorage1"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_table" "example" {
  name                 = "myexampletable"
  resource_group_name  = azurerm_resource_group.example.name
  storage_account_name = azurerm_storage_account.example.name
}

resource "azurerm_storage_table_entity" "example" {
  storage_account_name = azurerm_storage_account.example.name
  table_name           = azurerm_storage_table.example.name

  partition_key = "examplepartition"
  row_key       = "examplerow"

  entity = {
    example = "example"
  }
}
```

## Argument Reference

The following arguments are supported:

* `storage_account_name` - (Required) Specifies the storage account in which to create the storage table entity.
 Changing this forces a new resource to be created.

* `table_name` - (Required) The name of the storage table in which to create the storage table entity.
Changing this forces a new resource to be created.

* `partition_key` - (Required) The key for the partition where the entity will be inserted/merged. Changing this forces a new resource.

* `row_key` - (Required) The key for the row where the entity will be inserted/merged. Changing this forces a new resource.

* `entity` - (Required) A map of key/value pairs that describe the entity to be inserted/merged in to the storage table.


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the Entity within the Table in the Storage Account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Table Entity.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Table Entity.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Table Entity.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Table Entity.

## Import

Entities within a Table in an Azure Storage Account can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_table_entity.entity1 https://example.table.core.windows.net/table1(PartitionKey='samplepartition',RowKey='samplerow')
```
