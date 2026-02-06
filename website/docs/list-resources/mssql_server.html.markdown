---
subcategory: "mssql"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_server"
description: |-
    Lists mssql server resources.
---

# List resource: azurerm_mssql_server

Lists mssql server resources.

## Example Usage

### List all mssql servers 

```hcl
list "azurerm_mssql_server" "example" {
  provider = azurerm
  config {
  }
}
```

### List all mssql servers in a resource group

```hcl
list "azurerm_mssql_server" "example" {
  provider = azurerm
  config {
    resource_group_name = "example"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `subscription_id` - (Optional) The id of the mssql subscription to query.

* `resource_group_name` - (Optional) The name of the mssql resource group to query.

