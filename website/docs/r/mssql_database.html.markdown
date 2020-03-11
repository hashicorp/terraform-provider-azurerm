---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_database"
description: |-
  Manages a SQL Database.
---

# azurerm_mssql_database

Allows you to manage an Azure SQL Database via the `2017-10-01-preview` API which allows for `vCore` based configurations.

## Example Usage

```hcl
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

resource "azurerm_mssql_database" "example" {
  name            = "example-mssql-db"
  mssql_server_id = azurerm_sql_server.example.id
  collation       = "SQL_Latin1_General_CP1_CI_AS"
  license_type    = "LicenseIncluded"
  sample_name     = "AdventureWorksLT"

  general_purpose {
    capacity    = 2
    family      = "Gen5"
    max_size_gb = 4
    serverless {
      auto_pause_delay_in_minutes = 60
      min_capacity                = 0.5
    }
  }

  tags = {
    foo = "bar"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Ms SQL Database. Changing this forces a new resource to be created.

* `mssql_server_id` - (Required) The name of the Ms SQL Server on which to create the database. Changing this forces a new resource to be created.

* `business_critical` - (Optional) A `business_critical` block as defined below. Conflicts with `general_purpose` and `hyper_scale`.

* `collation` - (Optional) Specifies the collation of the database. Changing this forces a new resource to be created.

* `create_copy_mode` - (Optional) A `create_copy_mode` block as defined below. Conflicts with `create_pitr_mode` and `create_secondary_mode`.

* `create_pitr_mode` - (Optional) A `create_pitr_mode` block as defined below. Conflicts with `create_copy_mode` and `create_secondary_mode`.

* `create_secondary_mode` - (Optional) A `create_secondary_mode` block as defined below. Conflicts with `create_copy_mode` and `create_pitr_mode`.

* `elastic_pool_id` - (Optional) Specifies the id of the elastic pool containing this database. Changing this forces a new resource to be created.

* `general_purpose` - (Optional) A `general_purpose` block as defined below. Conflicts with `business_critical` and `hyper_scale`.

* `hyper_scale` - (Optional) A `hyper_scale` block as defined below. Changing this forces a new resource to be created. Conflicts with `business_critical` and `general_purpose`.

* `license_type` - (Optional) Specifies the license type to apply for this database. Possible values are `LicenseIncluded` and `BasePrice`.

* `sample_name` - (Optional) Specifies the name of the sample schema to apply when creating this database. Possible value is `AdventureWorksLT`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---
`business_critical` supports the following:

* `capacity` - (Required) The scale up/out capacity, representing server's compute units. For more information see the documentation for your Database configuration: [vCore-based](https://docs.microsoft.com/en-us/azure/sql-database/sql-database-service-tiers-vcore).

* `family` - (Required) The `family` of hardware `Gen4` or `Gen5`.

* `max_size_gb` - (Optional) The max size of the database in gigabytes.

* `read_scale` - (Optional) If enabled, connections that have application intent set to readonly in their connection string may be routed to a readonly secondary replica. Possible values are `Enabled` and `Disabled`.

* `zone_redundant` - (Optional) Whether or not this database is zone redundant, which means the replicas of this database will be spread across multiple availability zones.

---
`create_copy_mode` supports the following:

* `source_database_id` - (Required) The id of the source database to be copied to create the new database.

---
`create_pitr_mode` supports the following:

* `source_database_id` - (Required) The id of the source database to be restored to create the new database.

* `restore_point_in_time` - (Required) Specifies the point in time (ISO8601 format) of the source database that will be restored to create the new database.

---
`create_secondary_mode` supports the following:

* `source_database_id` - (Required) The id of the source database to be copied to create the new database in another location.

---
`general_purpose` supports the following:

* `capacity` - (Required) The scale up/out capacity, representing server's compute units. For more information see the documentation for your Database configuration: [vCore-based](https://docs.microsoft.com/en-us/azure/sql-database/sql-database-service-tiers-vcore).

* `family` - (Required) The `family` of hardware `Gen4` or `Gen5`.

* `max_size_gb` - (Optional) The max size of the database in gigabytes.

* `serverless` - (Optional) A `serverless` block as defined below.

---
`hyper_scale` supports the following:

* `capacity` - (Required) The scale up/out capacity, representing server's compute units. For more information see the documentation for your Database configuration: [vCore-based](https://docs.microsoft.com/en-us/azure/sql-database/sql-database-service-tiers-vcore).

* `family` - (Required) The `family` of hardware `Gen4` or `Gen5`.

* `read_replica_count` - (Optional) The number of readonly secondary replicas associated with the database to which readonly application intent connections may be routed. 

---
`serverless` supports the following:

* `auto_pause_delay_in_minutes` - (Optional) Time in minutes after which database is automatically paused. A value of -1 means that automatic pause is disabled. 

* `min_capacity` - (Optional) Minimal capacity that database will always have allocated, if not paused.

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
