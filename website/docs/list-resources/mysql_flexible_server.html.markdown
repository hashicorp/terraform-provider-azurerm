---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mysql_flexible_server"
description: |-
  Lists MySQL Flexible Servers.
---

# List resource: azurerm_mysql_flexible_server

Lists MySQL Flexible Server resources.

## Example Usage

### List all MySQL Flexible Servers in the subscription

```hcl
list "azurerm_mysql_flexible_server" "example" {
  provider = azurerm
  config {}
}
```

### List all MySQL Flexible Servers in a specific resource group

```hcl
list "azurerm_mysql_flexible_server" "example" {
  provider = azurerm
  config {
    resource_group_name = "example-rg"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `resource_group_name` - (Optional) The name of the resource group to query.

* `subscription_id` - (Optional) The Subscription ID to query. Defaults to the value specified in the Provider Configuration.
````
