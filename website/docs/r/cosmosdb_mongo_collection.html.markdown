---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_mongo_collection"
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
  resource_group_name = data.azurerm_cosmosdb_account.example.resource_group_name
  account_name        = data.azurerm_cosmosdb_account.example.name
}

resource "azurerm_cosmosdb_mongo_collection" "example" {
  name                = "tfex-cosmos-mongo-db"
  resource_group_name = data.azurerm_cosmosdb_account.example.resource_group_name
  account_name        = data.azurerm_cosmosdb_account.example.name
  database_name       = azurerm_cosmosdb_mongo_database.example.name

  default_ttl_seconds = "777"
  shard_key           = "uniqueKey"
  throughput          = 400
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Cosmos DB Mongo Collection. Changing this forces a new resource to be created.
* `resource_group_name` - (Required) The name of the resource group in which the Cosmos DB Mongo Collection is created. Changing this forces a new resource to be created.
* `database_name` - (Required) The name of the Cosmos DB Mongo Database in which the Cosmos DB Mongo Collection is created. Changing this forces a new resource to be created.
* `default_ttl_seconds` - (Required) The default Time To Live in seconds. If the value is `-1` or `0`, items are not automatically expired.
* `shard_key` - (Required) The name of the key to partition on for sharding. There must not be any other unique index keys.
* `analytical_storage_ttl` - (Optional) The default time to live of Analytical Storage for this Mongo Collection. If present and the value is set to `-1`, it is equal to infinity, and items don’t expire by default. If present and the value is set to some number `n` – items will expire `n` seconds after their last modified time.
* `index` - (Optional) One or more `index` blocks as defined below.
* `throughput` - (Optional) The throughput of the MongoDB collection (RU/s). Must be set in increments of `100`. The minimum value is `400`. This must be set upon database creation otherwise it cannot be updated without a manual terraform destroy-apply.
* `autoscale_settings` - (Optional) An `autoscale_settings` block as defined below. This must be set upon database creation otherwise it cannot be updated without a manual terraform destroy-apply. Requires `shard_key` to be set.

~> **Note:** Switching between autoscale and manual throughput is not supported via Terraform and must be completed via the Azure Portal and refreshed. 

---

An `autoscale_settings` block supports the following:

* `max_throughput` - (Optional) The maximum throughput of the MongoDB collection (RU/s). Must be between `4,000` and `1,000,000`. Must be set in increments of `1,000`. Conflicts with `throughput`.

---

The `index` block supports the following:

* `keys` - (Required) Specifies the list of user settable keys for each Cosmos DB Mongo Collection.

* `unique` - (Optional) Is the index unique or not? Defaults to `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Cosmos DB Mongo Collection.

* `system_indexes` - One or more `system_indexes` blocks as defined below.

---

The `system_indexes` block supports the following:

* `keys` - The list of system keys which are not settable for each Cosmos DB Mongo Collection.

* `unique` - Identifies whether the table contains no duplicate values.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the CosmosDB Mongo Collection.
* `update` - (Defaults to 30 minutes) Used when updating the CosmosDB Mongo Collection.
* `read` - (Defaults to 5 minutes) Used when retrieving the CosmosDB Mongo Collection.
* `delete` - (Defaults to 30 minutes) Used when deleting the CosmosDB Mongo Collection.

## Import

CosmosDB Mongo Collection can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cosmosdb_mongo_collection.collection1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.DocumentDB/databaseAccounts/account1/mongodbDatabases/db1/collections/collection1
```
