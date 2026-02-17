---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_nat_gateway"
description: |-
    Lists Nat Gateway resources.
---

# List resource: azurerm_nat_gateway

Lists Nat Gateway resources.

## Example Usage

### List all Nat Gateways

```hcl
list "azurerm_nat_gateway" "example" {
  provider = azurerm
  config {
  }
}
```

### List all Nat Gateways in a Resource Group

```hcl
list "azurerm_nat_gateway" "example" {
  provider = azurerm
  config {
    resource_group_name = "resource_group_name-example"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `subscription_id` - (Optional) The ID of the Subscription to query.

* `resource_group_name` - (Optional) The name of the Resource Group to query.
