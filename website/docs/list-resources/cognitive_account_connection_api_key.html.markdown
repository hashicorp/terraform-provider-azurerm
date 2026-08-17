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
data "azurerm_cognitive_account" "example" {
  name                = "example-account"
  resource_group_name = "example-resources"
}

list "cognitive_account_connection_api_key" "example" {
  provider = azurerm
  config {
    cognitive_account_id = data.azurerm_cognitive_account.example.id
  }
}
```
## Argument Reference
This list resource supports the following arguments:
* `cognitive_account_id` - (Required) The ID of the Cognitive Account to query.
