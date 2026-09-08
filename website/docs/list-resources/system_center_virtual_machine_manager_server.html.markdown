---
subcategory: "System Center Virtual Machine Manager"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_system_center_virtual_machine_manager_server"
description: |-
  Lists System Center Virtual Machine Manager Server resources.
---

# List resource: azurerm_system_center_virtual_machine_manager_server

Lists System Center Virtual Machine Manager Server resources.

## Example Usage

### List all System Center Virtual Machine Manager Servers in the subscription

```hcl
list "azurerm_system_center_virtual_machine_manager_server" "example" {
  provider = azurerm
  config {}
}
```

### List all System Center Virtual Machine Manager Servers in a specific resource group

```hcl
list "azurerm_system_center_virtual_machine_manager_server" "example" {
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
