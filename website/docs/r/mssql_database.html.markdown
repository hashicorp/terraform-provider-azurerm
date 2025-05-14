---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_database"
description: |-
  Manages a MS SQL Database.
---

# azurerm_mssql_database

Manages a MS SQL Database.

!> **Note:** To mitigate the possibility of accidental data loss it is highly recommended that you use the `prevent_destroy` lifecycle argument in your configuration file for this resource. For more information on the `prevent_destroy` lifecycle argument please see the [terraform documentation](https://developer.hashicorp.com/terraform/tutorials/state/resource-lifecycle#prevent-resource-deletion).

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_mssql_server" "example" {
  name                         = "example-sqlserver"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_mssql_database" "example" {
  name         = "example-db"
  server_id    = azurerm_mssql_server.example.id
  collation    = "SQL_Latin1_General_CP1_CI_AS"
  license_type = "LicenseIncluded"
  max_size_gb  = 2
  sku_name     = "S0"
  enclave_type = "VBS"

  tags = {
    foo = "bar"
  }

  # prevent the possibility of accidental data loss
  lifecycle {
    prevent_destroy = true
  }
}
```

## Example Usage for Transparent Data Encryption(TDE) with a Customer Managed Key(CMK) during Create
```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_user_assigned_identity" "example" {
  name                = "example-admin"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_storage_account" "example" {
  name                     = "examplesa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_server" "example" {
  name                         = "example-sqlserver"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"
}

resource "azurerm_mssql_database" "example" {
  name           = "example-db"
  server_id      = azurerm_mssql_server.example.id
  collation      = "SQL_Latin1_General_CP1_CI_AS"
  license_type   = "LicenseIncluded"
  max_size_gb    = 4
  read_scale     = true
  sku_name       = "S0"
  zone_redundant = true
  enclave_type   = "VBS"

  tags = {
    foo = "bar"
  }

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.example.id]
  }

  transparent_data_encryption_key_vault_key_id = azurerm_key_vault_key.example.id

  # prevent the possibility of accidental data loss
  lifecycle {
    prevent_destroy = true
  }
}

# Create a key vault with access policies which allow for the current user to get, list, create, delete, update, recover, purge and getRotationPolicy for the key vault key and also add a key vault access policy for the Microsoft Sql Server instance User Managed Identity to get, wrap, and unwrap key(s)
resource "azurerm_key_vault" "example" {
  name                        = "mssqltdeexample"
  location                    = azurerm_resource_group.example.location
  resource_group_name         = azurerm_resource_group.example.name
  enabled_for_disk_encryption = true
  tenant_id                   = azurerm_user_assigned_identity.example.tenant_id
  soft_delete_retention_days  = 7
  purge_protection_enabled    = true

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = ["Get", "List", "Create", "Delete", "Update", "Recover", "Purge", "GetRotationPolicy"]
  }

  access_policy {
    tenant_id = azurerm_user_assigned_identity.example.tenant_id
    object_id = azurerm_user_assigned_identity.example.principal_id

    key_permissions = ["Get", "WrapKey", "UnwrapKey"]
  }
}

resource "azurerm_key_vault_key" "example" {
  depends_on = [azurerm_key_vault.example]

  name         = "example-key"
  key_vault_id = azurerm_key_vault.example.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = ["unwrapKey", "wrapKey"]
}

```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the MS SQL Database. Changing this forces a new resource to be created. 

* `server_id` - (Required) The id of the MS SQL Server on which to create the database. Changing this forces a new resource to be created.

~> **Note:** This setting is still required for "Serverless" SKUs

* `auto_pause_delay_in_minutes` - (Optional) Time in minutes after which database is automatically paused. A value of `-1` means that automatic pause is disabled. This property is only settable for Serverless databases.

* `create_mode` - (Optional) The create mode of the database. Possible values are `Copy`, `Default`, `OnlineSecondary`, `PointInTimeRestore`, `Recovery`, `Restore`, `RestoreExternalBackup`, `RestoreExternalBackupSecondary`, `RestoreLongTermRetentionBackup` and `Secondary`. Mutually exclusive with `import`. Changing this forces a new resource to be created. Defaults to `Default`.

* `import` - (Optional) A `import` block as documented below. Mutually exclusive with `create_mode`.

* `creation_source_database_id` - (Optional) The ID of the source database from which to create the new database. This should only be used for databases with `create_mode` values that use another database as reference. Changing this forces a new resource to be created.

-> **Note:** When configuring a secondary database, please be aware of the constraints for the `sku_name` property, as noted below, for both the primary and secondary databases. The `sku_name` of the secondary database may be inadvertently changed to match that of the primary when an incompatible combination of SKUs is detected by the provider.

* `collation` - (Optional) Specifies the collation of the database. Changing this forces a new resource to be created.

* `elastic_pool_id` - (Optional) Specifies the ID of the elastic pool containing this database.

* `enclave_type` - (Optional) Specifies the type of enclave to be used by the elastic pool. When `enclave_type` is not specified (e.g., the default) enclaves are not enabled on the database. Once enabled (e.g., by specifying `Default` or `VBS`) removing the `enclave_type` field from the configuration file will force the creation of a new resource. Possible values are `Default` or `VBS`.

-> **Note:** `enclave_type` is currently not supported for DW (e.g, DataWarehouse) and DC-series SKUs.

-> **Note:** Geo Replicated and Failover databases must have the same `enclave_type`.

~> **Note:** The default value for the `enclave_type` field is unset not `Default`.

* `geo_backup_enabled` - (Optional) A boolean that specifies if the Geo Backup Policy is enabled. Defaults to `true`.

~> **Note:** `geo_backup_enabled` is only applicable for DataWarehouse SKUs (DW*). This setting is ignored for all other SKUs.

* `maintenance_configuration_name` - (Optional) The name of the Public Maintenance Configuration window to apply to the database. Valid values include `SQL_Default`, `SQL_EastUS_DB_1`, `SQL_EastUS2_DB_1`, `SQL_SoutheastAsia_DB_1`, `SQL_AustraliaEast_DB_1`, `SQL_NorthEurope_DB_1`, `SQL_SouthCentralUS_DB_1`, `SQL_WestUS2_DB_1`, `SQL_UKSouth_DB_1`, `SQL_WestEurope_DB_1`, `SQL_EastUS_DB_2`, `SQL_EastUS2_DB_2`, `SQL_WestUS2_DB_2`, `SQL_SoutheastAsia_DB_2`, `SQL_AustraliaEast_DB_2`, `SQL_NorthEurope_DB_2`, `SQL_SouthCentralUS_DB_2`, `SQL_UKSouth_DB_2`, `SQL_WestEurope_DB_2`, `SQL_AustraliaSoutheast_DB_1`, `SQL_BrazilSouth_DB_1`, `SQL_CanadaCentral_DB_1`, `SQL_CanadaEast_DB_1`, `SQL_CentralUS_DB_1`, `SQL_EastAsia_DB_1`, `SQL_FranceCentral_DB_1`, `SQL_GermanyWestCentral_DB_1`, `SQL_CentralIndia_DB_1`, `SQL_SouthIndia_DB_1`, `SQL_JapanEast_DB_1`, `SQL_JapanWest_DB_1`, `SQL_NorthCentralUS_DB_1`, `SQL_UKWest_DB_1`, `SQL_WestUS_DB_1`, `SQL_AustraliaSoutheast_DB_2`, `SQL_BrazilSouth_DB_2`, `SQL_CanadaCentral_DB_2`, `SQL_CanadaEast_DB_2`, `SQL_CentralUS_DB_2`, `SQL_EastAsia_DB_2`, `SQL_FranceCentral_DB_2`, `SQL_GermanyWestCentral_DB_2`, `SQL_CentralIndia_DB_2`, `SQL_SouthIndia_DB_2`, `SQL_JapanEast_DB_2`, `SQL_JapanWest_DB_2`, `SQL_NorthCentralUS_DB_2`, `SQL_UKWest_DB_2`, `SQL_WestUS_DB_2`, `SQL_WestCentralUS_DB_1`, `SQL_FranceSouth_DB_1`, `SQL_WestCentralUS_DB_2`, `SQL_FranceSouth_DB_2`, `SQL_SwitzerlandNorth_DB_1`, `SQL_SwitzerlandNorth_DB_2`, `SQL_BrazilSoutheast_DB_1`, `SQL_UAENorth_DB_1`, `SQL_BrazilSoutheast_DB_2`, `SQL_UAENorth_DB_2`, `SQL_SouthAfricaNorth_DB_1`, `SQL_SouthAfricaNorth_DB_2`, `SQL_WestUS3_DB_1`, `SQL_WestUS3_DB_2`, `SQL_SwedenCentral_DB_1`, `SQL_SwedenCentral_DB_2`. Defaults to `SQL_Default`.

~> **Note:** `maintenance_configuration_name` is only applicable if `elastic_pool_id` is not set.

* `ledger_enabled` - (Optional) A boolean that specifies if this is a ledger database. Defaults to `false`. Changing this forces a new resource to be created.

* `license_type` - (Optional) Specifies the license type applied to this database. Possible values are `LicenseIncluded` and `BasePrice`.

* `long_term_retention_policy` - (Optional) A `long_term_retention_policy` block as defined below.

* `max_size_gb` - (Optional) The max size of the database in gigabytes.

~> **Note:** This value should not be configured when the `create_mode` is `Secondary` or `OnlineSecondary`, as the sizing of the primary is then used as per [Azure documentation](https://docs.microsoft.com/azure/azure-sql/database/single-database-scale#geo-replicated-database).

* `min_capacity` - (Optional) Minimal capacity that database will always have allocated, if not paused. This property is only settable for Serverless databases.

* `restore_point_in_time` - (Optional) Specifies the point in time (ISO8601 format) of the source database that will be restored to create the new database. This property is only settable for `create_mode`= `PointInTimeRestore` databases.

* `recover_database_id` - (Optional) The ID of the database to be recovered. This property is only applicable when the `create_mode` is `Recovery`.

* `recovery_point_id` - (Optional) The ID of the Recovery Services Recovery Point Id to be restored. This property is only applicable when the `create_mode` is `Recovery`.

* `restore_dropped_database_id` - (Optional) The ID of the database to be restored. This property is only applicable when the `create_mode` is `Restore`.

* `restore_long_term_retention_backup_id` - (Optional) The ID of the long term retention backup to be restored. This property is only applicable when the `create_mode` is `RestoreLongTermRetentionBackup`.

* `read_replica_count` - (Optional) The number of readonly secondary replicas associated with the database to which readonly application intent connections may be routed. This property is only settable for Hyperscale edition databases.

* `read_scale` - (Optional) If enabled, connections that have application intent set to readonly in their connection string may be routed to a readonly secondary replica. This property is only settable for Premium and Business Critical databases.

* `sample_name` - (Optional) Specifies the name of the sample schema to apply when creating this database. Possible value is `AdventureWorksLT`.

* `short_term_retention_policy` - (Optional) A `short_term_retention_policy` block as defined below.

* `sku_name` - (Optional) Specifies the name of the SKU used by the database. For example, `GP_S_Gen5_2`,`HS_Gen4_1`,`BC_Gen5_2`, `ElasticPool`, `Basic`,`S0`, `P2` ,`DW100c`, `DS100`. Changing this from the HyperScale service tier to another service tier will create a new resource.

-> **Note:** A full list of supported SKU names by region can be retrieved using the Azure CLI: `az sql db list-editions -l <region> -o table`

-> **Note:** The default `sku_name` value may differ between Azure locations depending on local availability of Gen4/Gen5 capacity. When databases are replicated using the `creation_source_database_id` property, the source (primary) database cannot have a higher SKU service tier than any secondary databases. When changing the `sku_name` of a database having one or more secondary databases, this resource will first update any secondary databases as necessary. In such cases it's recommended to use the same `sku_name` in your configuration for all related databases, as not doing so may cause an unresolvable diff during subsequent plans.

* `storage_account_type` - (Optional) Specifies the storage account type used to store backups for this database. Possible values are `Geo`, `GeoZone`, `Local` and `Zone`. Defaults to `Geo`.

* `threat_detection_policy` - (Optional) Threat detection policy configuration. The `threat_detection_policy` block supports fields documented below.

* `identity` - (Optional) An `identity` block as defined below.

* `transparent_data_encryption_enabled` - (Optional) If set to true, Transparent Data Encryption will be enabled on the database. Defaults to `true`.

-> **Note:** `transparent_data_encryption_enabled` can only be set to `false` on DW (e.g, DataWarehouse) server SKUs.

* `transparent_data_encryption_key_vault_key_id` - (Optional) The fully versioned `Key Vault` `Key` URL (e.g. `'https://<YourVaultName>.vault.azure.net/keys/<YourKeyName>/<YourKeyVersion>`) to be used as the `Customer Managed Key`(CMK/BYOK) for the `Transparent Data Encryption`(TDE) layer.

~> **Note:** To successfully deploy a `Microsoft SQL Database` in CMK/BYOK TDE the `Key Vault` must have `Soft-delete` and `purge protection` enabled to protect from data loss due to accidental key and/or key vault deletion. The `Key Vault` and the `Microsoft SQL Server` `User Managed Identity Instance` must belong to the same `Azure Active Directory` `tenant`.

* `transparent_data_encryption_key_automatic_rotation_enabled` - (Optional) Boolean flag to specify whether TDE automatically rotates the encryption Key to latest version or not. Possible values are `true` or `false`. Defaults to `false`.

~> **Note:** When the `sku_name` is `DW100c`, the `transparent_data_encryption_key_automatic_rotation_enabled` and the `transparent_data_encryption_key_vault_key_id` properties should not be specified, as database-level CMK is not supported for Data Warehouse SKUs.

* `zone_redundant` - (Optional) Whether or not this database is zone redundant, which means the replicas of this database will be spread across multiple availability zones. This property is only settable for Premium and Business Critical databases.

* `secondary_type` - (Optional) How do you want your replica to be made? Valid values include `Geo` and `Named`. Defaults to `Geo`. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---
a `import` block supports the following:

* `storage_uri` - (Required) Specifies the blob URI of the .bacpac file.
* `storage_key` - (Required) Specifies the access key for the storage account.
* `storage_key_type` - (Required) Specifies the type of access key for the storage account. Valid values are `StorageAccessKey` or `SharedAccessKey`.
* `administrator_login` - (Required) Specifies the name of the SQL administrator.
* `administrator_login_password` - (Required) Specifies the password of the SQL administrator.
* `authentication_type` - (Required) Specifies the type of authentication used to access the server. Valid values are `SQL` or `ADPassword`.
* `storage_account_id` - (Optional) The resource id for the storage account used to store BACPAC file. If set, private endpoint connection will be created for the storage account. Must match storage account used for storage_uri parameter.

---
a `threat_detection_policy` block supports the following:

* `state` - (Optional) The State of the Policy. Possible values are `Enabled` or `Disabled`. Defaults to `Disabled`.
* `disabled_alerts` - (Optional) Specifies a list of alerts which should be disabled. Possible values include `Access_Anomaly`, `Sql_Injection` and `Sql_Injection_Vulnerability`.
* `email_account_admins` - (Optional) Should the account administrators be emailed when this alert is triggered? Possible values are `Enabled` or `Disabled`. Defaults to `Disabled`.
* `email_addresses` - (Optional) A list of email addresses which alerts should be sent to.
* `retention_days` - (Optional) Specifies the number of days to keep in the Threat Detection audit logs.
* `storage_account_access_key` - (Optional) Specifies the identifier key of the Threat Detection audit storage account. Required if `state` is `Enabled`.
* `storage_endpoint` - (Optional) Specifies the blob storage endpoint (e.g. <https://example.blob.core.windows.net>). This blob storage will hold all Threat Detection audit logs. Required if `state` is `Enabled`.

---

A `long_term_retention_policy` block supports the following:

* `weekly_retention` - (Optional) The weekly retention policy for an LTR backup in an ISO 8601 format. Valid value is between 1 to 520 weeks. e.g. `P1Y`, `P1M`, `P1W` or `P7D`. Defaults to `PT0S`.
* `monthly_retention` - (Optional) The monthly retention policy for an LTR backup in an ISO 8601 format. Valid value is between 1 to 120 months. e.g. `P1Y`, `P1M`, `P4W` or `P30D`. Defaults to `PT0S`.
* `yearly_retention` - (Optional) The yearly retention policy for an LTR backup in an ISO 8601 format. Valid value is between 1 to 10 years. e.g. `P1Y`, `P12M`, `P52W` or `P365D`. Defaults to `PT0S`.
* `week_of_year` - (Optional) The week of year to take the yearly backup. Value has to be between `1` and `52`.

---

A `short_term_retention_policy` block supports the following:

* `retention_days` - (Required) Point In Time Restore configuration. Value has to be between `1` and `35`.
* `backup_interval_in_hours` - (Optional) The hours between each differential backup. This is only applicable to live databases but not dropped databases. Value has to be `12` or `24`. Defaults to `12` hours.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this SQL Database. Possible value is `UserAssigned`.

* `identity_ids` - (Required) Specifies a list of User Assigned Managed Identity IDs to be assigned to this SQL Database.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the MS SQL Database.

---

A `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Identity of this SQL Database.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Identity of this SQL Database.

-> **Note:** You can access the Principal ID via `azurerm_mssql_database.example.identity[0].principal_id` and the Tenant ID via `azurerm_mssql_database.example.identity[0].tenant_id`

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the MS SQL Database.
* `read` - (Defaults to 5 minutes) Used when retrieving the MS SQL Database.
* `update` - (Defaults to 1 hour) Used when updating the MS SQL Database.
* `delete` - (Defaults to 1 hour) Used when deleting the MS SQL Database.

## Import

SQL Database can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_database.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/databases/example1
```
