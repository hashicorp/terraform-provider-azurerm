---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_endpoint"
description: |-
    Lists Private Endpoint resources.
---

# List resource: azurerm_private_endpoint

Lists Private Endpoint resources.

## Example Usage

### List all Private Endpoints in the subscription

```hcl
list "azurerm_private_endpoint" "example" {
  provider = azurerm
  config {
  }
}
```

### List all Private Endpoints in a Resource Group

```hcl
list "azurerm_private_endpoint" "example" {
  provider = azurerm
  config {
    resource_group_name = "example-rg"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `subscription_id` - (Optional) The ID of the Subscription to query. Defaults to the value specified in the Provider Configuration.

* `resource_group_name` - (Optional) The name of the Resource Group to query.
