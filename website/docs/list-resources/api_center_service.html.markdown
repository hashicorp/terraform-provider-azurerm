---
subcategory: "ApiCenter"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_center_service"
description: |-
    Lists Api Center Service resources.
---

# List resource: azurerm_api_center_service

Lists Api Center Service resources.

## Example Usage

### List all API Center Services in a Resource Group

```hcl
list "azurerm_api_center_service" "example" {
  provider = azurerm
  config {
    resource_group_name = "example-rg"
  }
}
```

### List all API Center Services in the subscription

```hcl
list "azurerm_api_center_service" "example" {
  provider = azurerm
  config {
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `resource_group_name` - (Optional) The name of the Resource Group to query.

* `subscription_id` - (Optional) The ID of the Subscription to query. Defaults to the value specified in the Provider Configuration.
