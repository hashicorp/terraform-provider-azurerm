---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_mongo_database"
description: |-
  Gets information about an existing Cosmos DB Mongo Database.
---

# Data Source: azurerm_cosmosdb_mongo_database

Use this data source to access information about an existing Cosmos DB Mongo Database.

## Example Usage

```hcl
data "azurerm_cosmosdb_mongo_database" "example" {
  name                = "test-cosmosdb-mongo-db"
  resource_group_name = "test-cosmosdb-account-rg"
  account_name        = "test-cosmosdb-account"
}

output "id" {
  value = data.azurerm_cosmosdb_mongo_database.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `account_name` - (Required) The name of the Cosmos DB Account where the Mongo Database exists.

* `name` - (Required) The name of this Cosmos DB Mongo Database.

* `resource_group_name` - (Required) The name of the Resource Group where the Cosmos DB Mongo Database exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Cosmos DB Mongo Database.

* `tags` - A mapping of tags assigned to the Cosmos DB Mongo Database.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Cosmos Mongo Database.
