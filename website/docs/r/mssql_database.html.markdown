---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_database"
description: |-
  Manages a MS SQL Database.
---

# azurerm_mssql_database

Manages a MS SQL Database.

~> **NOTE:** The Database Extended Auditing Policy Can be set inline here as well as with the [mssql_database_extended_auditing_policy resource](mssql_database_extended_auditing_policy.html) resource. You can only use one or the other and using both will cause a conflict.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplesa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
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
  name           = "acctest-db-d"
  server_id      = azurerm_sql_server.example.id
  collation      = "SQL_Latin1_General_CP1_CI_AS"
  license_type   = "LicenseIncluded"
  max_size_gb    = 4
  read_scale     = true
  sku_name       = "BC_Gen5_2"
  zone_redundant = true

  extended_auditing_policy {
    storage_endpoint                        = azurerm_storage_account.example.primary_blob_endpoint
    storage_account_access_key              = azurerm_storage_account.example.primary_access_key
    storage_account_access_key_is_secondary = true
    retention_in_days                       = 6
  }


  tags = {
    foo = "bar"
  }

}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Ms SQL Database. Changing this forces a new resource to be created.

* `server_id` - (Required) The id of the Ms SQL Server on which to create the database. Changing this forces a new resource to be created.

~> **NOTE:** This setting is still required for "Serverless" SKU's

* `auto_pause_delay_in_minutes` - (Optional) Time in minutes after which database is automatically paused. A value of `-1` means that automatic pause is disabled. This property is only settable for General Purpose Serverless databases.

* `create_mode` - (Optional) The create mode of the database. Possible values are `Copy`, `Default`, `OnlineSecondary`, `PointInTimeRestore`, `Recovery`, `Restore`, `RestoreExternalBackup`, `RestoreExternalBackupSecondary`, `RestoreLongTermRetentionBackup` and `Secondary`. 

* `creation_source_database_id` - (Optional) The id of the source database to be referred to create the new database. This should only be used for databases with `create_mode` values that use another database as reference. Changing this forces a new resource to be created.

* `collation` - (Optional) Specifies the collation of the database. Changing this forces a new resource to be created.

* `elastic_pool_id` - (Optional) Specifies the ID of the elastic pool containing this database.

* `extended_auditing_policy` - (Optional) A `extended_auditing_policy` block as defined below.

* `license_type` - (Optional) Specifies the license type applied to this database. Possible values are `LicenseIncluded` and `BasePrice`.

* `long_term_retention_policy` - (Optional) A `long_term_retention_policy` block as defined below.

* `max_size_gb` - (Optional) The max size of the database in gigabytes. 

* `min_capacity` - (Optional) Minimal capacity that database will always have allocated, if not paused. This property is only settable for General Purpose Serverless databases.

* `restore_point_in_time` - (Required) Specifies the point in time (ISO8601 format) of the source database that will be restored to create the new database. This property is only settable for `create_mode`= `PointInTimeRestore`  databases.

* `recover_database_id` - (Optional) The ID of the database to be recovered. This property is only applicable when the `create_mode` is `Recovery`.

* `restore_dropped_database_id` - (Optional) The ID of the database to be restored. This property is only applicable when the `create_mode` is `Restore`.

* `read_replica_count` - (Optional) The number of readonly secondary replicas associated with the database to which readonly application intent connections may be routed. This property is only settable for Hyperscale edition databases.

* `read_scale` - (Optional) If enabled, connections that have application intent set to readonly in their connection string may be routed to a readonly secondary replica. This property is only settable for Premium and Business Critical databases.

* `sample_name` - (Optional) Specifies the name of the sample schema to apply when creating this database. Possible value is `AdventureWorksLT`.

* `short_term_retention_policy` - (Optional) A `short_term_retention_policy` block as defined below.

* `sku_name` - (Optional) Specifies the name of the sku used by the database. Changing this forces a new resource to be created. For example, `GP_S_Gen5_2`,`HS_Gen4_1`,`BC_Gen5_2`, `ElasticPool`, `Basic`,`S0`, `P2` ,`DW100c`, `DS100`.

~> **NOTE** The default sku_name value may differ between Azure locations depending on local availability of Gen4/Gen5 capacity.

* `storage_account_type` - (Optional) Specifies the storage account type used to store backups for this database. Changing this forces a new resource to be created.  Possible values are `GRS`, `LRS` and `ZRS`.  The default value is `GRS`.

* `threat_detection_policy` - (Optional) Threat detection policy configuration. The `threat_detection_policy` block supports fields documented below.

* `zone_redundant` - (Optional) Whether or not this database is zone redundant, which means the replicas of this database will be spread across multiple availability zones. This property is only settable for Premium and Business Critical databases.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---
a `threat_detection_policy` block supports the following:

* `state` - (Required) The State of the Policy. Possible values are `Enabled`, `Disabled` or `New`.
* `disabled_alerts` - (Optional) Specifies a list of alerts which should be disabled. Possible values include `Access_Anomaly`, `Sql_Injection` and `Sql_Injection_Vulnerability`.
* `email_account_admins` - (Optional) Should the account administrators be emailed when this alert is triggered?
* `email_addresses` - (Optional) A list of email addresses which alerts should be sent to.
* `retention_days` - (Optional) Specifies the number of days to keep in the Threat Detection audit logs.
* `storage_account_access_key` - (Optional) Specifies the identifier key of the Threat Detection audit storage account. Required if `state` is `Enabled`.
* `storage_endpoint` - (Optional) Specifies the blob storage endpoint (e.g. https://MyAccount.blob.core.windows.net). This blob storage will hold all Threat Detection audit logs. Required if `state` is `Enabled`.
* `use_server_default` - (Optional) Should the default server policy be used? Defaults to `Disabled`.

---

A `extended_auditing_policy` block supports the following:

* `storage_account_access_key` - (Optional)  Specifies the access key to use for the auditing storage account.
* `storage_endpoint` - (Optional) Specifies the blob storage endpoint (e.g. https://MyAccount.blob.core.windows.net).
* `storage_account_access_key_is_secondary` - (Optional) Specifies whether `storage_account_access_key` value is the storage's secondary key.
* `retention_in_days` - (Optional) Specifies the number of days to retain logs for in the storage account.
* `monitor_enabled` - (Optional) Enable audit events to Azure Monitor? To enable audit events to Log Analytics, please refer to the example which can be found in [the `./examples/sql-azure/sql_auditing_log_analytics` directory within the Github Repository](https://github.com/terraform-providers/terraform-provider-azurerm/tree/master/examples/sql-azure/sql_auditing_log_analytics). To enable audit events to Eventhub, please refer to the example which can be found in [the `./examples/sql-azure/sql_auditing_eventhub` directory within the Github Repository](https://github.com/terraform-providers/terraform-provider-azurerm/tree/master/examples/sql-azure/sql_auditing_eventhub). 

---

A `long_term_retention_policy` block supports the following:

* `weekly_retention` - (Optional) The weekly retention policy for an LTR backup in an ISO 8601 format. Valid value is between 1 to 520 weeks. e.g. `P1Y`, `P1M`, `P1W` or `P7D`.
* `monthly_retention` - (Optional) The monthly retention policy for an LTR backup in an ISO 8601 format. Valid value is between 1 to 120 months. e.g. `P1Y`, `P1M`, `P4W` or `P30D`.
* `yearly_retention` - (Optional) The yearly retention policy for an LTR backup in an ISO 8601 format. Valid value is between 1 to 10 years. e.g. `P1Y`, `P12M`, `P52W` or `P365D`.
* `week_of_year` - (Optional) The week of year to take the yearly backup in an ISO 8601 format. Value has to be between `1` and `52`.

---

A `short_term_retention_policy` block supports the following:

* `retention_days` - (Required) Point In Time Restore configuration. Value has to be between `7` and `35`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the MS SQL Database.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the MS SQL Database.
* `update` - (Defaults to 60 minutes) Used when updating the MS SQL Database.
* `read` - (Defaults to 5 minutes) Used when retrieving the MS SQL Database.
* `delete` - (Defaults to 60 minutes) Used when deleting the MS SQL Database.

## Import

SQL Database can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_database.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/databases/example1
```
