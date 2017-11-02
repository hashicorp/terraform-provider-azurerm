---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mysql_database"
sidebar_current: "docs-azurerm-resource-database-mysql-database"
description: |-
  Creates a MySQL Database within a MySQL Server.
---

# azurerm_mysql_database

Creates a MySQL Database within a MySQL Server

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "api-rg-pro"
  location = "West Europe"
}

resource "azurerm_mysql_server" "test" {
  name                = "mysql-server-1"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name = "MYSQLB50"
    capacity = 50
    tier = "Basic"
  }

  administrator_login = "mysqladminun"
  administrator_login_password = "H@Sh1CoR3!"
  version = "5.7"
  storage_mb = "51200"
  ssl_enforcement = "Enabled"
}

resource "azurerm_mysql_database" "test" {
  name                = "exampledb"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_mysql_server.test.name}"
  charset             = "utf8"
  collation           = "utf8_unicode_ci"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the MySQL Database, which needs [to be a valid MySQL identifier](https://dev.mysql.com/doc/refman/5.7/en/identifiers.html). Changing this forces a new resource to be created.

* `server_name` - (Required) Specifies the name of the MySQL Server. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the MySQL Server exists. Changing this forces a new resource to be created.

* `charset` - (Required) Specifies the Charset for the MySQL Database, which needs [to be a valid MySQL Charset](https://dev.mysql.com/doc/refman/5.7/en/charset-charsets.html). Changing this forces a new resource to be created.

* `collation` - (Required) Specifies the Collation for the MySQL Database, which needs [to be a valid MySQL Collation](https://dev.mysql.com/doc/refman/5.7/en/charset-mysql.html). Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the MySQL Database.

## Import

MySQL Database's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mysql_database.database1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.DBforMySQL/servers/server1/databases/database1
```
