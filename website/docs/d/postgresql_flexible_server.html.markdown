---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_postgresql_flexible_server"
description: |-
  Gets information about an existing PostgreSQL Flexible Server.
---

# Data Source: azurerm_postgresql_flexible_server

Use this data source to access information about an existing PostgreSQL Flexible Server.

## Example Usage

```hcl
data "azurerm_postgresql_flexible_server" "example" {
  name                = "existing-postgresql-fs"
  resource_group_name = "existing-postgresql-resgroup"
}

output "id" {
  value = data.azurerm_postgresql_flexible_server.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this PostgreSQL Flexible Server.

* `resource_group_name` - (Required) The name of the Resource Group where the PostgreSQL Flexible Server exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the PostgreSQL Flexible Server.

* `location` - The Azure Region where the PostgreSQL Flexible Server exists.

* `administrator_login` - The Administrator Login for the PostgreSQL Flexible Server.

* `backup_retention_days` -  The backup retention days for the PostgreSQL Flexible Server.

* `cmk_enabled` - The status showing whether the data encryption is enabled with a customer-managed key.

* `delegated_subnet_id` - The ID of the virtual network subnet to create the PostgreSQL Flexible Server.

* `fqdn` - The FQDN of the PostgreSQL Flexible Server.

* `public_network_access_enabled` - Is public network access enabled?

* `sku_name` - The SKU Name for the PostgreSQL Flexible Server. The name of the SKU, follows the `tier` + `name` pattern (e.g. `B_Standard_B1ms`, `GP_Standard_D2s_v3`, `MO_Standard_E4s_v3`).

* `storage_mb` - The max storage allowed for the PostgreSQL Flexible Server.

* `version` - The version of PostgreSQL Flexible Server to use.

* `tags` - A mapping of tags assigned to the PostgreSQL Flexible Server.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the PostgreSQL Flexible Server.
