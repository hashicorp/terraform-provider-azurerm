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

  administrator_login          = "mysqladminun"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "5.7"
  ssl_enforcement              = "Enabled"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the MySQL Server. Changing this forces a new resource to be created. This needs to be globally unique within Azure.

* `resource_group_name` - (Required) The name of the resource group in which to create the MySQL Server. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku` - (Required) A `sku` block as defined below.

* `storage_profile` - (Required) A `storage_profile` block as defined below.

* `administrator_login` - (Required) The Administrator Login for the MySQL Server. Changing this forces a new resource to be created.

* `administrator_login_password` - (Required) The Password associated with the `administrator_login` for the MySQL Server.

* `version` - (Required) Specifies the version of MySQL to use. Valid values are `5.6` and `5.7`. Changing this forces a new resource to be created.

* `ssl_enforcement` - (Required) Specifies if SSL should be enforced on connections. Possible values are `Enabled` and `Disabled`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

`sku` supports the following:

* `name` - (Required) Specifies the SKU Name for this MySQL Server. The name of the SKU, follows the `tier` + `family` + `cores` pattern (e.g. B_Gen4_1, GP_Gen5_8). For more information see the [product documentation](https://docs.microsoft.com/en-us/rest/api/mysql/servers/create#sku).

* `capacity` - (Required) The scale up/out capacity, representing server's compute units.

* `tier` - (Required) The tier of the particular SKU. Possible values are `Basic`, `GeneralPurpose`, and `MemoryOptimized`. For more information see the [product documentation](https://docs.microsoft.com/en-us/azure/mysql/concepts-pricing-tiers).

* `family` - (Required) The `family` of hardware `Gen4` or `Gen5`, before selecting your `family` check the [product documentation](https://docs.microsoft.com/en-us/azure/mysql/concepts-pricing-tiers#compute-generations-vcores-and-memory) for availability in your region.

---

`storage_profile` supports the following:

* `storage_mb` - (Required) Max storage allowed for a server. Possible values are between `5120` MB(5GB) and `1048576` MB(1TB) for the Basic SKU and between `5120` MB(5GB) and `4194304` MB(4TB) for General Purpose/Memory Optimized SKUs. For more information see the [product documentation](https://docs.microsoft.com/en-us/rest/api/mysql/servers/create#StorageProfile).

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
