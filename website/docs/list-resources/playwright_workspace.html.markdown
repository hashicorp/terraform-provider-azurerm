---
subcategory: "Load Test"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_playwright_workspace"
description: |-
  Lists Playwright Workspace resources.
---

# List resource: azurerm_playwright_workspace

Lists Playwright Workspace resources.

## Example Usage

### List all Playwright Workspaces in the subscription

```hcl
list "azurerm_playwright_workspace" "example" {
  provider = azurerm
  config {}
}
```

### List all Playwright Workspaces in a specific resource group

```hcl
list "azurerm_playwright_workspace" "example" {
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
