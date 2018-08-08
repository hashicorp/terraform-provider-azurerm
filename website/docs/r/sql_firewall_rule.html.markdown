---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sql_firewall_rule"
sidebar_current: "docs-azurerm-resource-database-sql-firewall_rule"
description: |-
  Manages a SQL Firewall Rule.
---

# azurerm_sql_firewall_rule

Manages a SQL Firewall Rule

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  # ...
}

resource "azurerm_sql_server" "example" {
  # ...
}

resource "azurerm_sql_firewall_rule" "example" {
  name                = "example-rule"
  resource_group_name = "${azurerm_resource_group.example.name}"
  server_name         = "${azurerm_sql_server.example.name}"
  start_ip_address    = "10.0.17.62"
  end_ip_address      = "10.0.17.62"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the firewall rule.

* `resource_group_name` - (Required) The name of the resource group in which to create the sql server.

* `server_name` - (Required) The name of the SQL Server on which to create the Firewall Rule.

* `start_ip_address` - (Required) The starting IP address to allow through the firewall for this rule.

* `end_ip_address` - (Required) The ending IP address to allow through the firewall for this rule.

-> **NOTE:** The Azure feature `Allow access to Azure services` can be enabled by setting `start_ip_address` and `end_ip_address` to `0.0.0.0` which ([is documented in the Azure API Docs](https://docs.microsoft.com/en-us/rest/api/sql/firewallrules/createorupdate)).

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the SQL Firewall Rule.

## Import

SQL Firewall Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sql_firewall_rule.rule1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Sql/servers/myserver/firewallRules/rule1
```
