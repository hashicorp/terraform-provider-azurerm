---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_traffic_manager_profile"
description: |-
  Lists Traffic Manager Profile resources.
---

# List resource: azurerm_traffic_manager_profile

Lists Traffic Manager Profile resources.

## Example Usage

### List all Traffic Manager Profiles in the subscription

```hcl
list "azurerm_traffic_manager_profile" "example" {
  provider = azurerm
  config {}
}
```

### List all Traffic Manager Profiles in a specific resource group

```hcl
list "azurerm_traffic_manager_profile" "example" {
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