---
subcategory: "Storage Mover"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_mover"
description: |-
  Lists Storage Mover resources.
---

# List resource: azurerm_storage_mover

Lists Storage Mover resources.

## Example Usage

### List all Storage Movers in the subscription

```hcl
list "azurerm_storage_mover" "example" {
  provider = azurerm
  config {}
}
```

### List all Storage Movers in a specific resource group

```hcl
list "azurerm_storage_mover" "example" {
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