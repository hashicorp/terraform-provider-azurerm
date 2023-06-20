---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_outbound_firewall_rule"
description: |-
  Manages an Azure SQL Outbound Firewall Rule.
---

# azurerm_mssql_outbound_firewall_rule

Allows you to manage an Azure SQL Outbound Firewall Rule.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_mssql_server" "example" {
  name                         = "mysqlserver"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  version                      = "12.0"
  administrator_login          = "4dm1n157r470r"
  administrator_login_password = "4-v3ry-53cr37-p455w0rd"

  outbound_network_restriction_enabled = true
}

resource "azurerm_mssql_outbound_firewall_rule" "example" {
  name      = "sqlexamplefdqn.database.windows.net"
  server_id = azurerm_mssql_server.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the outbound firewall rule. This should be a FQDN. Changing this forces a new resource to be created.

* `server_id` - (Required) The resource ID of the SQL Server on which to create the Outbound Firewall Rule. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The SQL Outbound Firewall Rule ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the SQL Outbound Firewall Rule.
* `update` - (Defaults to 30 minutes) Used when updating the SQL Outbound Firewall Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the SQL Outbound Firewall Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the SQL Outbound Firewall Rule.

## Import

SQL Outbound Firewall Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_outbound_firewall_rule.rule1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Sql/servers/myserver/outboundFirewallRules/fqdn1
```
