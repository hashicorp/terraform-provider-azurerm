---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_failover_group"
description: |-
  Manages a Microsoft Azure SQL Failover Group.

---

# azurerm_mssql_failover_group

Manages a Microsoft Azure SQL Failover Group.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "database-rg"
  location = "West Europe"
}

resource "azurerm_mssql_server" "primary" {
  name                         = "mssqlserver-primary"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  version                      = "12.0"
  administrator_login          = "missadministrator"
  administrator_login_password = "thisIsKat11"
}

resource "azurerm_mssql_server" "secondary" {
  name                         = "mssqlserver-secondary"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = "North Europe"
  version                      = "12.0"
  administrator_login          = "missadministrator"
  administrator_login_password = "thisIsKat12"
}

resource "azurerm_mssql_database" "example" {
  name        = "exampledb"
  server_id   = azurerm_mssql_server.primary.id
  sku_name    = "S1"
  collation   = "SQL_Latin1_General_CP1_CI_AS"
  max_size_gb = "200"
}

resource "azurerm_mssql_failover_group" "example" {
  name      = "example"
  server_id = azurerm_mssql_server.primary.id
  databases = [
    azurerm_mssql_database.example.id
  ]

  partner_server {
    id = azurerm_mssql_server.secondary.id
  }

  read_write_endpoint_failover_policy {
    mode          = "Automatic"
    grace_minutes = 80
  }

  tags = {
    environment = "prod"
    database    = "example"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Failover Group. Changing this forces a new resource to be created.

* `server_id` - (Required) The ID of the primary SQL Server on which to create the failover group. Changing this forces a new resource to be created.

* `partner_server` - (Required) A `partner_server` block as defined below.

* `databases` - (Optional) A set of database names to include in the failover group.

* `readonly_endpoint_failover_policy_enabled` - (Optional) Whether failover is enabled for the readonly endpoint. Defaults to `false`.

* `read_write_endpoint_failover_policy` - (Required) A `read_write_endpoint_failover_policy` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `partner_server` block supports the following:

* `id` - (Required) The ID of a partner SQL server to include in the failover group.

---

The `read_write_endpoint_failover_policy` block supports the following:

* `mode` - (Required) The failover policy of the read-write endpoint for the failover group. Possible values are `Automatic` or `Manual`.

* `grace_minutes` - (Optional) The grace period in minutes, before failover with data loss is attempted for the read-write endpoint. Required when `mode` is `Automatic`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Failover Group.

* `partner_server` - A `partner_server` block as defined below.

---

A `partner_server` block exports the following:

* `location` - The location of the partner server.

* `role` - The replication role of the partner server. Possible values include `Primary` or `Secondary`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Failover Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Failover Group.
* `update` - (Defaults to 30 minutes) Used when updating the Failover Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Failover Group.

## Import

Failover Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_failover_group.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Sql/servers/server1/failoverGroups/failoverGroup1
```
