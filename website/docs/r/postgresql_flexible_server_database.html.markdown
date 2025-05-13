---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_postgresql_flexible_server_database"
description: |-
  Manages a PostgreSQL Flexible Server Database.
---

# azurerm_postgresql_flexible_server_database

Manages a PostgreSQL Flexible Server Database.

!> **Note:** To mitigate the possibility of accidental data loss it is highly recommended that you use the `prevent_destroy` lifecycle argument in your configuration file for this resource. For more information on the `prevent_destroy` lifecycle argument please see the [terraform documentation](https://developer.hashicorp.com/terraform/tutorials/state/resource-lifecycle#prevent-resource-deletion).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_postgresql_flexible_server" "example" {
  name                   = "example-psqlflexibleserver"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  version                = "12"
  administrator_login    = "psqladmin"
  administrator_password = "H@Sh1CoR3!"
  storage_mb             = 32768
  sku_name               = "GP_Standard_D4s_v3"
}

resource "azurerm_postgresql_flexible_server_database" "example" {
  name      = "exampledb"
  server_id = azurerm_postgresql_flexible_server.example.id
  collation = "en_US.utf8"
  charset   = "UTF8"

  # prevent the possibility of accidental data loss
  lifecycle {
    prevent_destroy = true
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the PostgreSQL Database, which needs [to be a valid PostgreSQL identifier](https://www.postgresql.org/docs/current/sql-syntax-lexical.html#SQL-SYNTAX-IDENTIFIERS). Changing this forces a new Azure PostgreSQL Flexible Server Database to be created.

* `server_id` - (Required) The ID of the Azure PostgreSQL Flexible Server from which to create this PostgreSQL Flexible Server Database. Changing this forces a new Azure PostgreSQL Flexible Server Database to be created.

* `charset` - (Optional) Specifies the Charset for the Azure PostgreSQL Flexible Server Database, which needs [to be a valid PostgreSQL Charset](https://www.postgresql.org/docs/current/static/multibyte.html). Defaults to `UTF8`. Changing this forces a new Azure PostgreSQL Flexible Server Database to be created.

* `collation` - (Optional) Specifies the Collation for the Azure PostgreSQL Flexible Server Database, which needs [to be a valid PostgreSQL Collation](https://www.postgresql.org/docs/current/static/collation.html). Defaults to `en_US.utf8`. Changing this forces a new Azure PostgreSQL Flexible Server Database to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure PostgreSQL Flexible Server Database.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure PostgreSQL Flexible Server Database.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure PostgreSQL Flexible Server Database.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure PostgreSQL Flexible Server Database.

## Import

Azure PostgreSQL Flexible Server Database can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_postgresql_flexible_server_database.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.DBforPostgreSQL/flexibleServers/flexibleServer1/databases/database1
```
