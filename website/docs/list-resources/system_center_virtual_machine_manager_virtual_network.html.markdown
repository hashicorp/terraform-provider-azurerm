---
subcategory: "System Center Virtual Machine Manager"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_system_center_virtual_machine_manager_virtual_network"
description: |-
  Lists System Center Virtual Machine Manager Virtual Network resources.
---

# List resource: azurerm_system_center_virtual_machine_manager_virtual_network

Lists System Center Virtual Machine Manager Virtual Network resources.

## Example Usage

### List all System Center Virtual Machine Manager Virtual Networks in the subscription

```hcl
list "azurerm_system_center_virtual_machine_manager_virtual_network" "example" {
  provider = azurerm
  config {}
}
```

### List all System Center Virtual Machine Manager Virtual Networks in a specific resource group

```hcl
list "azurerm_system_center_virtual_machine_manager_virtual_network" "example" {
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