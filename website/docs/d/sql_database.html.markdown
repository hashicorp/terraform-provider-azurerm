---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sql_database"
description: |-
  Gets information about an existing SQL Azure Database.
---

# Data Source: azurerm_sql_database

Use this data source to access information about an existing SQL Azure Database.

## Example Usage

```hcl
data "azurerm_sql_database" "example" {
  name                = "example_db"
  server_name         = "example_db_server"
  resource_group_name = "example-resources"
}

output "sql_database_id" {
  value = data.azurerm_sql_database.example.id
}
```

## Argument Reference

* `name` - The name of the SQL Database.

* `server_name` - The name of the SQL Server.

* `resource_group_name` - Specifies the name of the Resource Group where the Azure SQL Database exists.

## Attributes Reference

* `id` - The SQL Database ID.

* `collation` - The name of the collation.

* `creation_date` - The creation date of the SQL Database.

* `default_secondary_location` - The default secondary location of the SQL Database.

* `edition` - The edition of the database.

* `elastic_pool_name` - The name of the elastic database pool the database belongs to.

* `failover_group_id` - The ID of the failover group the database belongs to.

* `location` - The location of the Resource Group in which the SQL Server exists.

* `name` - The name of the database.

* `read_scale` - Indicate if read-only connections will be redirected to a high-available replica.

* `requested_service_objective_id` - The ID pertaining to the performance level of the database.

* `requested_service_objective_name` - The name pertaining to the performance level of the database. 

* `resource_group_name` - The name of the resource group in which the database resides. This will always be the same resource group as the Database Server.

* `server_name` - The name of the SQL Server on which to create the database.

* `tags` - A mapping of tags assigned to the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the SQL Azure Database.
