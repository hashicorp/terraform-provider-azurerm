---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_database"
description: |-
  Manages a MS SQL Database.
---

# azurerm_mssql_database

Manages a MS SQL Database.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_sql_server" "example" {
  name                         = "example-sqlserver"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}
resource "azurerm_mssql_database" "test" {
  name           = "acctest-db-%d"
  sql_server_id  = azurerm_sql_server.test.id
  collation      = "SQL_Latin1_General_CP1_CI_AS"
  license_type   = "LicenseIncluded"
  max_size_gb    = 4
  read_scale     = "Enabled"
  sample_name    = "AdventureWorksLT"
  sku_name       = "BC_Gen5_2"
  zone_redundant = true

  tags = {
    foo = "bar"
  }

}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Ms SQL Database. Changing this forces a new resource to be created.

* `sql_server_id` - (Required) The id of the Ms SQL Server on which to create the database. Changing this forces a new resource to be created.

* `auto_pause_delay_in_minutes` - (Optional) Time in minutes after which database is automatically paused. A value of -1 means that automatic pause is disabled. 

* `create_mode` - (Optional) The create mode of the database. Possible values are `Copy`, `Default`, `OnlineSecondary`, `PointInTimeRestore`, `Restore`, `RestoreExternalBackup`, `RestoreExternalBackupSecondary`, `RestoreLongTermRetentionBackup` and `Secondary`. 

* `collation` - (Optional) Specifies the collation of the database. Changing this forces a new resource to be created.

* `elastic_pool_id` - (Optional) Specifies the id of the elastic pool containing this database. Changing this forces a new resource to be created.

* `license_type` - (Optional) Specifies the license type to apply for this database. Possible values are `LicenseIncluded` and `BasePrice`.

* `max_size_gb` - (Optional) The max size of the database in gigabytes.

* `min_capacity` - (Optional) Minimal capacity that database will always have allocated, if not paused.

* `restore_point_in_time` - (Required) Specifies the point in time (ISO8601 format) of the source database that will be restored to create the new database.

* `read_replica_count` - (Optional) The number of readonly secondary replicas associated with the database to which readonly application intent connections may be routed. 

* `read_scale` - (Optional) If enabled, connections that have application intent set to readonly in their connection string may be routed to a readonly secondary replica. Possible values are `Enabled` and `Disabled`.

* `sample_name` - (Optional) Specifies the name of the sample schema to apply when creating this database. Possible value is `AdventureWorksLT`.

* `sku_name` - (Optional) Specifies the name of the sku of the database. Changing this forces a new resource to be created. For example, `GP_Gen5_2`,`HS_Gen4_1`,`BC_Gen5_2`, `ElasticPool`, `Basic`, `Standard`, `Premium`, `S0`, `P2`.

* `source_database_id` - (Required) The id of the source database to be copied to create the new database.

* `zone_redundant` - (Optional) Whether or not this database is zone redundant, which means the replicas of this database will be spread across multiple availability zones.

* `tags` - (Optional) A mapping of tags to assign to the resource.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the MS SQL Database.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the MS SQL Database.
* `update` - (Defaults to 30 minutes) Used when updating the MS SQL Database.
* `read` - (Defaults to 5 minutes) Used when retrieving the MS SQL Database.
* `delete` - (Defaults to 30 minutes) Used when deleting the MS SQL Database.

## Import

SQL Database can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_database.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/databases/example1
```
