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
  name            = "example-mssql-db"
  mssql_server_id = "example-mssql-server-id"
}

output "database_id" {
  value = data.azurerm_mssql_database.example.id
}
```

## Argument Reference

* `name` - The name of the Ms SQL Database.

* `mssql_server_id` - The id of the Ms SQL Server on which to create the database.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the SQL database.
