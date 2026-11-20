---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_actions_task_definition"
description: |-
  Lists Storage Actions Task Definition resources.
---

# List resource: azurerm_storage_actions_task_definition

Lists Storage Actions Task Definition resources.

## Example Usage

### List all Storage Actions Task Definitions in the subscription

```hcl
list "azurerm_storage_actions_task_definition" "example" {
  provider = azurerm
  config {}
}
```

### List all Storage Actions Task Definitions in a specific resource group

```hcl
list "azurerm_storage_actions_task_definition" "example" {
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
