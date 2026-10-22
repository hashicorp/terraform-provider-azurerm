---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_workspace"
description: |-
    Lists Monitor Workspace resources.
---

# List resource: azurerm_monitor_workspace

Lists Monitor Workspace resources.

## Example Usage

### List all Monitor Workspaces in the subscription

```hcl
list "azurerm_monitor_workspace" "example" {
  provider = azurerm
  config {}
}
```

### List all Monitor Workspaces in a specific resource group

```hcl
list "azurerm_monitor_workspace" "example" {
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
