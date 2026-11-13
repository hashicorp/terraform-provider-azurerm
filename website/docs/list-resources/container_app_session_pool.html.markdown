---
subcategory: "Container Apps"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_container_app_session_pool"
description: |-
  Lists Container App Session Pool resources.
---

# List resource: azurerm_container_app_session_pool

Lists Container App Session Pool resources.

## Example Usage

### List all Container App Session Pools in the subscription

```hcl
list "azurerm_container_app_session_pool" "example" {
  provider = azurerm
}
```

### List all Container App Session Pools in a specific Resource Group

```hcl
list "azurerm_container_app_session_pool" "example" {
  provider = azurerm
  config {
    resource_group_name = "example-resources"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `resource_group_name` - (Optional) The name of the Resource Group to query.

* `subscription_id` - (Optional) The Subscription ID to query. Defaults to the value specified in the Provider Configuration.
