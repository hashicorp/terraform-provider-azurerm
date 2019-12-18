---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_mongo_collection"
sidebar_current: "docs-azurerm-resource-cosmosdb-mongo-collection"
description: |-
  Manages a Mongo Collection within a Cosmos DB Account.
---

# azurerm_cosmosdb_mongo_collection

Manages a Mongo Collection within a Cosmos DB Account.

## Example Usage

```hcl
data "azurerm_cosmosdb_account" "example" {
  name                = "tfex-cosmosdb-account"
  resource_group_name = "tfex-cosmosdb-account-rg"
}

resource "azurerm_cosmosdb_mongo_database" "example" {
  name                = "tfex-cosmos-mongo-db"
  resource_group_name = "${data.azurerm_cosmosdb_account.example.resource_group_name}"
  account_name        = "${data.azurerm_cosmosdb_account.example.name}"
}

resource "azurerm_cosmosdb_mongo_collection" "example" {
  name                = "tfex-cosmos-mongo-db"
  resource_group_name = "${data.azurerm_cosmosdb_account.example.resource_group_name}"
  account_name        = "${data.azurerm_cosmosdb_account.example.name}"
  database_name       = "${azurerm_cosmosdb_mongo_database.example.name}"

  default_ttl_seconds = "777"
  shard_key           = "uniqueKey"
  throughput          = 400

  indexes {
    key    = "aKey"
    unique = false
  }

  indexes {
    key    = "uniqueKey"
    unique = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Cosmos DB Mongo Collection. Changing this forces a new resource to be created.
* `resource_group_name` - (Required) The name of the resource group in which the Cosmos DB Mongo Collection is created. Changing this forces a new resource to be created.
* `database_name` - (Required) The name of the Cosmos DB Mongo Database in which the Cosmos DB Mongo Collection is created. Changing this forces a new resource to be created.
* `default_ttl_seconds` - (Required) The default Time To Live in seconds. If the value is `-1` items are not automatically expired.
* `shard_key` - (Required) The name of the key to partition on for sharding. There must not be any other unique index keys.
* `throughput` - (Optional) The throughput of the MongoDB collection (RU/s). Must be set in increments of `100`. The minimum value is `400`. This must be set upon database creation otherwise it cannot be updated without a manual terraform destroy-apply.
* `indexes` - (Optional) One or more `indexes` blocks as defined below.

---

An `indexes` block supports the following:

* `key` - (Required) The name of the key to use for this index.
* `unique` - (Required) Whether the index key should be unique.

## Attributes Reference

The following attributes are exported:

* `id` - the Cosmos DB Mongo Collection ID.

## Import

Cosmos DB Mongo Collection can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cosmosdb_mongo_collection.collection1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.DocumentDB/databaseAccounts/account1/apis/mongodb/databases/db1/collections/collection1
```
