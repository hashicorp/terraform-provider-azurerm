---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_postgresql_server"
sidebar_current: "docs-azurerm-resource-postgresql-server"
description: |-
  Create a PostgreSQL Server.
---

# azurerm\_postgresql\_server

Create a PostgreSQL Server.

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
    name = "PGSQLB50"
    capacity = 50
    tier = "Basic"
  }

  administrator_login = "psqladminun"
  administrator_login_password = "H@Sh1CoR3!"
  version = "9.5"
  storage_mb = "51200"
  ssl_enforcement = "Enabled"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the PostgreSQL Server. Changing this forces a
    new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to
    create the PostgreSQL Server.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku` - (Required) A `sku` block as defined below.

* `administrator_login` - (Required) The Administrator Login for the PostgreSQL Server. Changing this forces a new resource to be created.

* `administrator_login_password` - (Required) The Password associated with the `administrator_login` for the PostgreSQL Server.

* `version` - (Required) Specifies the version of PostgreSQL to use. Valid values are `9.5` and `9.6`. Changing this forces a new resource to be created.

* `storage_mb` - (Required) Specifies the amount of storage for the PostgreSQL Server in Megabytes. Possible values are `51200` and `102400`. Changing this forces a new resource to be created.

* `ssl_enforcement` - (Required) Specifies if SSL should be enforced on connections. Possible values are `Enforced` and `Disabled`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

* `sku` supports the following:

* `name` - (Optional) Specifies the SKU Name for this PostgreSQL Server. Possible values are `PGSQLB50` and `PGSQLB100`.
* `capacity` - (Optional) Specifies the DTU's for this PostgreSQL Server. Possible values are `50` and `100` DTU's.
* `tier` - (Optional) Specifies the SKU Tier for this PostgreSQL Server. At this time the only allowed value is `Basic`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the PostgreSQL Server.

* `fqdn` - The FQDN of the PostgreSQL Server.

## Import

PostgreSQL Server's can be imported using the `resource id`, e.g.

```
terraform import azurerm_postgresql_server.server1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.DBforPostgreSQL/servers/server1
```
