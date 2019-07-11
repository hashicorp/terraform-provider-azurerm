---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_table_entity"
sidebar_current: "docs-azurerm-resource-storage-table-entity"
description: |-
  Manages an Entity within a Table in an Azure Storage Account.
---

# azurerm_storage_table_entity

Manages an Entity within a Table in an Azure Storage Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "azuretest"
  location = "westus"
}

resource "azurerm_storage_account" "test" {
  name                     = "azureteststorage1"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "westus"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_table" "test" {
  name                 = "mysampletable"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"
}

resource "azurerm_storage_table_entity" "test" {
  storage_account_name = "${azurerm_storage_account.test.name}"
  table_name           = "${azurerm_storage_table.test.name}"
  
  partition_key = "samplepartition"
  row_key       = "samplerow"

  entity = {
    sample = "entity"
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

## Import

Entities within a Table in an Azure Storage Account can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_table_entity.entity1 https://example.table.core.windows.net/table1(PartitionKey='samplepartition',RowKey='samplerow')
```


