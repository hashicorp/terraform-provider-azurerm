---
subcategory: "Spring Cloud"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_spring_cloud_app_mysql_association"
description: |-
  Associates a [Spring Cloud Application](spring_cloud_app.html) with a [MySQL Database](mysql_database.html).
---

# azurerm_spring_cloud_app_mysql_association

Associates a [Spring Cloud Application](spring_cloud_app.html) with a [MySQL Database](mysql_database.html).

!> **Note:** Azure Spring Apps is now deprecated and will be retired on 2028-05-31 - as such the `azurerm_spring_cloud_app_mysql_association` resource is deprecated and will be removed in a future major version of the AzureRM Provider. See https://aka.ms/asaretirement for more information.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_spring_cloud_service" "example" {
  name                = "example-springcloud"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_spring_cloud_app" "example" {
  name                = "example-springcloudapp"
  resource_group_name = azurerm_resource_group.example.name
  service_name        = azurerm_spring_cloud_service.example.name
}

resource "azurerm_mysql_flexible_server" "example" {
  name                   = "example-fsserver"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  sku_name               = "B_Standard_B1ms"
  zone                   = "2"
}

resource "azurerm_mysql_flexible_database" "example" {
  name                = "exampledb"
  resource_group_name = azurerm_resource_group.example.name
  server_name         = azurerm_mysql_flexible_server.example.name
  charset             = "utf8"
  collation           = "utf8_unicode_ci"
}

resource "azurerm_spring_cloud_app_mysql_association" "example" {
  name                = "example-bind"
  spring_cloud_app_id = azurerm_spring_cloud_app.example.id
  mysql_server_id     = azurerm_mysql_flexible_server.example.id
  database_name       = azurerm_mysql_flexible_database.example.name
  username            = azurerm_mysql_flexible_server.example.administrator_login
  password            = azurerm_mysql_flexible_server.example.administrator_login_password
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Spring Cloud Application Association. Changing this forces a new resource to be created.

* `spring_cloud_app_id` - (Required) Specifies the ID of the Spring Cloud Application where this Association is created. Changing this forces a new resource to be created.

* `mysql_server_id` - (Required) Specifies the ID of the MySQL Server. Changing this forces a new resource to be created.

* `database_name` - (Required) Specifies the name of the MySQL Database which the Spring Cloud App should be associated with.

* `username` - (Required) Specifies the username which should be used when connecting to the MySQL Database from the Spring Cloud App.

* `password` - (Required) Specifies the password which should be used when connecting to the MySQL Database from the Spring Cloud App.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Spring Cloud Application MySQL Association.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spring Cloud Application MySQL Association.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spring Cloud Application MySQL Association.
* `update` - (Defaults to 30 minutes) Used when updating the Spring Cloud Application MySQL Association.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spring Cloud Application MySQL Association.

## Import

Spring Cloud Application MySQL Association can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_spring_cloud_app_mysql_association.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourcegroup1/providers/Microsoft.AppPlatform/spring/service1/apps/app1/bindings/bind1
```
