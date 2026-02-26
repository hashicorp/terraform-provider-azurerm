---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mysql_flexible_server"
description: |-
  Lists MySQL Flexible Servers.
---

# List resource: azurerm_mysql_flexible_server_firewall_rule

Lists MySQL Flexible Server Firewall Rule resources.

## Example Usage

### List all firewall rules deployed in a specific MySQL Flexible Server

```hcl
list "azurerm_mysql_flexible_server_firewall_rule" "example" {
  provider = azurerm
  config {
    flexible_server_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DBforMySQL/flexibleServers/flexibleServer1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `flexible_server_id` - (Required) The full ID of an existing Azure MySQL Flexible Server.

````
