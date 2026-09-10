---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_action_group"
description: |-
    Lists Monitor Action Group resources.
---

# List resource: azurerm_monitor_action_group

Lists Monitor Action Group resources.

## Example Usage

### List all Monitor Action Groups in the subscription

```hcl
list "azurerm_monitor_action_group" "example" {
  provider = azurerm
  config {}
}
```

### List all Monitor Action Groups in a specific resource group

```hcl
list "azurerm_monitor_action_group" "example" {
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
