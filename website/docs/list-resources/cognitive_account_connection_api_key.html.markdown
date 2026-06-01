---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_account_connection_api_key"
description: |-
  Lists Cognitive Services Account Connection with API Key authentication resources.
---

# List resource: azurerm_cognitive_account_connection_api_key

Lists Cognitive Services Account Connection with API Key authentication resources.

## Example Usage

```hcl
list "azurerm_cognitive_account_connection_api_key" "example" {
  provider = azurerm
  config {
    resource_group_name = "example-resources"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `resource_group_name` - (Required) The name of the Resource Group to query.

* `subscription_id` - (Optional) The Subscription ID to query. Defaults to the value specified in the Provider Configuration.