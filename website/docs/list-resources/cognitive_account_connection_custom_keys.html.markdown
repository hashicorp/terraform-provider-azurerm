---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_account_connection_custom_keys"
description: |-
  Lists Cognitive Services Account Connection with Custom Keys authentication resources.
---

# List resource: azurerm_cognitive_account_connection_custom_keys

Lists Cognitive Services Account Connection with Custom Keys authentication resources.

## Example Usage

### List all Cognitive Services Account Connection with Custom Keys authentication resources in the subscription

```hcl
list "azurerm_cognitive_account_connection_custom_keys" "example" {
  provider = azurerm
  config {}
}
```

### List all Cognitive Services Account Connection with Custom Keys authentication resources in a specific Cognitive Account

```hcl
list "azurerm_cognitive_account_connection_custom_keys" "example" {
  provider = azurerm
  config {
    cognitive_account_name = "example-cognitive-account"
    resource_group_name    = "example-resources"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `cognitive_account_name` - (Optional) The name of the Cognitive Account to query. When specified, `resource_group_name` must also be specified.

* `resource_group_name` - (Optional) The name of the Resource Group to query.

* `subscription_id` - (Optional) The Subscription ID to query. Defaults to the value specified in the Provider Configuration.
