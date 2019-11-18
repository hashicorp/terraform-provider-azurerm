---
subcategory: "MS SQL"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_database_blob_auditing_policies"
sidebar_current: "docs-azurerm-datasource-mssql-database-blob-auditing-policies"
description: |-
  Gets information about an existing MS SQL Azure Database Blob Auditing Policies.
---

# Data Source: azurerm_mssql_database_blob_auditing_policies

Use this data source to access information about an existing MS SQL Azure Database Blob Auditing Policies.

## Example Usage

```hcl
data "azurerm_mssql_database_blob_auditing_policies" "test"{
server_name                              = "example-server-name"
database_name                            = "example-database-name"
resource_group_name                      = "example-resources"
}
output "mssql_database_blob_auditing_policies_id" {
  value = "${data.azurerm_mssql_database_blob_auditing_policies.test.id}"
}
```
## Argument Reference

* `server_name` - (Required) The name of the SQL Server.

* `database_name` - (Required) The name of the database.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the SQL Server exists.
