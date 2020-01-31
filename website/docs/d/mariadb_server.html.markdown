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

* `name` - (Required) The name of the MariaDB Server to retrieve information about.

* `resource_group_name` - (Required) The name of the resource group where the MariaDB Server exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the MariaDB Server.

* `fqdn` - The FQDN of the MariaDB Server.

* `location` - The Azure location where the resource exists.

* `sku_name` - The SKU Name for this MariaDB Server. 

* `storage_profile` - A `storage_profile` block as defined below.

* `administrator_login` - The Administrator Login for the MariaDB Server.

* `administrator_login_password` - The password associated with the `administrator_login` for the MariaDB Server.

* `version` - The version of MariaDB being used.

* `ssl_enforcement` - The SSL being enforced on connections.

* `tags` - A mapping of tags assigned to the resource.
---

A `storage_profile` block exports the following:

* `storage_mb` - The max storage allowed for a server.

* `backup_retention_days` - Backup retention days for the server.

* `geo_redundant_backup` - Whether Geo-redundant is enabled or not for server backup.

* `auto_grow` - Whether autogrow is enabled or disabled for the storage.

### Timeouts

~> **Note:** Custom Timeouts are available [as an opt-in Beta in version 1.43 of the Azure Provider](/docs/providers/azurerm/guides/2.0-beta.html) and will be enabled by default in version 2.0 of the Azure Provider.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the   Gets information about a MariaDB Server.
