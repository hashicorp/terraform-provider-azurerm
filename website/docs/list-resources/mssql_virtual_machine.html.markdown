---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_virtual_machine"
description: |-
    Lists mssql virtual machine resources.
---

# List resource: azurerm_mssql_virtual_machine

Lists mssql virtual machine resources.

## Example Usage

### List all mssql virtual machines 

```hcl
list "azurerm_mssql_virtual_machine" "example" {
  provider = azurerm
  config {
  }
}
```

### List all mssql virtual machines in a resource group

```hcl
list "azurerm_mssql_virtual_machine" "example" {
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

