---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_service_plan"
description: |-
  Lists Service Plan resources.
---

# List resource: azurerm_service_plan

Lists Service Plan resources.

## Example Usage

### List all Service Plans in the subscription

```hcl
list "azurerm_service_plan" "example" {
  provider = azurerm
  config {}
}
```

### List all Service Plans in a specific resource group

```hcl
list "azurerm_service_plan" "example" {
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
