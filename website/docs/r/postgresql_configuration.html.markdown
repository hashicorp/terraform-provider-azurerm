---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_postgresql_configuration"
description: |-
  Sets a PostgreSQL Configuration value on a PostgreSQL Server.
---

# azurerm_postgresql_configuration

Sets a PostgreSQL Configuration value on a PostgreSQL Server.

## Disclaimers

~> **Note:** Since this resource is provisioned by default, the Azure Provider will not check for the presence of an existing resource prior to attempting to create it.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "api-rg-pro"
  location = "West Europe"
}

resource "azurerm_postgresql_server" "example" {
  name                = "postgresql-server-1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku_name = "B_Gen5_2"

  storage_mb                   = 5120
  backup_retention_days        = 7
  geo_redundant_backup_enabled = false
  auto_grow_enabled            = true

  administrator_login          = "psqladminun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "9.5"
  ssl_enforcement_enabled      = true
}

resource "azurerm_postgresql_configuration" "example" {
  name                = "backslash_quote"
  resource_group_name = azurerm_resource_group.example.name
  server_name         = azurerm_postgresql_server.example.name
  value               = "on"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the PostgreSQL Configuration, which needs [to be a valid PostgreSQL configuration name](https://www.postgresql.org/docs/current/static/sql-syntax-lexical.html#SQL-SYNTAX-IDENTIFIER). Changing this forces a new resource to be created.

* `server_name` - (Required) Specifies the name of the PostgreSQL Server. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the PostgreSQL Server exists. Changing this forces a new resource to be created.

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
terraform import azurerm_postgresql_configuration.backslash_quote /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.DBforPostgreSQL/servers/server1/configurations/backslash_quote
```
