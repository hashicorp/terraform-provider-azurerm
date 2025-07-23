---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mysql_flexible_server_configuration"
description: |-
  Sets a MySQL Flexible Server Configuration value on a MySQL Flexible Server.
---

# azurerm_mysql_flexible_server_configuration

Sets a MySQL Flexible Server Configuration value on a MySQL Flexible Server.

## Disclaimers

~> **Note:** Since this resource is provisioned by default, the Azure Provider will not check for the presence of an existing resource prior to attempting to create it.

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
  name                   = "example-fs"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  administrator_login    = "adminTerraform"
  administrator_password = "H@Sh1CoR3!"
  sku_name               = "GP_Standard_D2ds_v4"
}

resource "azurerm_mysql_flexible_server_configuration" "example" {
  name                = "interactive_timeout"
  resource_group_name = azurerm_resource_group.example.name
  server_name         = azurerm_mysql_flexible_server.example.name
  value               = "600"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the MySQL Flexible Server Configuration, which needs [to be a valid MySQL configuration name](https://dev.mysql.com/doc/refman/5.7/en/server-configuration.html). Changing this forces a new resource to be created.

* `server_name` - (Required) Specifies the name of the MySQL Flexible Server. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the MySQL Flexible Server exists. Changing this forces a new resource to be created.

* `value` - (Required) Specifies the value of the MySQL Flexible Server Configuration. See the MySQL documentation for valid values.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the MySQL Flexible Server Configuration.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the MySQL Flexible Server Configuration.
* `read` - (Defaults to 5 minutes) Used when retrieving the MySQL Flexible Server Configuration.
* `update` - (Defaults to 30 minutes) Used when updating the MySQL Flexible Server Configuration.
* `delete` - (Defaults to 30 minutes) Used when deleting the MySQL Flexible Server Configuration.

## Import

MySQL Flexible Server Configurations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mysql_flexible_server_configuration.interactive_timeout /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DBforMySQL/flexibleServers/flexibleServer1/configurations/interactive_timeout
```
