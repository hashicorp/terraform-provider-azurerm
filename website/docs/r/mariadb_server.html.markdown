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
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_mariadb_server" "example" {
  name                = "example-mariadb-server"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  administrator_login          = "mariadbadmin"
  administrator_login_password = "H@Sh1CoR3!"

  sku_name   = "B_Gen5_2"
  storage_mb = 5120
  version    = "10.2"

  auto_grow_enabled             = true
  backup_retention_days         = 7
  geo_redundant_backup_enabled  = false
  public_network_access_enabled = false
  ssl_enforcement_enabled       = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the MariaDB Server. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the MariaDB Server. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku_name` - (Required) Specifies the SKU Name for this MariaDB Server. The name of the SKU, follows the `tier` + `family` + `cores` pattern (e.g. `B_Gen4_1`, `GP_Gen5_8`). For more information see the [product documentation](https://docs.microsoft.com/en-us/rest/api/mariadb/servers/create#sku).

* `version` - (Required) Specifies the version of MariaDB to use. Possible values are `10.2` and `10.3`. Changing this forces a new resource to be created.

* `administrator_login` - (Required) The Administrator Login for the MariaDB Server. Changing this forces a new resource to be created.

* `administrator_login_password` - (Required) The Password associated with the `administrator_login` for the MariaDB Server.

* `auto_grow_enabled` - (Optional) Enable/Disable auto-growing of the storage. Storage auto-grow prevents your server from running out of storage and becoming read-only. If storage auto grow is enabled, the storage automatically grows without impacting the workload. The default value if not explicitly specified is `true`.

* `backup_retention_days` - (Optional) Backup retention days for the server, supported values are between `7` and `35` days.

* `create_mode` - (Optional) The creation mode. Can be used to restore or replicate existing servers. Possible values are `Default`, `Replica`, `GeoRestore`, and `PointInTimeRestore`. Defaults to `Default`.

* `creation_source_server_id` - (Optional) For creation modes other than `Default`, the source server ID to use.

* `geo_redundant_backup_enabled` - (Optional) Turn Geo-redundant server backups on/off. This allows you to choose between locally redundant or geo-redundant backup storage in the General Purpose and Memory Optimized tiers. When the backups are stored in geo-redundant backup storage, they are not only stored within the region in which your server is hosted, but are also replicated to a paired data center. This provides better protection and ability to restore your server in a different region in the event of a disaster. This is not supported for the Basic tier.

* `public_network_access_enabled` - (Optional) Whether or not public network access is allowed for this server. Defaults to `true`.

* `restore_point_in_time` - (Optional) When `create_mode` is `PointInTimeRestore`, specifies the point in time to restore from `creation_source_server_id`.

* `ssl_enforcement_enabled` - (Required) Specifies if SSL should be enforced on connections. Possible values are `true` and `false`.

* `storage_mb` - (Required) Max storage allowed for a server. Possible values are between `5120` MB (5GB) and `1024000`MB (1TB) for the Basic SKU and between `5120` MB (5GB) and `4096000` MB (4TB) for General Purpose/Memory Optimized SKUs. For more information see the [product documentation](https://docs.microsoft.com/en-us/rest/api/mariadb/servers/create#storageprofile).

* `tags` - (Optional) A mapping of tags to assign to the resource.

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
