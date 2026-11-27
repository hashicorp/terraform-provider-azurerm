---
subcategory: "Discovery"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_discovery_bookshelf"
description: |-
  Lists Discovery Bookshelf resources.
---

# List resource: azurerm_discovery_bookshelf

Lists Discovery Bookshelf resources.

## Example Usage

### List all Discovery Bookshelves in the subscription

```hcl
list "azurerm_discovery_bookshelf" "example" {
  provider = azurerm
  config {}
}
```

### List all Discovery Bookshelves in a specific resource group

```hcl
list "azurerm_discovery_bookshelf" "example" {
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
