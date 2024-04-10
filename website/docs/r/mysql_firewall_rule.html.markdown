---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mysql_firewall_rule"
description: |-
  Manages a Firewall Rule for a MySQL Server.
---

# azurerm_mysql_firewall_rule

Manages a Firewall Rule for a MySQL Server.

~> **Note:** Azure Database for MySQL Single Server and its sub resources are scheduled for retirement by 2024-09-16 and will migrate to using Azure Database for MySQL Flexible Server: https://go.microsoft.com/fwlink/?linkid=2216041. The `azurerm_mysql_firewall_rule` resource is deprecated and will be removed in v4.0 of the AzureRM Provider. Please use the `azurerm_mysql_flexible_server_firewall_rule` resource instead.

## Example Usage (Single IP Address)

```hcl
resource "azurerm_resource_group" "example" {
  name     = "api-rg-pro"
  location = "West Europe"
}

resource "azurerm_mysql_server" "example" {
  name                    = "example"
  location                = azurerm_resource_group.example.location
  resource_group_name     = azurerm_resource_group.example.name
  version                 = "5.7"
  sku_name                = "GP_Gen5_2"
  ssl_enforcement_enabled = true
}

resource "azurerm_mysql_firewall_rule" "example" {
  name                = "office"
  resource_group_name = azurerm_resource_group.example.name
  server_name         = azurerm_mysql_server.example.name
  start_ip_address    = "40.112.8.12"
  end_ip_address      = "40.112.8.12"
}
```

## Example Usage (IP Range)

```hcl
resource "azurerm_resource_group" "example" {
  name     = "api-rg-pro"
  location = "West Europe"
}

resource "azurerm_mysql_server" "example" {
  # ...
}

resource "azurerm_mysql_firewall_rule" "example" {
  name                = "office"
  resource_group_name = azurerm_resource_group.example.name
  server_name         = azurerm_mysql_server.example.name
  start_ip_address    = "40.112.0.0"
  end_ip_address      = "40.112.255.255"
}
```

## Example Usage (Allow access to Azure services)

```hcl
resource "azurerm_resource_group" "example" {
  name     = "api-rg-pro"
  location = "West Europe"
}

resource "azurerm_mysql_server" "example" {
  # ...
}

resource "azurerm_mysql_firewall_rule" "example" {
  name                = "office"
  resource_group_name = azurerm_resource_group.example.name
  server_name         = azurerm_mysql_server.example.name
  start_ip_address    = "0.0.0.0"
  end_ip_address      = "0.0.0.0"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the MySQL Firewall Rule. Changing this forces a new resource to be created.

* `server_name` - (Required) Specifies the name of the MySQL Server. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the MySQL Server exists. Changing this forces a new resource to be created.

* `start_ip_address` - (Required) Specifies the Start IP Address associated with this Firewall Rule. 

* `end_ip_address` - (Required) Specifies the End IP Address associated with this Firewall Rule. 

-> **NOTE:** The Azure feature `Allow access to Azure services` can be enabled by setting `start_ip_address` and `end_ip_address` to `0.0.0.0` which ([is documented in the Azure API Docs](https://docs.microsoft.com/rest/api/sql/firewallrules/createorupdate)).

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the MySQL Firewall Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the MySQL Firewall Rule.
* `update` - (Defaults to 30 minutes) Used when updating the MySQL Firewall Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the MySQL Firewall Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the MySQL Firewall Rule.

## Import

MySQL Firewall Rule's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mysql_firewall_rule.rule1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.DBforMySQL/servers/server1/firewallRules/rule1
```
