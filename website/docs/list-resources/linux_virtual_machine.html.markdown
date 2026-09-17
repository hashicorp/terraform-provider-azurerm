---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_linux_virtual_machine"
description: |-
    Lists Linux Virtual Machine resources.
---

# List resource: azurerm_linux_virtual_machine

Lists Linux Virtual Machine resources.

## Example Usage

### List all Linux Virtual Machines in the subscription

```hcl
list "azurerm_linux_virtual_machine" "example" {
  provider = azurerm
  config {}
}
```

### List all Linux Virtual Machines in a specific resource group

```hcl
list "azurerm_linux_virtual_machine" "example" {
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
