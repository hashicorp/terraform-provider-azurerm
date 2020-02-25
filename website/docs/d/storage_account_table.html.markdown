---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_table"
description: |-
  Fetch Azure Table records stored under an existing Storage Account.

---

# Data Source: azurerm_storage_table

Use this data source to fetch data from Azure Table stored under an existing Storage Account.

Data source allows to "select" data in three ways and always returns a collection.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "westus"
}

resource "azurerm_storage_account" "example" {
  name                     = "example"
  location                 = azurerm_resource_group.example.location
  resource_group_name      = azurerm_resource_group.example.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_table" "cats" {
  name                 = "cats"
  storage_account_name = azurerm_storage_account.example.name
}

resource "azurerm_storage_table_entity" "cat1" {
  storage_account_name = azurerm_storage_account.example.name
  table_name           = azurerm_storage_table.cats.name
  partition_key        = "mypartition"
  row_key              = "1"

  entity = {
    breed     = "Scottish Fold"
    life_span = "11 - 14"
    img_url   = "https://cdn2.thecatapi.com/images/4UTnt4f74.jpg"
  }
}

resource "azurerm_storage_table_entity" "cat2" {
  storage_account_name = azurerm_storage_account.example.name
  table_name           = azurerm_storage_table.cats.name
  partition_key        = "mypartition"
  row_key              = "2"

  entity = {
    breed     = "Munchkin"
    life_span = "10 - 15"
    img_url   = "https://cdn2.thecatapi.com/images/hxlto6Z4I.jpg"
  }
}

data "azurerm_storage_table" "by_resource_id" {
  resource_id = azurerm_storage_table_entity.cat1.id
}

data "azurerm_storage_table" "by_key" {
  key {
    storage_account_name = azurerm_storage_account.example.name
    table_name           = azurerm_storage_table.cats.name
    partition_key        = "mypartition"
    row_key              = "2"
  }
}

data "azurerm_storage_table" "by_query" {
  key {
    storage_account_name = azurerm_storage_account.example.name
    table_name           = azurerm_storage_table.cats.name
    filter               = "PartitionKey eq 'mypartition'"
    top                  = 10
  }
}

output "data_storage_table_by_resource_id" {
  value = data.azurerm_storage_table.by_resource_id.entities
}
output "data_storage_table_by_key" {
  value = data.azurerm_storage_table.by_key.entities
}
output "data_storage_table_by_query" {
  value = data.azurerm_storage_table.by_query.entities
}
```

## Argument Reference

To fetch entity (entities) from Azure Table you must specify at-least one of the following configuration arguments:
`resource_id`, `key` (block), `query` (block). Please note that you cannot use two or more different configuration arguments
at the same time.

* `resource_id` - The ID string that will be used to fetch Table entity

---

* `key` is a set of parameters that are used to fetch single Table entity. A `key` block contains:

  * `storage_account_name` - (Required) The name of Storage Account holding Azure Table

  * `table_name` - (Required) The name of Azure Table

  * `partition_key` - (Required) The key for the partition where the Table entity is stored

  * `row_key` - (Required) The key for the row where the Table entity is stored

---

* `query` is a set of parameters that are used to fetch collection of Table entities. The collection is filtered by `filter` parameter. A `query` block contains:

  * `storage_account_name` - (Required) The name of Storage Account holding Azure Table

  * `table_name` - (Required) The name of Azure Table

  * `filter` - (Required) Filter configuration that will be used to configure query

  * `top` - (Optional) Integer to set maximum number of returned Table entities

Refer to the [Supported Query Options](https://docs.microsoft.com/en-us/rest/api/storageservices/querying-tables-and-entities#supported-query-options)
on syntax of `filter`.

## Attributes Reference

The following attributes are exported:

* `entities` - The array of Table entities fetched by data source.

* `id` - The ID of the Entity within the Table in the Storage Account.
