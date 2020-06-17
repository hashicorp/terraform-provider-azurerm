---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_database_short_term_retention_policy"
description: |-
  Manages Backup Short term retention policy for a SQL Database.
---

# azurerm_mssql_database_short_term_retention_policy

Allows you to manage Backup Short term retention policy for Azure SQL Database

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "acceptanceTestResourceGroup1"
  location = "uksouth"
}

resource "azurerm_sql_server" "example" {
  name                         = "testsqlserver"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = "uksouth"
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_sql_database" "example" {
  name                = "mysqldatabase"
  resource_group_name = azurerm_resource_group.example.name
  location            = "uksouth"
  server_name         = azurerm_sql_server.example.name

  tags = {
    environment = "production"
  }
}

resource "azurerm_mssql_database_short_term_retention_policy" "example" {
  database_name       = azurerm_sql_database.example.name
  resource_group_name = azurerm_resource_group.example.name
  server_name         = azurerm_sql_server.example.name

  retention_days = 7
}
```

## Argument Reference

The following arguments are supported:

* `database_name` - (Required) The name of the database.

* `resource_group_name` - (Required) The name of the resource group in which database resides.  This must be the same as Database Server resource group currently.

* `server_name` - (Required) The name of the associated SQL Server.

* `retention_days` - (Required) Point In Time Restore configuration. Value has to be between `7` and `35`.

## Attributes Reference

The following attributes are exported:

* `id` - Short term policy ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating Short Term retention policy for SQL Database.
* `update` - (Defaults to 5 minutes) Used when updating Short Term retention policy for SQL Database.
* `read` - (Defaults to 5 minutes) Used when retrieving Short Term retention policy for SQL Database.
* `delete` - (Defaults to 5 minutes) Used when deleting Short Term retention policy for SQL Database.

## Import

Short Term retention policy can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_database_short_term_retention_policy.policy1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Sql/servers/myserver/databases/database1/backupShortTermRetentionPolicies/default
```
