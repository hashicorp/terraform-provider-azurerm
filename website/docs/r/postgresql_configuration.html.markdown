---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_postgresql_configuration"
sidebar_current: "docs-azurerm-resource-database-postgresql-configuration"
description: |-
  Sets a PostgreSQL Configuration value on a PostgreSQL Server.
---

# azurerm_postgresql_configuration

Sets a PostgreSQL Configuration value on a PostgreSQL Server.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  # ...
}

resource "azurerm_postgresql_server" "example" {
  # ...
}

resource "azurerm_postgresql_configuration" "example" {
  name                = "backslash_quote"
  resource_group_name = "${azurerm_resource_group.example.name}"
  server_name         = "${azurerm_postgresql_server.example.name}"
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

## Import

PostgreSQL Configurations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_postgresql_configuration.backslash_quote /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.DBforPostgreSQL/servers/server1/configurations/backslash_quote
```
