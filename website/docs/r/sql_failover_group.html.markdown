---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sql_failover_group"
description: |-
  Manages a SQL Failover Group.
---

# azurerm_sql_failover_group

Create a failover group of databases on a collection of Azure SQL servers.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "uksouth"
}

resource "azurerm_sql_server" "primary" {
  name                         = "sql-primary"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  version                      = "12.0"
  administrator_login          = "sqladmin"
  administrator_login_password = "pa$$w0rd"
}

resource "azurerm_sql_server" "secondary" {
  name                         = "sql-secondary"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = "northeurope"
  version                      = "12.0"
  administrator_login          = "sqladmin"
  administrator_login_password = "pa$$w0rd"
}

resource "azurerm_sql_database" "db1" {
  name                = "db1"
  resource_group_name = azurerm_sql_server.primary.resource_group_name
  location            = azurerm_sql_server.primary.location
  server_name         = azurerm_sql_server.primary.name
}

resource "azurerm_sql_failover_group" "example" {
  name                = "example-failover-group"
  resource_group_name = azurerm_sql_server.primary.resource_group_name
  server_name         = azurerm_sql_server.primary.name
  databases           = [azurerm_sql_database.db1.id]
  partner_servers {
    id = azurerm_sql_server.secondary.id
  }

  read_write_endpoint_failover_policy {
    mode          = "Automatic"
    grace_minutes = 60
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the failover group. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group containing the SQL server

* `server_name` - (Required) The name of the primary SQL server. Changing this forces a new resource to be created.

* `databases` - A list of database ids to add to the failover group

-> **NOTE:** The failover group will create a secondary database for each database listed in `databases`. If the secondary databases need to be managed through Terraform, they should be defined as resources and a dependency added to the failover group to ensure the secondary databases are created first. Please refer to the detailed example which can be found in [the `./examples/sql-azure/failover_group` directory within the Github Repository](https://github.com/terraform-providers/terraform-provider-azurerm/tree/master/examples/sql-azure/failover_group)

* `partner_servers` - (Required) A list of secondary servers as documented below

* `read_write_endpoint_failover_policy` - (Required) A read/write policy as documented below

* `readonly_endpoint_failover_policy` - (Optional) a read-only policy as documented below

* `tags` - (Optional) A mapping of tags to assign to the resource.

`partner_servers` supports the following:

* `id` - (Required) the SQL server ID

`read_write_endpoint_failover_policy` supports the following:

* `mode` - (Required) the failover mode. Possible values are `Manual`, `Automatic`

* `grace_minutes` - Applies only if `mode` is `Automatic`. The grace period in minutes before failover with data loss is attempted

`readonly_endpoint_failover_policy` supports the following:

* `mode` - Failover policy for the read-only endpoint. Possible values are `Enabled`, and `Disabled`

## Attribute Reference

The following attributes are exported:

* `id` - The failover group ID.
* `location` - the location of the failover group.
* `server_name` - the name of the primary SQL Database Server.
* `role` - local replication role of the failover group instance.
* `databases` - list of databases in the failover group.
* `partner_servers` - list of partner server information for the failover group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the SQL Failover Group.
* `update` - (Defaults to 30 minutes) Used when updating the SQL Failover Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the SQL Failover Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the SQL Failover Group.

## Import

SQL Failover Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sql_failover_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Sql/servers/myserver/failovergroups/group1
```
