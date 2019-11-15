---
subcategory: "MS SQL"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_server_blob_extended_auditing_policies"
sidebar_current: "docs-azurerm-datasource-mssql-server-blob-extended-auditing-policies"
description: |-
  Gets information about an existing MS SQL Azure Server Blob Extended Auditing Policies.
---

# Data Source: azurerm_mssql_server_blob_extended_auditing_policies

Use this data source to access information about an existing MS SQL Azure Extended Server Blob Auditing Policies.

## Example Usage

```hcl
data "azurerm_mssql_server_blob_extended_auditing_policies" "test"{
server_name                              = "example-server-name"
resource_group_name                      = "example-resources"
}
output "mssql_server_id" {
  value = "${data.azurerm_mssql_server_blob_extended_auditing_policies.test.id}"
}
```
## Argument Reference

* `server_name` - (Required) The name of the SQL Server.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the SQL Server exists.

