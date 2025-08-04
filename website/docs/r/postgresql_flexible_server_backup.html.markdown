---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_postgresql_flexible_server_backup"
description: |-
  Manages a PostgreSQL Flexible Server Backup.
---

# azurerm_postgresql_flexible_server_backup

Manages a PostgreSQL Flexible Server Backup.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_postgresql_flexible_server" "example" {
  name                   = "example-fs"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  administrator_login    = "adminTerraform"
  administrator_password = "QAZwsx123"
  version                = "12"
  sku_name               = "GP_Standard_D2s_v3"
  zone                   = "2"
}

resource "azurerm_postgresql_flexible_server_backup" "example" {
  name      = "example-pfsb"
  server_id = azurerm_postgresql_flexible_server.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of this PostgreSQL Flexible Server Backup. Changing this forces a new resource to be created.

* `server_id` - (Required) The ID of the PostgreSQL Flexible Server from which to create this PostgreSQL Flexible Server Backup. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the PostgreSQL Flexible Server Backup.

* `completed_time` - The Time (ISO8601 format) at which the backup was completed.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating this PostgreSQL Flexible Server Backup.
* `delete` - (Defaults to 30 minutes) Used when deleting this PostgreSQL Flexible Server Backup.
* `read` - (Defaults to 5 minutes) Used when retrieving this PostgreSQL Flexible Server Backup.

## Import

An existing PostgreSQL Flexible Server Backup can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_postgresql_flexible_server_backup.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DBforPostgreSQL/flexibleServers/fs1/backups/backup1
```
