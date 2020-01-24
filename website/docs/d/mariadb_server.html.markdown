---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mariadb_server"
description: |-
  Gets information about a MariaDB Server.
---

# Data Source: azurerm_mariadb_server

Use this data source to access information about a MariaDB Server.

## Example Usage

```hcl
data "azurerm_mariadb_server" "db_server" {
  name                = "mariadb-server"
  resource_group_name = "${azurerm_mariadb_server.example.resource_group_name}"
}
output "mariadb_server_id" {
  value = "${data.azurerm_mariadb_server.example.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the MariaDB Server to retrieve information about.

* `resource_group_name` - (Required) The name of the resource group where the MariaDB Server exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the MariaDB Server.

* `fqdn` - The FQDN of the MariaDB Server.

* `location` - The Azure location where the resource exists.

* `sku_name` - Specifies the SKU Name for this MariaDB Server. The name of the SKU, follows the `tier` + `family` + `cores` pattern (e.g. `B_Gen4_1`, `GP_Gen5_8`). For more information see the [product documentation](https://docs.microsoft.com/en-us/rest/api/mariadb/servers/create#sku).

* `storage_profile` - A `storage_profile` block as defined below.

* `administrator_login` - The Administrator Login for the MariaDB Server.

* `administrator_login_password` - The Password associated with the `administrator_login` for the MariaDB Server.

* `version` - Specifies the version of MariaDB being used. Possible values are `10.2` and `10.3`.

* `ssl_enforcement` - Specifies the SSL being enforced on connections. Possible values are `Enabled` and `Disabled`.

* `tags` - A mapping of tags assigned to the resource.
---

A `storage_profile` block exports the following:

* `storage_mb` - Max storage allowed for a server. Possible values are between `5120` MB (5GB) and `1024000`MB (1TB) for the Basic SKU and between `5120` MB (5GB) and `4096000` MB (4TB) for General Purpose/Memory Optimized SKUs. For more information see the [product documentation](https://docs.microsoft.com/en-us/rest/api/mariadb/servers/create#storageprofile).

* `backup_retention_days` - Backup retention days for the server, supported values are between `7` and `35` days.

* `geo_redundant_backup` - Specifies whether Geo-redundant is enabled or not for server backup. Valid values for this property are `Enabled` or `Disabled`.

* `auto_grow` - Specifies whether autogrow is enabled or disabled for the storage. Valid values are `Enabled` or `Disabled`.
