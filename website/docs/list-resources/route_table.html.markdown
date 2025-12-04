---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_route_table"
description: |-
  Lists Route Table resources.
---

# List resource: azurerm_route_table

~> **Note:** The `azurerm_route_table` List Resource is in beta. Its interface and behaviour may change as the feature evolves, and breaking changes are possible. It is offered as a technical preview without compatibility guarantees until Terraform 1.14 is generally available.

Lists Route Table resources.

## Example Usage

### List all Route Tables in the subscription

```hcl
list "azurerm_route_table" "example" {
  provider = azurerm
  config {}
}
```

### List all Route Tables in a specific resource group

```hcl
list "azurerm_route_table" "example" {
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
