---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_firewall_rule"
description: |-
  Manages an Azure SQL Firewall Rule.
---

# azurerm_mssql_firewall_rule

Allows you to manage an Azure SQL Firewall Rule.

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
}

resource "azurerm_mssql_firewall_rule" "example" {
  name             = "FirewallRule1"
  server_id        = azurerm_mssql_server.example.id
  start_ip_address = "10.0.17.62"
  end_ip_address   = "10.0.17.62"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the firewall rule. Changing this forces a new resource to be created.

* `server_id` - (Required) The resource ID of the SQL Server on which to create the Firewall Rule. Changing this forces a new resource to be created.

* `start_ip_address` - (Required) The starting IP address to allow through the firewall for this rule.

* `end_ip_address` - (Required) The ending IP address to allow through the firewall for this rule.

-> **Note:** The Azure feature `Allow access to Azure services` can be enabled by setting `start_ip_address` and `end_ip_address` to `0.0.0.0` which ([is documented in the Azure API Docs](https://docs.microsoft.com/rest/api/sql/firewallrules/createorupdate)).

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The SQL Firewall Rule ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the SQL Firewall Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the SQL Firewall Rule.
* `update` - (Defaults to 30 minutes) Used when updating the SQL Firewall Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the SQL Firewall Rule.

## Import

SQL Firewall Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_firewall_rule.rule1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Sql/servers/myserver/firewallRules/rule1
```
