---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_account_capability_host"
description: |-
  Lists Cognitive Account Capability Host resources.
---

# List resource: azurerm_cognitive_account_capability_host

Lists Cognitive Account Capability Host resources.

## Example Usage

```hcl
data "azurerm_cognitive_account" "example" {
  name                = "example-account"
  resource_group_name = "example-resources"
}

list "azurerm_cognitive_account_capability_host" "example" {
  provider = azurerm
  config {
    cognitive_account_id = data.azurerm_cognitive_account.example.id
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `cognitive_account_id` - (Required) The ID of the Cognitive Account to query.
