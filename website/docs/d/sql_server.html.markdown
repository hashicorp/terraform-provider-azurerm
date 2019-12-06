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

* `tags` - A mapping of tags assigned to the resource.

---

An `identity` block exports the following:

* `principal_id` - The ID of the Principal (Client) in Azure Active Directory.

* `tenant_id` - The ID of the Azure Active Directory Tenant.

* `type` - The identity type of the SQL Server.
