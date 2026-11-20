---
subcategory: "Machine Learning"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_machine_learning_registry"
description: |-
  Lists Machine Learning Registry resources.
---

# List resource: azurerm_machine_learning_registry

Lists Machine Learning Registry resources.

## Example Usage

### List all Machine Learning Registries in the subscription

```hcl
list "azurerm_machine_learning_registry" "example" {
  provider = azurerm
  config {}
}
```

### List all Machine Learning Registries in a specific resource group

```hcl
list "azurerm_machine_learning_registry" "example" {
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
