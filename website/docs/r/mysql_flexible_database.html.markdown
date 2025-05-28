---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mysql_flexible_database"
description: |-
  Manages a MySQL Database within a MySQL Flexible Server.
---

# azurerm_mysql_flexible_database

Manages a MySQL Database within a MySQL Flexible Server

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_mysql_flexible_server" "example" {
  name                   = "example-mysql-flexible-server"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  administrator_login    = "mysqladminun"
  administrator_password = "H@Sh1CoR3!"
  sku_name               = "B_Standard_B1ms"
}

resource "azurerm_mysql_flexible_database" "example" {
  name                = "exampledb"
  resource_group_name = azurerm_resource_group.example.name
  server_name         = azurerm_mysql_flexible_server.example.name
  charset             = "utf8"
  collation           = "utf8_unicode_ci"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the MySQL Database, which needs [to be a valid MySQL identifier](https://dev.mysql.com/doc/refman/5.7/en/identifiers.html). Changing this forces a new resource to be created.

* `server_name` - (Required) Specifies the name of the MySQL Flexible Server. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the MySQL Server exists. Changing this forces a new resource to be created.

* `charset` - (Required) Specifies the Charset for the MySQL Database, which needs [to be a valid MySQL Charset](https://dev.mysql.com/doc/refman/5.7/en/charset-charsets.html). Changing this forces a new resource to be created.

* `collation` - (Required) Specifies the Collation for the MySQL Database, which needs [to be a valid MySQL Collation](https://dev.mysql.com/doc/refman/5.7/en/charset-mysql.html). Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the MySQL Database.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the MySQL Database.
* `read` - (Defaults to 5 minutes) Used when retrieving the MySQL Database.
* `delete` - (Defaults to 1 hour) Used when deleting the MySQL Database.

## Import

MySQL Database's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mysql_flexible_database.database1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.DBforMySQL/flexibleServers/flexibleserver1/databases/database1
```
