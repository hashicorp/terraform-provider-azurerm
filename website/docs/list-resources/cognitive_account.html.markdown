---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_account"
description: |-
  Lists Cognitive Account resources.
---
    
# List resource: azurerm_cognitive_account

Lists Cognitive Account resources. 

## Example Usage

### List all Cognitive Accounts in the subscription

```hcl
list "azurerm_cognitive_account" "example" {
  provider = azurerm
  config {}
}
```

### List all Cognitive Accounts in a specific resource group

```hcl
list "azurerm_cognitive_account" "example" {
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
