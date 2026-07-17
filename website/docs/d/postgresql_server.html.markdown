---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_postgresql_server"
description: |-
  Gets information about an existing PostgreSQL Azure Database Server.
---

# Data Source: azurerm_postgresql_server

Use this data source to access information about an existing PostgreSQL Azure Database Server.

~> **Note:** The `azurerm_postgresql_server` data source is deprecated and will be removed in v5.0 of the AzureRM Provider. Azure Database for PostgreSQL Single Server and its sub resources have been retired as of 2025-03-28, please use the `azurerm_postgresql_flexible_server` data source instead. For more information, see https://techcommunity.microsoft.com/blog/adforpostgresql/retiring-azure-database-for-postgresql-single-server-in-2025/3783783.

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

## Arguments Reference

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

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 10 minutes) Used when retrieving the PostgreSQL Azure Database Server.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.DBforPostgreSQL` - 2017-12-01
