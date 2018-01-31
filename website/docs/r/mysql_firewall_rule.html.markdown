---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mysql_firewall_rule"
sidebar_current: "docs-azurerm-resource-database-mysql-firewall-rule"
description: |-
  Creates a Firewall Rule for a MySQL Server.
---

# azurerm_mysql_firewall_rule

Creates a Firewall Rule for a MySQL Server

## Example Usage (Single IP Address)

```hcl
resource "azurerm_resource_group" "test" {
  name     = "api-rg-pro"
  location = "West Europe"
}

resource "azurerm_mysql_server" "test" {
  # ...
}

resource "azurerm_mysql_firewall_rule" "test" {
  name                = "office"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_mysql_server.test.name}"
  start_ip_address    = "40.112.8.12"
  end_ip_address      = "40.112.8.12"
}
```

## Example Usage (IP Range)

```hcl
resource "azurerm_resource_group" "test" {
  name     = "api-rg-pro"
  location = "West Europe"
}

resource "azurerm_mysql_server" "test" {
  # ...
}

resource "azurerm_mysql_firewall_rule" "test" {
  name                = "office"
  resource_group_name = "${azurerm_resource_group.test.name}"
  server_name         = "${azurerm_mysql_server.test.name}"
  start_ip_address    = "40.112.0.0"
  end_ip_address      = "40.112.255.255"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the MySQL Firewall Rule. Changing this forces a new resource to be created.

* `server_name` - (Required) Specifies the name of the MySQL Server. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the MySQL Server exists. Changing this forces a new resource to be created.

* `start_ip_address` - (Required) Specifies the Charset for the MySQL Database. Changing this forces a new resource to be created.

* `end_ip_address` - (Required) Specifies the End IP Address associated with this Firewall Rule. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the MySQL Firewall Rule.

## Import

MySQL Firewall Rule's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mysql_firewall_rule.rule1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.DBforMySQL/servers/server1/firewallRules/rule1
```
