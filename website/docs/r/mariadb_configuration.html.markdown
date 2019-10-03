---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mariadb_configuration"
sidebar_current: "docs-azurerm-resource-database-mariadb-configuration"
description: |-
  Sets a MariaDB Configuration value on a MariaDB Server.
---

# azurerm_mariadb_configuration

Sets a MariaDB Configuration value on a MariaDB Server.

-> **NOTE** MariaDB Server is currently in Public Preview. You can find more information, including [how to register for the Public Preview here](https://azure.microsoft.com/en-us/updates/mariadb-public-preview/).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "api-rg-pro"
  location = "West Europe"
}

resource "azurerm_mariadb_server" "example" {
  name                = "mariadb-server-1"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"

  sku {
    name     = "B_Gen5_2"
    capacity = 2
    tier     = "Basic"
    family   = "Gen5"
  }

  storage_profile {
    storage_mb            = 5120
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "mariadbadmin"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "10.2"
  ssl_enforcement              = "Enabled"
}

resource "azurerm_mariadb_configuration" "example" {
  name                = "interactive_timeout"
  resource_group_name = "${azurerm_resource_group.example.name}"
  server_name         = "${azurerm_mariadb_server.example.name}"
  value               = "600"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the MariaDB Configuration, which needs [to be a valid MariaDB configuration name](https://mariadb.com/kb/en/library/server-system-variables/). Changing this forces a new resource to be created.

* `server_name` - (Required) Specifies the name of the MariaDB Server. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the MariaDB Server exists. Changing this forces a new resource to be created.

* `value` - (Required) Specifies the value of the MariaDB Configuration. See the MariaDB documentation for valid values.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the MariaDB Configuration.

## Import

MariaDB Configurations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mariadb_configuration.interactive_timeout /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.DBforMariaDB/servers/server1/configurations/interactive_timeout
```
