---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_server"
description: |-
    Lists MSSQL Server resources.
---

# List resource: azurerm_mssql_server

Lists MSSQL Server resources.

## Example Usage

### List all MSSQL Servers in the subscription

```hcl
list "azurerm_mssql_server" "example" {
  provider = azurerm
  config {}
}
```

### List all MSSQL Servers in a specific resource group

```hcl
list "azurerm_mssql_server" "example" {
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

