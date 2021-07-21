---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_postgresql_flexible_server_configuration"
description: |-
  Sets a PostgreSQL Configuration value on a Azure PostgreSQL Flexible Server.
---

# azurerm_postgresql_flexible_server_configuration

Sets a PostgreSQL Configuration value on a Azure PostgreSQL Flexible Server.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_postgresql_flexible_server" "example" {
  name                   = "example-psqlflexibleserver"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  version                = "12"
  administrator_login    = "psqladminun"
  administrator_password = "H@Sh1CoR3!"

  storage_mb = 32768

  sku_name = "GP_Standard_D4s_v3"
}

resource "azurerm_postgresql_flexible_server_configuration" "example" {
  name      = "backslash_quote"
  server_id = azurerm_postgresql_flexible_server.example.id
  value     = "on"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the PostgreSQL Configuration, which needs [to be a valid PostgreSQL configuration name](https://www.postgresql.org/docs/current/static/sql-syntax-lexical.html#SQL-SYNTAX-IDENTIFIER). Changing this forces a new resource to be created.

* `server_id` - (Required) The ID of the PostgreSQL Flexible Server where we want to change configuration. Changing this forces a new PostgreSQL Flexible Server Configuration resource.

* `value` - (Required) Specifies the value of the PostgreSQL Configuration. See the PostgreSQL documentation for valid values.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the PostgreSQL Configuration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the PostgreSQL Configuration.
* `update` - (Defaults to 30 minutes) Used when updating the PostgreSQL Configuration.
* `read` - (Defaults to 5 minutes) Used when retrieving the PostgreSQL Configuration.
* `delete` - (Defaults to 30 minutes) Used when deleting the PostgreSQL Configuration.

## Import

PostgreSQL Configurations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_postgresql_flexible_server_configuration.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.DBforPostgreSQL/flexibleServers/server1/configurations/configuration1
```
