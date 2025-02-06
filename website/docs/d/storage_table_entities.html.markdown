---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_table_entities"
description: |-
  Gets all existing entities from Storage Tablethat match a filter.
---

# Data Source: azurerm_storage_table_entity

Use this data source to access information about an existing Storage Table Entity.

## Example Usage

```hcl
data "azurerm_storage_table_entities" "example" {
  storage_table_id = azurerm_storage_table.example.id
  filter           = "PartitionKey eq 'example'"
}
```

## Argument Reference

The following arguments are supported:

* `storage_table_id` - The Storage Table ID where the entities exist.

* `filter` - The filter used to retrieve the entities.

* `select` - (Optional) A list of properties to select from the returned Storage Table Entities.

## Attributes Reference

* `id` - The ID of the storage table entity.

* `items` - A list of `items` blocks as defined below.

---

Each element in `items` block exports the following:

* `partition_key` - Partition Key of the Entity.

* `row_key` - Row Key of the Entity.

* `properties` - A map of any additional properties in key-value format.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Table Entity.
