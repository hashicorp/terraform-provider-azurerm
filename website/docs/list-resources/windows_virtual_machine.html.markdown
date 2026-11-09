---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_windows_virtual_machine"
description: |-
    Lists Windows Virtual Machine resources.
---

# List resource: azurerm_windows_virtual_machine

Lists Windows Virtual Machine resources.

## Example Usage

### List all Windows Virtual Machines in the subscription

```hcl
list "azurerm_windows_virtual_machine" "example" {
  provider = azurerm
  config {}
}
```

### List all Windows Virtual Machines in a specific resource group

```hcl
list "azurerm_windows_virtual_machine" "example" {
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
