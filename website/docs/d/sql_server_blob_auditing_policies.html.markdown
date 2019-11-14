---
subcategory:""
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sql_server_blob_auditing_policies"
sidebar_current: "docs-azurerm-datasource-sql-server_blob_auditing_policies"
description: |-
  Gets information about an existing SQL Azure Database Server Blob Auditing Policies.
---

# Data Source: azurerm_sql_server_blob_auditing_policies

Use this data source to access information about an existing SQL Azure Database Server Blob Auditing Policies.

## Example Usage

```hcl
data "azurerm_sql_server_blob_auditing_policies" "test"{
server_name                              = "example-server-name"
resource_group_name                      = "example-resources"
state                                    = "Enabled"
storage_endpoint                         = "example-storage-account-primary_blob_endpoint"
storage_account_access_key               = "example-storage-account-primary_access_key"
retention_days                           = 0
is_storage_secondary_key_in_use          = true
audit_actions_and_groups                 = "example-audit_actions_and_groups"
is_azure_monitor_target_enabled          = true
storage_account_subscription_id          ="00000000-0000-0000-0000-000000000000"
}
output "sql_server_id" {
  value = "${data.azurerm_sql_server_blob_auditing_policies.test.id}"
}
```
## Argument Reference

* `server_name` - (Required) The name of the SQL Server.

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


