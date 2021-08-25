---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_postgresql_flexible_server_firewall_rule"
description: |-
  Manages a PostgreSQL Flexible Server Firewall Rule.
---

# azurerm_postgresql_flexible_server_firewall_rule

Manages a PostgreSQL Flexible Server Firewall Rule.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_postgresql_flexible_server" "example" {
  name                   = "example-psqlflexibleserver"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  version                = "12"
  administrator_login    = "psqladminun"
  administrator_password = "H@Sh1CoR3!"

  storage_mb = 32768

  sku_name = "GP_Standard_D4s_v3"
}

resource "azurerm_postgresql_flexible_server_firewall_rule" "example" {
  name             = "example-fw"
  server_id        = azurerm_postgresql_flexible_server.example.id
  start_ip_address = "122.122.0.0"
  end_ip_address   = "122.122.0.0"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this PostgreSQL Flexible Server Firewall Rule. Changing this forces a new PostgreSQL Flexible Server Firewall Rule to be created.

* `server_id` - (Required) The ID of the PostgreSQL Flexible Server from which to create this PostgreSQL Flexible Server Firewall Rule. Changing this forces a new PostgreSQL Flexible Server Firewall Rule to be created.

* `start_ip_address` - (Required) The Start IP Address associated with this PostgreSQL Flexible Server Firewall Rule.

* `end_ip_address` - (Required) The End IP Address associated with this PostgreSQL Flexible Server Firewall Rule.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the PostgreSQL Flexible Server Firewall Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the PostgreSQL Flexible Server Firewall Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the PostgreSQL Flexible Server Firewall Rule.
* `update` - (Defaults to 30 minutes) Used when updating the PostgreSQL Flexible Server Firewall Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the PostgreSQL Flexible Server Firewall Rule.

## Import

PostgreSQL Flexible Server Firewall Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_postgresql_flexible_server_firewall_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DBforPostgreSQL/flexibleServers/flexibleServer1/firewallRules/firewallRule1
```
