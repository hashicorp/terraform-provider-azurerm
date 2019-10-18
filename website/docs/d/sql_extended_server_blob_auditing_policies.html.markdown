---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sql_extended_server_blob_auditing_policies"
sidebar_current: "docs-azurerm-datasource-sql-extended-server-blob-auditing-policies"
description: |-
  Gets information about an existing SQL Azure Database Server.
---

# Data Source: azurerm_sql_server_blob_auditing_policies

Use this data source to access information about an existing SQL Azure Database Server.

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
predicate_expression                     ="object_name = 'SensitiveData'"
}
output "sql_server_id" {
  value = "${data.azurerm_sql_server_blob_auditing_policies.test.id}"
}
```

## Argument Reference

* `server_name` - (Required) The name of the SQL Server.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the SQL Server exists.

* `state` - (Required) Specifies the state "Enabled"/"Disabled" of Blob Auditing Policies.

* `storage_endpoint` - (Required) Specifies the endpoint of the storage account to be connected.

* `storage_account_access_key` - (Required) Specifies the access key of the storage account to be connected.

## Attributes Reference

* `retention_days` - Specifies the rentention days of the Blob Auditing Policies.

* `is_storage_secondary_key_in_use` - Specifies if the storage secondary key is in use.

* `audit_actions_and_groups` - List of the audit actions and groups ("SUCCESSFUL_DATABASE_AUTHENTICATION_GROUP","FAILED_DATABASE_AUTHENTICATION_GROUP","BATCH_COMPLETED_GROUP").

* `is_azure_monitor_target_enabled` - Specifies if the azure monitor target is enabled.

* `storage_account_subscription_id` - Specifies the subscription id of the storage account to be connected.

* `predicate_expression` - Specifies the predicate expression of the Blob Auditing Policies.
