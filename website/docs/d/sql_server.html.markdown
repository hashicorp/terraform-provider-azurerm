---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sql_server"
description: |-
  Gets information about an existing SQL Azure Database Server.
---

# Data Source: azurerm_sql_server

Use this data source to access information about an existing SQL Azure Database Server.

-> **Note:** The `azurerm_sql_server` data source is deprecated in version 3.0 of the AzureRM provider and will be removed in version 4.0. Please use the [`azurerm_mssql_server`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/data-sources/mssql_server) data source instead.

## Example Usage

```hcl
data "azurerm_sql_server" "example" {
  name                = "examplesqlservername"
  resource_group_name = "example-resources"
}

output "sql_server_id" {
  value = data.azurerm_sql_server.example.id
}
```

## Argument Reference

* `name` - The name of the SQL Server.

* `resource_group_name` - Specifies the name of the Resource Group where the SQL Server exists.

## Attributes Reference

* `id` - The id of the SQL Server resource.

* `location` - The location of the Resource Group in which the SQL Server exists.

* `fqdn` - The fully qualified domain name of the SQL Server.

* `version` - The version of the SQL Server.

* `administrator_login` - The administrator username of the SQL Server.

* `identity` - An `identity` block as defined below.

* `tags` - A mapping of tags assigned to the resource.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

* `type` - The identity type of this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the SQL Azure Database Server.
