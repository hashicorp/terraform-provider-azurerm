---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_postgresql_server"
description: |-
  Gets information about an existing PostgreSQL Azure Database Server.
---

# Data Source: azurerm_postgresql_server

Use this data source to access information about an existing PostgreSQL Azure Database Server.

## Example Usage

```hcl
data "azurerm_postgresql_server" "example" {
  name                = "postgresql-server-1"
  resource_group_name = "api-rg-pro"
}

output "postgresql_server_id" {
  value = data.azurerm_postgresql_server.example.id
}
```

## Argument Reference

* `name` - The name of the PostgreSQL Server.

* `resource_group_name` - Specifies the name of the Resource Group where the PostgreSQL Server exists.

## Attributes Reference

* `id` - The ID of the PostgreSQL Server.

* `location` - The location of the Resource Group in which the PostgreSQL Server exists.

* `fqdn` - The fully qualified domain name of the PostgreSQL Server.

* `identity` - An `identity` block as defined below.

* `version` - The version of the PostgreSQL Server.

* `administrator_login` - The administrator username of the PostgreSQL Server.

* `sku_name` - The SKU name of the PostgreSQL Server.

* `tags` - A mapping of tags assigned to the resource.

---

An `identity` block exports the following:

* `principal_id` - The ID of the System Managed Service Principal assigned to the PostgreSQL Server.

* `tenant_id` - The ID of the Tenant of the System Managed Service Principal assigned to the PostgreSQL Server.

* `type` - The identity type of the Managed Identity assigned to the PostgreSQL Server.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 10 minutes) Used when retrieving the PostgreSQL Azure Database Server.
