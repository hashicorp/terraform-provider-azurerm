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

  create_mode = "Default"
  administrator_login = "mysqladminun"
  administrator_login_password = "H@Sh1CoR3!"
  version = "5.7"
  ssl_enforcement = "Enabled"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the MySQL Server. Changing this forces a new resource to be created. This needs to be globally unique within Azure.

* `resource_group_name` - (Required) The name of the resource group in which to create the MySQL Server.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku` - (Required) A `sku` block as defined below.

* `storage_profile` - (Required) A `storage_profile` block as defined below.

* `create_mode` - (Optional) The mode to create a new server, supported values `Default` or `PointInTimeRestore`. 

* `administrator_login` - (Required) The Administrator Login for the MySQL Server. Changing this forces a new resource to be created.

* `administrator_login_password` - (Required) The Password associated with the `administrator_login` for the MySQL Server.

* `version` - (Required) Specifies the version of MySQL to use. Valid values are `5.6` and `5.7`. Changing this forces a new resource to be created.

* `ssl_enforcement` - (Required) Specifies if SSL should be enforced on connections. Possible values are `Enforced` and `Disabled`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

* `sku` supports the following:

* `name` - (Required) Specifies the SKU Name for this MySQL Server. See the `name` values by `tier` and `capacity` table below for valid values.
* `capacity` - (Required) The scale up/out capacity, representing server's compute units. Valid values depends on the `tier` of the server. See the `name` values by `tier` and `capacity` table below for valid values.
* `tier` - (Required) The tier of the particular SKU. Possible values are `Basic`, `GeneralPurpose`, and `MemoryOptimized`.
* `family` - (Required) The `family` of hardware, `Gen4` (Intel E5-2673 v3 (Haswell) 2.4 GHz processors) or
`Gen5` (Intel E5-2673 v4 (Broadwell) 2.3 GHz processors).


Suported `name` values by `tier` and `capacity`:

**`Tier`: Basic**
`name` | `family` | `capacity` | Storage Type
-- | -- | -- | --
`B_Gen4_1` | `Gen4` | 1 | Standard Storage
`B_Gen4_2` | `Gen4` | 2 | Standard Storage
`B_Gen5_1` | `Gen5` | 1 | Standard Storage
`B_Gen5_2` | `Gen5` | 2 | Standard Storage

**`Tier`: Gereral Purpose**
`name` | `family` | `capacity` | Storage Type
-- | -- | -- | --
`GP_Gen4_2` | `Gen4` | 2 | Premium Storage
`GP_Gen4_4` | `Gen4` | 4 | Premium Storage
`GP_Gen4_8` | `Gen4` | 8 | Premium Storage
`GP_Gen4_16` | `Gen4` | 16 | Premium Storage
`GP_Gen4_32` | `Gen4` | 32 | Premium Storage
`GP_Gen5_2` | `Gen5` | 2 | Premium Storage
`GP_Gen5_4` | `Gen5` | 4 | Premium Storage
`GP_Gen5_8` | `Gen5` | 8 | Premium Storage
`GP_Gen5_16` | `Gen5` | 16 | Premium Storage
`GP_Gen5_32` | `Gen5` | 32 | Premium Storage

**`Tier`: MemoryOptimized**
`name` | `family` | `capacity` | Storage Type
-- | -- | -- | --
`MO_Gen5_2` | `Gen5` | 2 | Premium Storage
`MO_Gen5_4` | `Gen5` | 4 | Premium Storage
`MO_Gen5_8` | `Gen5` | 8 | Premium Storage
`MO_Gen5_16` | `Gen5` | 16 | Premium Storage

---

* `storage_profile` supports the following:

* `storage_mb` - (Required) Max storage allowed for a server, possible values are between `5120` (5GB) and `1048576` (1TB). The step for this value must be in `1024` (1GB) increments.
* `backup_retention_days` - (Optional) Backup retention days for the server, supported values are between `7` and `35` days.
* `geo_redundant_backup` - (Optional) Enable Geo-redundant or not for server backup. Valid values for this property are `Enabled` or `Disabled`, not supported for the `basic` tier.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the MySQL Server.

* `fqdn` - The FQDN of the MySQL Server.

## Import

MySQL Server's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mysql_server.server1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.DBforMySQL/servers/server1
```
