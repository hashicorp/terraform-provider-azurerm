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

### List all Container App Session Pools in a specific resource group

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

* `resource_group_name` - (Optional) The name of the resource group to query.

* `subscription_id` - (Optional) The Subscription ID to query. Defaults to the value specified in the Provider Configuration.

~> **Note:** The `value` of each `secret` block is not returned by the Azure API, so listed Container App Session Pools do not include secret values.
