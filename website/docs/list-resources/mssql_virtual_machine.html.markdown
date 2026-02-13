---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_virtual_machine"
description: |-
    Lists Mssql Virtual Machine resources.
---

# List resource: azurerm_mssql_virtual_machine

Lists Mssql Virtual Machine resources.

## Example Usage

### List all Mssql Virtual Machines

```hcl
list "azurerm_mssql_virtual_machine" "example" {
  provider = azurerm
  config {
  }
}
```

### List all Mssql Virtual Machines in a Resource Group

```hcl
list "azurerm_mssql_virtual_machine" "example" {
  provider = azurerm
  config {
    resource_group_name = "resource_group_name-example"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `subscription_id` - (Optional) The ID of the Subscription to query.

* `resource_group_name` - (Optional) The name of the Resource Group to query.
