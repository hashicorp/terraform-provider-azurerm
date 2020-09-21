---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_mssql_server"
description: |-
  Gets information about an existing Microsoft SQL Server.
---

# Data Source: azurerm_mssql_server

Use this data source to access information about an existing Microsoft SQL Server.

## Example Usage

```hcl
data "azurerm_mssql_server" "example" {
  name                = "existingMsSqlServer"
  resource_group_name = "existingResGroup"
}

output "id" {
  value = data.azurerm_mssql_server.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Microsoft SQL Server.

* `resource_group_name` - (Required) The name of the Resource Group where the Microsoft SQL Server exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Microsoft SQL Server.

* `administrator_login` - The server's administrator login name.

* `fully_qualified_domain_name` - The fully qualified domain name of the Azure SQL Server.

* `identity` - A `identity` block as defined below.

* `location` - The Azure Region where the Microsoft SQL Server exists.

* `restorable_dropped_database_ids` - A list of dropped restorable database IDs on the server.

* `tags` - A mapping of tags assigned to this Microsoft SQL Server.

* `version` - This servers MS SQL version.

---

A `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Identity of this SQL Server.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Identity of this SQL Server.

* `type` - The identity type of the Microsoft SQL Server.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Microsoft SQL Server.
