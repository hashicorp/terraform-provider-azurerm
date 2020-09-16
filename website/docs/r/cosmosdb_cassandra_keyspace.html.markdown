---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_cassandra_keyspace"
description: |-
  Manages a Cassandra KeySpace within a Cosmos DB Account.
---

# azurerm_cosmosdb_cassandra_keyspace

Manages a Cassandra KeySpace within a Cosmos DB Account.

## Example Usage

```hcl
data "azurerm_resource_group" "example" {
  name = "tflex-cosmosdb-account-rg"
}

resource "azurerm_cosmosdb_account" "example" {
  name                = "tfex-cosmosdb-account"
  resource_group_name = data.azurerm_resource_group.example.name
  location            = data.azurerm_resource_group.example.location
  offer_type          = "Standard"

  capabilities {
    name = "EnableCassandra"
  }

  consistency_policy {
    consistency_level = "Strong"
  }

  geo_location {
    location          = "West US"
    failover_priority = 0
  }
}

resource "azurerm_cosmosdb_cassandra_keyspace" "example" {
  name                = "tfex-cosmos-cassandra-keyspace"
  resource_group_name = data.azurerm_cosmosdb_account.example.resource_group_name
  account_name        = azurerm_cosmosdb_account.example.name
  throughput          = 400
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Cosmos DB Cassandra KeySpace. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Cosmos DB Cassandra KeySpace is created. Changing this forces a new resource to be created.

* `account_name` - (Required) The name of the Cosmos DB Cassandra KeySpace to create the table within. Changing this forces a new resource to be created.

* `throughput` - (Optional) The throughput of Cassandra KeySpace (RU/s). Must be set in increments of `100`. The minimum value is `400`. This must be set upon database creation otherwise it cannot be updated without a manual terraform destroy-apply.

* `autoscale_settings` - (Optional) An `autoscale_settings` block as defined below. This must be set upon database creation otherwise it cannot be updated without a manual terraform destroy-apply.

~> **Note:** Switching between autoscale and manual throughput is not supported via Terraform and must be completed via the Azure Portal and refreshed. 

---

An `autoscale_settings` block supports the following:

* `max_throughput` - (Optional) The maximum throughput of the Cassandra KeySpace (RU/s). Must be between `4,000` and `1,000,000`. Must be set in increments of `1,000`. Conflicts with `throughput`.


## Attributes Reference

The following attributes are exported:

* `id` - the ID of the CosmosDB Cassandra KeySpace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the CosmosDB Cassandra KeySpace.
* `update` - (Defaults to 30 minutes) Used when updating the CosmosDB Cassandra KeySpace.
* `read` - (Defaults to 5 minutes) Used when retrieving the CosmosDB Cassandra KeySpace.
* `delete` - (Defaults to 30 minutes) Used when deleting the CosmosDB Cassandra KeySpace.

## Import

Cosmos Cassandra KeySpace can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cosmosdb_cassandra_keyspace.ks1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.DocumentDB/databaseAccounts/account1/cassandraKeyspaces/ks1
```
