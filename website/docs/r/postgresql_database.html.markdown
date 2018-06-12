---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_postgresql_database"
sidebar_current: "docs-azurerm-resource-database-postgresql-database"
description: |-
  Manages a PostgreSQL Database within a PostgreSQL Server.
---

# azurerm_postgresql_database

Manages a PostgreSQL Database within a PostgreSQL Server

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "api-rg-pro"
  location = "West Europe"
}

resource "azurerm_postgresql_server" "test" {
  name                = "postgresql-server-1"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name = "B_Gen4_2"
    capacity = 2
    tier = "Basic"
    family = "Gen4"
  }

  storage_profile {
    storage_mb = 5120
    backup_retention_days = 7
    geo_redundant_backup = "Disabled"
  }

  administrator_login = "psqladminun"
  administrator_login_password = "H@Sh1CoR3!"
  version = "9.5"
  ssl_enforcement = "Enabled"
}

resource "azurerm_postgresql_database" "test" {
  name                = "exampledb"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_postgresql_server.test.name}"
  charset             = "UTF8"
  collation           = "English_United States.1252"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the PostgreSQL Database, which needs [to be a valid PostgreSQL identifier](https://www.postgresql.org/docs/current/static/sql-syntax-lexical.html#SQL-SYNTAX-IDENTIFIERS). Changing this forces a
    new resource to be created.

* `server_name` - (Required) Specifies the name of the PostgreSQL Server. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the PostgreSQL Server exists. Changing this forces a new resource to be created.

* `charset` - (Required) Specifies the Charset for the PostgreSQL Database, which needs [to be a valid PostgreSQL Charset](https://www.postgresql.org/docs/current/static/multibyte.html). Changing this forces a new resource to be created.

* `collation` - (Required) Specifies the Collation for the PostgreSQL Database, which needs [to be a valid PostgreSQL Collation](https://www.postgresql.org/docs/current/static/collation.html). Note that Microsoft uses different [notation](https://msdn.microsoft.com/library/windows/desktop/dd373814.aspx) - en-US instead of en_US. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the PostgreSQL Database.

## Import

PostgreSQL Database's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_postgresql_database.database1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.DBforPostgreSQL/servers/server1/databases/database1
```
