---
subcategory: "MS SQL"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_database_blob_auditing_policies"
sidebar_current: "docs-azurerm-resource-azurerm-mssql_database_blob_auditing_policies"
description: |-
  Manages a MS SQL Azure Database Blob Auditing Policies.
---

# azurerm_mssql_database_blob_auditing_policies

Manages a MS SQL Azure Database Blob Auditing Policies.

~> **Note:** All arguments including the administrator login and password will be stored in the raw state as plain-text.
[Read more about sensitive data in state](/docs/state/sensitive-data.html).

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
 name     = "database-rg"
 location = "West US"
}
resource "azurerm_sql_server" "test" {
 name                         = "mysqlserver"
 resource_group_name          = "${azurerm_resource_group.test.name}"
 location                     = "${azurerm_resource_group.test.location}"
 version                      = "12.0"
 administrator_login          = "mradministrator"
 administrator_login_password = "thisIsDog11"
}
resource "azurerm_sql_database" "test" {
  name                             = "mysqldb"
  resource_group_name              = "${azurerm_resource_group.test.name}"
  server_name                      = "${azurerm_sql_server.test.name}"
  location                         = "${azurerm_resource_group.test.location}"
  edition                          = "Standard"
  collation                        = "SQL_Latin1_General_CP1_CI_AS"
  max_size_bytes                   = "1073741824"
  requested_service_objective_name = "S0"
}
resource "azurerm_storage_account" "test" {
 name                     = "mystorageaccount"
 resource_group_name      = "${azurerm_resource_group.test.name}"
 location                 = "${azurerm_resource_group.test.location}"
 account_tier             = "Standard"
 account_replication_type = "GRS"
}
resource "azurerm_mssql_database_blob_auditing_policies" "test"{
resource_group_name               = "${azurerm_resource_group.test.name}"
server_name                       = "${azurerm_sql_server.test.name}"
database_name                     = "${azurerm_sql_database.test.name}"
state                             = "Enabled"
storage_endpoint                  = "${azurerm_storage_account.test.primary_blob_endpoint}"
storage_account_access_key        = "${azurerm_storage_account.test.primary_access_key}"
retention_days                    = 6
is_storage_secondary_key_in_use   = true
audit_actions_and_groups          = "SUCCESSFUL_DATABASE_AUTHENTICATION_GROUP,FAILED_DATABASE_AUTHENTICATION_GROUP"
is_azure_monitor_target_enabled   = true
storage_account_subscription_id   = "00000000-0000-0000-3333-000000000000"
predicate_expression              ="object_name = 'SensitiveData'"
}
```
## Argument Reference

The following arguments are supported:

* `server_name` - (Required) The name of the SQL Server.

* `database_name` - (Required) The name of the database.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the SQL Server exists.

* `state` - (Required) Specifies the state of the policy. If state is Enabled, storageEndpoint or isAzureMonitorTargetEnabled are required. Possible values include: 'Enabled', 'Disabled'

* `storage_endpoint` - (Required) Specifies the blob storage endpoint (e.g. https://MyAccount.blob.core.windows.net). If state is Enabled, storageEndpoint is required.

* `storage_account_access_key` - (Required)Specifies the identifier key of the auditing storage account. If state is Enabled and storageEndpoint is specified, storageAccountAccessKey is required.

## Attributes Reference

* `retention_days` - Specifies the number of days to keep in the audit logs in the storage account.

* `storage_account_subscription_id` - Specifies the blob storage subscription Id.

* `is_storage_secondary_key_in_use` - Specifies whether storageAccountAccessKey value is the storage's secondary key.

* `audit_actions_and_groups` - Specifies the Actions-Groups and Actions to audit.For more information, see [Database-Level Audit Actions](https://docs.microsoft.com/en-us/sql/relational-databases/security/auditing/sql-server-audit-action-groups-and-actions#database-level-audit-actions).

* `is_azure_monitor_target_enabled` - Specifies whether audit events are sent to Azure Monitor.For more information, see [Diagnostic Settings REST API](https://go.microsoft.com/fwlink/?linkid=2033207) or [Diagnostic Settings PowerShell](https://go.microsoft.com/fwlink/?linkid=2033043).

* `predicate_expression` - Specifies condition of where clause when creating an audit.