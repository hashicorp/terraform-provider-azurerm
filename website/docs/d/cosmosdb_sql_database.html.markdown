---
subcategory: "CosmosDB (DocumentDB)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cosmosdb_sql_database"
description: |-
  Gets information about an existing CosmosDB SQL Database.
---

# Data Source: azurerm_cosmosdb_sql_database

Use this data source to access information about an existing CosmosDB SQL Database.

## Example Usage

```hcl
data "azurerm_cosmosdb_sql_database" "example" {
  name                = "tfex-cosmosdb-sql-database"
  resource_group_name = "tfex-cosmosdb-sql-database-rg"
  account_name        = "tfex-cosmosdb-sql-database-account-name"
}

```

## Argument Reference

The following arguments are supported:

* `name` - Specifies the name of the Cosmos DB SQL Database.

* `resource_group_name` - The name of the resource group in which the Cosmos DB SQL Database is created.

* `account_name` - The name of the Cosmos DB SQL Database to create the table within.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the CosmosDB SQL Database.

* `throughput` -The throughput of SQL database (RU/s).

* `autoscale_settings` - An `autoscale_settings` block as defined below.

---

An `autoscale_settings` block supports the following:

* `max_throughput` - The maximum throughput of the SQL database (RU/s).

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the CosmosDB SQL Database.
