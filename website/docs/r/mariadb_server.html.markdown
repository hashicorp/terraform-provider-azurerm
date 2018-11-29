---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mariadb_server"
sidebar_current: "docs-azurerm-resource-database-mariadb-server"
description: |-
  Manages a MariaDB Server.
---

# azurerm_mariadb_server

Manage a MariaDB Server.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "api-rg-pro"
  location = "West Europe"
}

resource "azurerm_mariadb_server" "test" {
  name                = "mariadb-server-1"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name     = "B_Gen5_2"
    capacity = 2
    tier     = "Basic"
  }

  storage_profile {
    storage_mb            = 5120
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
  }

  administrator_login          = "mariadbadmin"
  administrator_login_password = "H@Sh1CoR3!"
  ssl_enforcement              = "Enabled"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the MariaDB Server. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the MariaDB Server. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku` - (Required) A `sku` block as defined below.

* `storage_profile` - (Required) A `storage_profile` block as defined below.

* `administrator_login` - (Required) The Administrator Login for the MariaDB Server. Changing this forces a new resource to be created.

* `administrator_login_password` - (Required) The Password associated with the `administrator_login` for the MariaDB Server.

* `ssl_enforcement` - (Required) Specifies if SSL should be enforced on connections. Possible values are `Enabled` and `Disabled`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

`sku` supports the following:

* `name` - (Required) Specifies the SKU Name for this MariaDB Server. The name of the SKU, follows the `tier` + `family` + `cores` pattern (e.g. B_Gen5_1, GP_Gen5_8). For more information see the [product documentation](https://docs.microsoft.com/en-us/rest/api/mariadb/servers/create#sku).

* `capacity` - (Required) The scale up/out capacity, representing server's compute units.

* `tier` - (Required) The tier of the particular SKU. Possible values are `Basic`, `GeneralPurpose`, and `MemoryOptimized`. For more information see the [product documentation](https://docs.microsoft.com/en-us/azure/mariadb/concepts-pricing-tiers). 

  **NOTE:** `family` has been omitted from the SKU since MariaDB only supports `Gen5` hardware.
---

`storage_profile` supports the following:

* `storage_mb` - (Required) Max storage allowed for a server. Possible values are between `5120` MB(5GB) and `1048576` MB(1TB) for the Basic SKU and between `5120` MB(5GB) and `4194304` MB(4TB) for General Purpose/Memory Optimized SKUs. For more information see the [product documentation](https://docs.microsoft.com/en-us/rest/api/mariadb/servers/create#storageprofile).

* `backup_retention_days` - (Optional) Backup retention days for the server, supported values are between `7` and `35` days.

* `geo_redundant_backup` - (Optional) Enable Geo-redundant or not for server backup. Valid values for this property are `Enabled` or `Disabled`, not supported for the `basic` tier.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the MariaDB Server.

* `fqdn` - The FQDN of the MariaDB Server.

* `version` - The Version of the MariaDB Server.

## Import

MariaDB Server's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mariadb_server.server1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.DBforMariaDB/servers/server1
```
