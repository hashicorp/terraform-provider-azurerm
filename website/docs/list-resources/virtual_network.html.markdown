---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_network"
description: |-
  Lists Virtual Network resources.
---

# List resource: azurerm_virtual_network

Lists Virtual Network resources.

## Example Usage

```hcl
list "azurerm_storage_account" "test" {
  provider = azurerm
  config {
    resource_group_name = "example-rg"
  }
}
```

## Argument Reference

This list resource supports the following attributes:

* `resource_group_name` - (Required) The name of the resource group to query.
