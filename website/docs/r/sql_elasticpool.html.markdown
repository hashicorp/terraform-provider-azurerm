---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sql_elasticpool"
description: |-
  Manages a SQL Elastic Pool.
---

# azurerm_sql_elasticpool

Allows you to manage an Azure SQL Elastic Pool.

~> **NOTE:** -  This version of the `Elasticpool` resource is being **deprecated** and should no longer be used. Please use the [azurerm_mssql_elasticpool](./mssql_elasticpool.html) version instead.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "my-resource-group"
  location = "West Europe"
}

resource "azurerm_sql_server" "example" {
  name                         = "my-sql-server" # NOTE: needs to be globally unique
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_sql_elasticpool" "example" {
  name                = "test"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  server_name         = azurerm_sql_server.example.name
  edition             = "Basic"
  dtu                 = 50
  db_dtu_min          = 0
  db_dtu_max          = 5
  pool_size           = 5000
}
```

~> **NOTE on `azurerm_sql_elasticpool`:** -  The values of `edition`, `dtu`, and `pool_size` must be consistent with the [Azure SQL Database Service Tiers](https://docs.microsoft.com/en-gb/azure/sql-database/sql-database-service-tiers#elastic-pool-service-tiers-and-performance-in-edtus). Any inconsistent argument configuration will be rejected.

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the elastic pool. This needs to be globally unique. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the elastic pool. This must be the same as the resource group of the underlying SQL server.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `server_name` - (Required) The name of the SQL Server on which to create the elastic pool. Changing this forces a new resource to be created.

* `edition` - (Required) The edition of the elastic pool to be created. Valid values are `Basic`, `Standard`, and `Premium`. Refer to [Azure SQL Database Service Tiers](https://docs.microsoft.com/en-gb/azure/sql-database/sql-database-service-tiers#elastic-pool-service-tiers-and-performance-in-edtus) for details. Changing this forces a new resource to be created.

* `dtu` - (Required) The total shared DTU for the elastic pool. Valid values depend on the `edition` which has been defined. Refer to [Azure SQL Database Service Tiers](https://docs.microsoft.com/en-gb/azure/sql-database/sql-database-service-tiers#elastic-pool-service-tiers-and-performance-in-edtus) for valid combinations.

* `db_dtu_min` - (Optional) The minimum DTU which will be guaranteed to all databases in the elastic pool to be created.

* `db_dtu_max` - (Optional) The maximum DTU which will be guaranteed to all databases in the elastic pool to be created.

* `pool_size` - (Optional) The maximum size in MB that all databases in the elastic pool can grow to. The maximum size must be consistent with combination of `edition` and `dtu` and the limits documented in [Azure SQL Database Service Tiers](https://docs.microsoft.com/en-gb/azure/sql-database/sql-database-service-tiers#elastic-pool-service-tiers-and-performance-in-edtus). If not defined when creating an elastic pool, the value is set to the size implied by `edition` and `dtu`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The SQL Elastic Pool ID.

* `creation_date` - The creation date of the SQL Elastic Pool.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the SQL Elastic Pool.
* `update` - (Defaults to 30 minutes) Used when updating the SQL Elastic Pool.
* `read` - (Defaults to 5 minutes) Used when retrieving the SQL Elastic Pool.
* `delete` - (Defaults to 30 minutes) Used when deleting the SQL Elastic Pool.

## Import

SQL Elastic Pool's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sql_elasticpool.pool1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Sql/servers/myserver/elasticPools/pool1
```
