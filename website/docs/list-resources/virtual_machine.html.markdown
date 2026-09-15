---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_machine"
description: |-
    Lists Virtual Machine resources.
---

# List resource: azurerm_virtual_machine

Lists Virtual Machine resources.

## Example Usage

### List all Virtual Machines in the subscription

```hcl
list "azurerm_virtual_machine" "example" {
  provider = azurerm
  config {}
}
```

### List all Virtual Machines in a specific resource group

```hcl
list "azurerm_virtual_machine" "example" {
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
