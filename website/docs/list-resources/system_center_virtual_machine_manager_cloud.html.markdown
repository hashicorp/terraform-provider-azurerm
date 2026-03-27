---
subcategory: "System Center Virtual Machine Manager"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_system_center_virtual_machine_manager_cloud"
description: |-
  Lists System Center Virtual Machine Manager Cloud resources.
---

# List resource: azurerm_system_center_virtual_machine_manager_cloud

Lists System Center Virtual Machine Manager Cloud resources.

## Example Usage

### List all System Center Virtual Machine Manager Clouds in the subscription

```hcl
list "azurerm_system_center_virtual_machine_manager_cloud" "example" {
  provider = azurerm
  config {}
}
```

### List all System Center Virtual Machine Manager Clouds in a specific resource group

```hcl
list "azurerm_system_center_virtual_machine_manager_cloud" "example" {
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
