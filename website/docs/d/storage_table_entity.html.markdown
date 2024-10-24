---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_table_entity"
description: |-
  Gets information about an existing Storage Table Entity.
---

# Data Source: azurerm_storage_table_entity

Use this data source to access information about an existing Storage Table Entity.

## Example Usage

```hcl
data "azurerm_storage_table_entity" "example" {
  storage_table_id = azurerm_storage_table.example.id
  partition_key    = "example-partition-key"
  row_key          = "example-row-key"
}
```

## Argument Reference

The following arguments are supported:

* `storage_table_id` - (Required) The Storage Table ID where the entity exists.

* `partition_key` - The key for the partition where the entity will be retrieved.

* `row_key` - The key for the row where the entity will be retrieved.

## Attributes Reference

* `id` - The ID of the storage table entity.

* `entity` - A map of key/value pairs that describe the entity to be stored in the storage table.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Table Entity.
