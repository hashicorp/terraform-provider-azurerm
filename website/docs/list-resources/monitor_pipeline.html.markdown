---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_pipeline"
description: |-
  Lists Azure Monitor Pipeline resources.
---

# List resource: azurerm_monitor_pipeline

Lists Azure Monitor Pipeline resources.

## Example Usage

### List all Azure Monitor Pipeline resources in the subscription

```hcl
list "azurerm_monitor_pipeline" "example" {
  provider = azurerm
  config {}
}
```

### List all Azure Monitor Pipeline resources in a specific resource group

```hcl
list "azurerm_monitor_pipeline" "example" {
  provider = azurerm
  config {
    resource_group_name = "example-rg"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `resource_group_name` - (Optional) The name of the Resource Group to query.

* `subscription_id` - (Optional) The Subscription ID to query. Defaults to the value specified in the Provider Configuration.
