---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mysql_configuration"
sidebar_current: "docs-azurerm-resource-database-mysql-configuration"
description: |-
  Sets a MySQL Configuration value on a MySQL Server.
---

# azurerm_mysql_configuration

Sets a MySQL Configuration value on a MySQL Server.

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

  administrator_login = "psqladminun"
  administrator_login_password = "H@Sh1CoR3!"
  version = "5.7"
  storage_mb = "51200"
  ssl_enforcement = "Enabled"
}

resource "azurerm_mysql_configuration" "test" {
  name                = "interactive_timeout"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_mysql_server.test.name}"
  value               = "600"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the MySQL Configuration, which needs [to be a valid MySQL configuration name](https://dev.mysql.com/doc/refman/5.7/en/server-configuration.html). Changing this forces a new resource to be created.

* `server_name` - (Required) Specifies the name of the MySQL Server. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the MySQL Server exists. Changing this forces a new resource to be created.

* `value` - (Required) Specifies the value of the MySQL Configuration. See the MySQL documentation for valid values.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the MySQL Configuration.

## Import

MySQL Configurations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mysql_configuration.interactive_timeout /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.DBforMySQL/servers/server1/configurations/interactive_timeout
```
