---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_database_long_term_retention_policy"
description: |-
  Manages Backup Long term retention policy for a SQL Database.
---

# azurerm_mssql_database_long_term_retention_policy

Allows you to manage Backup Long term retention policy for Azure SQL Database

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

resource "azurerm_mssql_database_long_term_retention_policy" "example" {
  database_name       = azurerm_sql_database.example.name
  resource_group_name = azurerm_resource_group.example.name
  server_name         = azurerm_sql_server.example.name

  weekly_retention  = "P1W"
  monthly_retention = "P1M"
  yearly_retention  = "P1Y"
  week_of_year      = 1
}
```

## Argument Reference

The following arguments are supported:

* `database_name` - (Required) The name of the database.

* `resource_group_name` - (Required) The name of the resource group in which database resides.  This must be the same as Database Server resource group currently.

* `server_name` - (Required) The name of the associated SQL Server.

* `weekly_retention` - (Optional) The weekly retention policy for an LTR backup in an ISO 8601 format. Valid value is between 1 to 520 weeks. e.g. `P1Y`, `P1M`, `P1W` or `P7D`.

* `monthly_retention` - (Optional) The monthly retention policy for an LTR backup in an ISO 8601 format. Valid value is between 1 to 120 months. e.g. `P1Y`, `P1M`, `P4W` or `P30D`.

* `yearly_retention` - (Optional) The yearly retention policy for an LTR backup in an ISO 8601 format. Valid value is between 1 to 10 years. e.g. `P1Y`, `P12M`, `P52W` or `P365D`.

* `week_of_year` - (Optional) The week of year to take the yearly backup in an ISO 8601 format. Value has to be between `1` and `52`.

## Attributes Reference

The following attributes are exported:

* `id` - Long term policy ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating Long Term retention policy for SQL Database.
* `update` - (Defaults to 5 minutes) Used when updating Long Term retention policy for SQL Database.
* `read` - (Defaults to 5 minutes) Used when retrieving Long Term retention policy for SQL Database.
* `delete` - (Defaults to 5 minutes) Used when deleting Long Term retention policy for SQL Database.

## Import

Long Term retention policy can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_database_long_term_retention_policy.policy1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Sql/servers/myserver/databases/database1/backupLongTermRetentionPolicies/default
```
