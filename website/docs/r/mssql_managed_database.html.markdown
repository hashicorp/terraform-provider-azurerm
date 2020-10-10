---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_managed_database"
description: |-
  Manages a MS SQL Managed Database.
---

# azurerm_mssql_managed_database

Manages a MS SQL Managed database.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "westus"
}

resource "azurerm_network_security_group" "example" {
  name                = "acceptanceTestSecurityGroup1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_virtual_network" "example" {
  name                = "sql_mi-network"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "internal"
  virtual_network_name = azurerm_virtual_network.example.name
  resource_group_name  = azurerm_resource_group.example.name
  address_prefixes     = ["10.0.1.0/24"]
  delegation {
    name = "miDelegation"

    service_delegation {
      name    = "Microsoft.Sql/managedInstances"
    }
  }
}

resource "azurerm_subnet_network_security_group_association" "example" {
  subnet_id                 = azurerm_subnet.example.id
  network_security_group_id = azurerm_network_security_group.example.id
}

resource "azurerm_route_table" "example" {
  name                = "example-routetable"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  route {
    name                   = "example"
    address_prefix         = "10.100.0.0/14"
    next_hop_type          = "VirtualAppliance"
    next_hop_in_ip_address = "10.10.1.1"
  }
}

resource "azurerm_subnet_route_table_association" "example" {
  subnet_id      = azurerm_subnet.example.id
  route_table_id = azurerm_route_table.example.id
}

resource "azurerm_mssql_managed_instance" "example" {
  name                 = "sql-mi"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  administrator_login = "demoReadUser"
  administrator_login_password = "ReadUser@123456"
  subnet_id = azurerm_subnet.example.id
  identity {
    type = "SystemAssigned"
  }
   sku {
        capacity = 8
        family = "Gen5"
        name = "GP_Gen5"
        tier = "GeneralPurpose"
      }
      license_type = "LicenseIncluded"
      collation =  "SQL_Latin1_General_CP1_CI_AS"
      proxy_override = "Redirect"
      storage_size_gb = 64
      vcores = 8
      public_data_endpoint_enabled = false
      timezone_id = "Central America Standard Time"
      minimal_tls_version = "1.2"
}

# Default create mode
resource "azurerm_sql_managed_database" "example" {
  name                 = "sql-managed-db"
  managed_instance_id = azurerm_mssql_managed_instance.example.id
  collation            = "SQL_Latin1_General_CP1_CI_AS"
  catalog_collation = "SQL_Latin1_General_CP1_CI_AS"
}

/* Point in time create mode. This block creates a managed database which is replica of an existing managed database to a specified point in time.
*/
resource "azurerm_sql_managed_database" "example" {
  name                 = "sql-managed-db"
  managed_instance_id = azurerm_mssql_managed_instance.example.id
  collation            = "SQL_Latin1_General_CP1_CI_AS"
  catalog_collation = "SQL_Latin1_General_CP1_CI_AS"
  create_mode = "PointInTimeRestore"
  source_database_id = "/subscriptions/00000000-1111-2222-3333-444444444444/resourceGroups/Default-SQL-SouthEastAsia/providers/Microsoft.Sql/managedInstances/testmi/databases/testdb"
  restore_point_in_time = "2017-07-14T05:35:31.503Z"
}

# Point in time create mode. This block creates a new db which is replica of a deleted managed database to a specified point in time.
resource "azurerm_sql_managed_database" "example" {
  name                 = "sql-managed-db"
  managed_instance_id = azurerm_mssql_managed_instance.example.id
  create_mode = "PointInTimeRestore"
  restorable_dropped_database_id = "/subscriptions/00000000-1111-2222-3333-444444444444/resourceGroups/Default-SQL-SouthEastAsia/providers/Microsoft.Sql/managedInstances/testmi/restorableDroppedDatabases/testdb"
  restore_point_in_time = "2017-07-14T05:35:31.503Z"
}

# Recovery create mode. This block creates a new db from a geo-replicated backup. 
resource "azurerm_sql_managed_database" "example" {
  name                 = "sql-managed-db"
  managed_instance_id = azurerm_mssql_managed_instance.example.id
  create_mode = "Recovery"
  recoverable_database_id = "/subscriptions/00000000-1111-2222-3333-444444444444/resourceGroups/Default-SQL-SouthEastAsia/providers/Microsoft.Sql/managedInstances/testmi/recoverableDatabases/testdb"
}

/* RestoreLongTermRetentionBackup create mode. This block creates a new db from long term retention backup(LTR backup policies should be applied first on the managed instance and there should be a LTR backup available).
*/
resource "azurerm_sql_managed_database" "example" {
  name                 = "sql-managed-db"
  managed_instance_id = azurerm_mssql_managed_instance.example.id
  create_mode = "RestoreLongTermRetentionBackup"
  longterm_retention_backup_id = "/subscriptions/00000000-1111-2222-3333-444444444444/providers/Microsoft.Sql/locations/japaneast/longTermRetentionManagedInstances/testInstance/longTermRetentionDatabases/testDatabase/longTermRetentionManagedInstanceBackups/55555555-6666-7777-8888-999999999999;131637960820000000"
}

/* RestoreExternalBackup create mode. This block creates a new db from a .bak file located in a storage account. The storage account SAS token should have Read and List permissions to access the .bak file and the .bak file should have checksum enabled.
*/
resource "azurerm_sql_managed_database" "example" {
  name                 = "sql-managed-db"
  managed_instance_id = azurerm_mssql_managed_instance.example.id
  create_mode = "RestoreExternalBackup"
  collation = "SQL_Latin1_General_CP1_CI_AS"
  storage_container_uri = "https://myaccountname.blob.core.windows.net/backups"
  storage_container_sas_token = "sv=2015-12-11&sr=c&sp=rl&sig=1234"
  last_backup_name = "mydb.bak"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the MS SQL Managed database. Changing this forces a new resource to be created.

* `managed_instance_id` - (Required) The Managed instance id to which the database belongs to. Changing this forces a new resource to be created.

* `collation` - (Optional) The collation of the Managed database. Changing this forces a new resource to be created. Defaults to `SQL_Latin1_General_CP1_CI_AS`.

* `restore_point_in_time` - (Optional) The restore point in time value if the create mode is `PointInTimeRestore`. Changing this forces a new resource to be created.

* `catalog_collation` - (Optional) Collation of the metadata catalog. Possible values include `DATABASE_DEFAULT` and `SQL_Latin1_General_CP1_CI_AS`.

* `create_mode` - (Optional) Specifies the mode of database creation. Possible values are `Default` (regular empty database creation),  `PointInTimeRestore` (creates a database by restoring a set of backups to specific point in time. `restore_point_in_time` and `source_database_id` must be specified to copy an existing database. `restore_point_in_time` and `restorable_dropped_database_id` must be specified to copy a deleted database), `Recovery` (creation from a geo-recover backup. `recoverable_database_id` must also be specified), `RestoreLongTermRetentionBackup`(creation from a LTR backup. `longterm_retention_backup_id` must be specified.) and `RestoreExternalBackup` (creation from a .bak file located in a storage account. `storage_container_uri`, `collation`, `storage_container_sas_token` and  `last_backup_name` must be cpecified. The .bak file should have checksum enabled). Changing this forces a new resource to be created.

* `storage_container_uri` - (Optional) Specifies the uri of the storage container where backups for this restore are stored. If `create_mode` is `RestoreExternalBackup`, this value is required along with `collation`, `storage_container_sas_token` and  `last_backup_name`. Changing this forces a new resource to be created.

* `source_database_id` - (Optional) The resource identifier of the source database associated with create operation of this database. If `create_mode` is `PointInTimeRestore`, this value is required along with `restore_point_in_time`. Changing this forces a new resource to be created.

* `restorable_dropped_database_id` - (Optional) The deleted database resource id to restore when creating this database. If `create_mode` is `PointInTimeRestore`, this value is required along with `restore_point_in_time`. Changing this forces a new resource to be created.

* `storage_container_sas_token` - (Optional) Specifies the uri of the storage container SAS token. If `create_mode` is `RestoreExternalBackup`, this value is required along with `collation`, `storage_container_uri` and  `last_backup_name`.

* `recoverable_database_id` - (Optional) The resource identifier of the geo replicated database backup associated with `Recovery` createMode operation. Changing this forces a new resource to be created.

* `last_backup_name` - (Optional) Specifies the .bak file name located in a storage account to restore from. If `create_mode` is `RestoreExternalBackup`, this value is required along with `collation`, `storage_container_sas_token` and  `storage_container_uri`. Changing this forces a new resource to be created.

* `longterm_retention_backup_id` - (Optional) TThe name of the Long Term Retention backup to be used for restore of this managed database. If `create_mode` is `RestoreLongTermRetentionBackup`, this value is required. Changing this forces a new resource to be created.

* `auto_complete_restore` - (Optional) whether to auto complete restore of this managed database.

* `tags` - (Optional) resource tags.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the managed database.

* `type` - The resource type of the managed database.

* `status` - The status of the managed database creation.

* `creation_date` - Managed database creation date

* `earliest_restore_point` - Earliest restore point in time of the created managed database

* `default_secondary_location` - The default secondary location of the managed database

* `failover_group_id` - The failover group id of the managed database

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 120 minutes) Used when creating the Managed database depending on the size of the backup for all create modes except Default. 
* `update` - (Defaults to 120 minutes) Used when updating the Managed database. 
* `read` - (Defaults to 5 minutes) Used when retrieving the Managed database.
* `delete` - (Defaults to 5 minutes) Used when deleting the Managed database.

## Import

Managed database can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sql_managed_database.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Sql/managedInstances/sql-mi/databases/sql-managed-db
```
