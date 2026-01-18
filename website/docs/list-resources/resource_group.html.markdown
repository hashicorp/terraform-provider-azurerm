---
subcategory: "Base"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_resource_group"
description: |-
  Lists Resource Group resources.
---

# List resource: azurerm_resource_group

~> **Note:** The `azurerm_resource_group` List Resource is in beta. Its interface and behaviour may change as the feature evolves, and breaking changes are possible. It is offered as a technical preview without compatibility guarantees until Terraform 1.14 is generally available.

Lists Resource Group resources.

## Example Usage

### List all Resource Groups in the subscription

```hcl
list "azurerm_resource_group" "example" {
  provider = azurerm
  config {}
}
```

### List all Resource Groups in the subscription matching a filter 

```hcl
list "azurerm_resource_group" "example" {
  provider = azurerm
  config {
    filter = "tagName eq 'terraform' and tagValue eq 'example'"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `subscription_id` - (Optional) The ID of the subscription to query. Defaults to the value specified in the Provider Configuration.

* `filter` - (Optional) A filter expression to filter the results by.
