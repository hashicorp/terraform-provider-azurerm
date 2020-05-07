---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sql_server"
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

* `principal_id` - The ID of the Principal (Client) in Azure Active Directory.

* `tenant_id` - The ID of the Azure Active Directory Tenant.

* `type` - The identity type of the SQL Server.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the SQL Azure Database Server.
