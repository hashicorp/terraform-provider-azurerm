---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sql_server"
sidebar_current: "docs-azurerm-datasource-sql-server"
description: |-
  Gets information about an existing SQL Azure Database Server.
---

# Data Source: azurerm_sql_server

Use this data source to access information about an existing SQL Azure Database Server.

## Example Usage

```hcl
data "azurerm_sql_server" "example" {
  name                = "examplesqlservername"
  resource_group_name = "example-resources"
}

output "sql_server_id" {
  value = "${data.azurerm_sql_server.example.id}"
}
```

## Argument Reference

* `name` - (Required) The name of the SQL Server.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the SQL Server exists.

## Attributes Reference

* `location` - The location of the Resource Group in which the SQL Server exists.

* `fqdn` - The fully qualified domain name of the SQL Server.

* `version` - The version of the SQL Server.

* `administrator_login` - The administrator username of the SQL Server.

* `identity` - An `identity` block as defined below.

* `blob_extended_auditing_policy` - An `blob_extended_auditing_policy` block as defined below.

* `tags` - A mapping of tags assigned to the resource.

---

An `identity` block exports the following:

* `principal_id` - The ID of the Principal (Client) in Azure Active Directory.

* `tenant_id` - The ID of the Azure Active Directory Tenant.

* `type` - The identity type of the SQL Server.

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
