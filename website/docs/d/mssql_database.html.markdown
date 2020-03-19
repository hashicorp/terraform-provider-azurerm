---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_database"
description: |-
  Gets information about an existing SQL database.
---

# Data Source: azurerm_mssql_database

Use this data source to access information about an existing SQL database.

## Example Usage

```hcl
data "azurerm_mssql_database" "example" {
  name      = "example-mssql-db"
  server_id = "example-mssql-server-id"
}

output "database_id" {
  value = data.azurerm_mssql_database.example.id
}
```

## Argument Reference

* `name` - The name of the Ms SQL Database.

* `server_id` - The id of the Ms SQL Server on which to create the database.

## Attribute Reference

* `collation` - The collation of the database. 

* `elastic_pool_id` - The id of the elastic pool containing this database.

* `license_type` - The license type to apply for this database.

* `max_size_gb` - The max size of the database in gigabytes.

* `read_replica_count` - The number of readonly secondary replicas associated with the database to which readonly application intent connections may be routed. 

* `read_scale` - If enabled, connections that have application intent set to readonly in their connection string may be routed to a readonly secondary replica.

* `sku_name` - The name of the sku of the database.

* `zone_redundant` - Whether or not this database is zone redundant, which means the replicas of this database will be spread across multiple availability zones.

* `tags` -  A mapping of tags to assign to the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the SQL database.
