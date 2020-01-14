---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sql_database"
description: |-
  Gets information about an existing SQL Azure Database.
---

# Data Source: azurerm_sql_database

Use this data source to access information about an existing SQL Azure Database.

## Example Usage

```hcl
data "azurerm_sql_database" "example" {
  name                = "example_db"
  server_name         = "example_db_server"
  resource_group_name = "example-resources"
}

output "sql_database_id" {
  value = data.azurerm_sql_database.example.id
}
```

## Argument Reference

* `name` - (Required) The name of the SQL Database.

* `server_name` - (Required) The name of the SQL Server.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Azure SQL Database exists.

## Attributes Reference

* `id` - The SQL Database ID.

* `collation` - The name of the collation. 
    
* `creation_date` - The creation date of the SQL Database.

* `default_secondary_location` - The default secondary location of the SQL Database.

* `edition` - The edition of the database.

* `elastic_pool_name` - The name of the elastic database pool the database belongs to.

* `failover_group_id` - The ID of the failover group the database belongs to.

* `location` - The location of the Resource Group in which the SQL Server exists.

* `name` - The name of the database.
    
* `read_scale` - Indicate if read-only connections will be redirected to a high-available replica.

* `requested_service_objective_id` - The ID pertaining to the performance level of the database.
 
* `requested_service_objective_name` - The name pertaining to the performance level of the database. 

* `resource_group_name` - The name of the resource group in which the database resides. This will always be the same resource group as the Database Server.

* `server_name` - The name of the SQL Server on which to create the database.
    
* `tags` - A mapping of tags assigned to the resource.

* `blob_extended_auditing_policy` - An `blob_extended_auditing_policy` block as defined below.

---

An `blob_extended_auditing_policy` block exports the following:

* `state` - Specifies the state of the policy. If state is Enabled, storageEndpoint or isAzureMonitorTargetEnabled are required. Possible values include: 'Enabled', 'Disabled'

* `storage_endpoint` - Specifies the blob storage endpoint (e.g. https://MyAccount.blob.core.windows.net). If state is Enabled, storageEndpoint is required.

* `storage_account_access_key` - Specifies the identifier key of the auditing storage account. If state is Enabled and storageEndpoint is specified, storageAccountAccessKey is required.

* `retention_days` - Specifies the number of days to keep in the audit logs in the storage account.

* `storage_account_subscription_id` - Specifies the blob storage subscription Id.

* `is_storage_secondary_key_in_use` - Specifies whether storageAccountAccessKey value is the storage's secondary key.

* `audit_actions_and_groups` - Specifies the Actions-Groups and Actions to audit.For more information, see [Database-Level Audit Actions](https://docs.microsoft.com/en-us/sql/relational-databases/security/auditing/sql-server-audit-action-groups-and-actions#database-level-audit-actions).

* `predicate_expression` - Specifies condition of where clause when creating an audit.
