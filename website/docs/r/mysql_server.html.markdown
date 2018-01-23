---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mysql_server"
sidebar_current: "docs-azurerm-resource-database-mysql-server"
description: |-
  Manages a MySQL Server.

---

# azurerm_mysql_server

Manages a MySQL Server.

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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the MySQL Server. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the MySQL Server.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku` - (Required) A `sku` block as defined below.

* `administrator_login` - (Required) The Administrator Login for the MySQL Server. Changing this forces a new resource to be created.

* `administrator_login_password` - (Required) The Password associated with the `administrator_login` for the MySQL Server.

* `version` - (Required) Specifies the version of MySQL to use. Valid values are `5.6` and `5.7`. Changing this forces a new resource to be created.

* `storage_mb` - (Required) Specifies the amount of storage for the MySQL Server in Megabytes. Possible values are shown below. Changing this forces a new resource to be created.

Possible values for `storage_mb` when using a SKU Name of `Basic` are:
- `51200` (50GB)
- `179200` (175GB)
- `307200` (300GB)
- `435200` (425GB)
- `563200` (550GB)
- `691200` (675GB)
- `819200` (800GB)
- `947200` (925GB)

Possible values for `storage_mb` when using a SKU Name of `Standard` are:
- `128000` (125GB)
- `256000` (256GB)
- `384000` (384GB)
- `512000` (512GB)
- `640000` (640GB)
- `768000` (768GB)
- `896000` (896GB)
- `1024000` (1TB)

* `ssl_enforcement` - (Required) Specifies if SSL should be enforced on connections. Possible values are `Enforced` and `Disabled`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

* `sku` supports the following:

* `name` - (Optional) Specifies the SKU Name for this MySQL Server. Possible values are: `MYSQLB50`, `MYSQLB100`, `MYSQLS100`, `MYSQLS200`, `MYSQLS400` and `MYSQLS800`.
* `capacity` - (Optional) Specifies the DTU's for this MySQL Server. Possible values are `50` and `100` DTU's when using a `Basic` SKU and `100`, `200`, `400` or `800` when using the `Standard` SKU.
* `tier` - (Optional) Specifies the SKU Tier for this MySQL Server. Possible values are `Basic` and `Standard`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the MySQL Server.

* `fqdn` - The FQDN of the MySQL Server.

## Import

MySQL Server's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mysql_server.server1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.DBforMySQL/servers/server1
```
