---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sql_database"
description: |-
  Manages a SQL Database.
---

# azurerm_sql_database

Allows you to manage an Azure SQL Database

!>**IMPORTANT:** To mitigate the possibility of accidental data loss it is highly recommended that you use the `prevent_destroy` lifecycle argument in your configuration file for this resource. For more information on the `prevent_destroy` lifecycle argument please see the [terraform documentation](https://developer.hashicorp.com/terraform/tutorials/state/resource-lifecycle#prevent-resource-deletion).

->**NOTE:** The `azurerm_sql_database` resource is deprecated in version 3.0 of the AzureRM provider and will be removed in version 4.0. Please use the [`azurerm_mssql_database`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/mssql_database) resource instead.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_sql_server" "example" {
  name                         = "myexamplesqlserver"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"

  tags = {
    environment = "production"
  }
}

resource "azurerm_storage_account" "example" {
  name                     = "examplesa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_sql_database" "example" {
  name                = "myexamplesqldatabase"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  server_name         = azurerm_sql_server.example.name

  tags = {
    environment = "production"
  }

  # prevent the possibility of accidental data loss
  lifecycle {
    prevent_destroy = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the database. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the database. This must be the same as Database Server resource group currently. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `server_name` - (Required) The name of the SQL Server on which to create the database. Changing this forces a new resource to be created.

* `create_mode` - (Optional) Specifies how to create the database. Valid values are: `Default`, `Copy`, `OnlineSecondary`, `NonReadableSecondary`, `PointInTimeRestore`, `Recovery`, `Restore` or `RestoreLongTermRetentionBackup`. Must be `Default` to create a new database. Defaults to `Default`. Please see [Azure SQL Database REST API](https://docs.microsoft.com/rest/api/sql/databases/createorupdate#createmode)

* `import` - (Optional) A `import` block as documented below. `create_mode` must be set to `Default`.

* `source_database_id` - (Optional) The URI of the source database if `create_mode` value is not `Default`.

* `restore_point_in_time` - (Optional) The point in time for the restore. Only applies if `create_mode` is `PointInTimeRestore`, it should be provided in [RFC3339](https://www.rfc-editor.org/rfc/rfc3339) format, e.g. `2013-11-08T22:00:40Z`.

* `edition` - (Optional) The edition of the database to be created. Applies only if `create_mode` is `Default`. Valid values are: `Basic`, `Standard`, `Premium`, `DataWarehouse`, `Business`, `BusinessCritical`, `Free`, `GeneralPurpose`, `Hyperscale`, `Premium`, `PremiumRS`, `Standard`, `Stretch`, `System`, `System2`, or `Web`. Please see [Azure SQL database models](https://docs.microsoft.com/azure/azure-sql/database/purchasing-models?view=azuresql).

* `collation` - (Optional) The name of the collation. Applies only if `create_mode` is `Default`. Azure default is `SQL_LATIN1_GENERAL_CP1_CI_AS`. Changing this forces a new resource to be created.

* `max_size_bytes` - (Optional) The maximum size that the database can grow to. Applies only if `create_mode` is `Default`. Please see [Azure SQL database models](https://docs.microsoft.com/azure/azure-sql/database/purchasing-models?view=azuresql).

* `requested_service_objective_id` - (Optional) A GUID/UUID corresponding to a configured Service Level Objective for the Azure SQL database which can be used to configure a performance level.
.
* `requested_service_objective_name` - (Optional) The service objective name for the database. Valid values depend on edition and location and may include `S0`, `S1`, `S2`, `S3`, `P1`, `P2`, `P4`, `P6`, `P11` and `ElasticPool`. You can list the available names with the CLI: `shell az sql db list-editions -l westus -o table`. For further information please see [Azure CLI - az sql db](https://docs.microsoft.com/cli/azure/sql/db?view=azure-cli-latest#az-sql-db-list-editions).

* `source_database_deletion_date` - (Optional) The deletion date time of the source database. Only applies to deleted databases where `create_mode` is `PointInTimeRestore`.

* `elastic_pool_name` - (Optional) The name of the elastic database pool.

* `threat_detection_policy` - (Optional) Threat detection policy configuration. The `threat_detection_policy` block supports fields documented below.

* `read_scale` - (Optional) Read-only connections will be redirected to a high-available replica. Please see [Use read-only replicas to load-balance read-only query workloads](https://docs.microsoft.com/azure/sql-database/sql-database-read-scale-out).

* `zone_redundant` - (Optional) Whether or not this database is zone redundant, which means the replicas of this database will be spread across multiple availability zones.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `import` block supports the following:

* `storage_uri` - (Required) Specifies the blob URI of the .bacpac file.
* `storage_key` - (Required) Specifies the access key for the storage account.
* `storage_key_type` - (Required) Specifies the type of access key for the storage account. Valid values are `StorageAccessKey` or `SharedAccessKey`.
* `administrator_login` - (Required) Specifies the name of the SQL administrator.
* `administrator_login_password` - (Required) Specifies the password of the SQL administrator.
* `authentication_type` - (Required) Specifies the type of authentication used to access the server. Valid values are `SQL` or `ADPassword`.
* `operation_mode` - (Optional) Specifies the type of import operation being performed. The only allowable value is `Import`. Defaults to `Import`.

---

The `threat_detection_policy` block supports the following:

* `state` - (Optional) The State of the Policy. Possible values are `Enabled`, `Disabled` or `New`. Defaults to `Disabled`.
* `disabled_alerts` - (Optional) Specifies a list of alerts which should be disabled. Possible values include `Access_Anomaly`, `Sql_Injection` and `Sql_Injection_Vulnerability`.
* `email_account_admins` - (Optional) Should the account administrators be emailed when this alert is triggered? Possible values are `Disabled` and `Enabled`. Defaults to `Disabled`.
* `email_addresses` - (Optional) A list of email addresses which alerts should be sent to.
* `retention_days` - (Optional) Specifies the number of days to keep in the Threat Detection audit logs.
* `storage_account_access_key` - (Optional) Specifies the identifier key of the Threat Detection audit storage account. Required if `state` is `Enabled`.
* `storage_endpoint` - (Optional) Specifies the blob storage endpoint (e.g. <https://example.blob.core.windows.net>). This blob storage will hold all Threat Detection audit logs. Required if `state` is `Enabled`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The SQL Database ID.
* `creation_date` - The creation date of the SQL Database.
* `default_secondary_location` - The default secondary location of the SQL Database.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the SQL Database.
* `update` - (Defaults to 60 minutes) Used when updating the SQL Database.
* `read` - (Defaults to 5 minutes) Used when retrieving the SQL Database.
* `delete` - (Defaults to 60 minutes) Used when deleting the SQL Database.

## Import

SQL Databases can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sql_database.database1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Sql/servers/myserver/databases/database1
```
