---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mariadb_server"
description: |-
  Manages a MariaDB Server.
---

# azurerm_mariadb_server

Manages a MariaDB Server.

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

  sku_name = "B_Gen5_2"

  storage_profile {
    storage_mb            = 5120
    backup_retention_days = 7
    geo_redundant_backup  = "Disabled"
    storage_autogrow      = "Disabled"
  }

  administrator_login          = "mariadbadmin"
  administrator_login_password = "H@Sh1CoR3!"
  version                      = "10.2"
  ssl_enforcement              = "Enabled"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the MariaDB Server. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the MariaDB Server. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku_name` - (Required) Specifies the SKU Name for this MariaDB Server. The name of the SKU, follows the `tier` + `family` + `cores` pattern (e.g. `B_Gen4_1`, `GP_Gen5_8`). For more information see the [product documentation](https://docs.microsoft.com/en-us/rest/api/mariadb/servers/create#sku).

* `storage_profile` - (Required) A `storage_profile` block as defined below.

* `administrator_login` - (Required) The Administrator Login for the MariaDB Server. Changing this forces a new resource to be created.

* `administrator_login_password` - (Required) The Password associated with the `administrator_login` for the MariaDB Server.

* `version` - (Required) Specifies the version of MariaDB to use. Possible values are `10.2` and `10.3`. Changing this forces a new resource to be created.

* `ssl_enforcement` - (Required) Specifies if SSL should be enforced on connections. Possible values are `Enabled` and `Disabled`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `storage_profile` block supports the following:

* `storage_mb` - (Required) Max storage allowed for a server. Possible values are between `5120` MB (5GB) and `1024000`MB (1TB) for the Basic SKU and between `5120` MB (5GB) and `4096000` MB (4TB) for General Purpose/Memory Optimized SKUs. For more information see the [product documentation](https://docs.microsoft.com/en-us/rest/api/mariadb/servers/create#storageprofile).

* `backup_retention_days` - (Optional) Backup retention days for the server, supported values are between `7` and `35` days.

* `geo_redundant_backup` - (Optional) Enable Geo-redundant or not for server backup. Valid values for this property are `Enabled` or `Disabled`.

* `auto_grow` - (Optional) Defines whether autogrow is enabled or disabled for the storage. Valid values are `Enabled` or `Disabled`.

-> **NOTE:** Geo Redundant Backups cannot be configured when using the `Basic` tier.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the MariaDB Server.

* `fqdn` - The FQDN of the MariaDB Server.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the MariaDB Server.
* `update` - (Defaults to 60 minutes) Used when updating the MariaDB Server.
* `read` - (Defaults to 5 minutes) Used when retrieving the MariaDB Server.
* `delete` - (Defaults to 60 minutes) Used when deleting the MariaDB Server.

## Import

MariaDB Server's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mariadb_server.server1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.DBforMariaDB/servers/server1
```
