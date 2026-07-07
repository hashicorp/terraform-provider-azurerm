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
```hcl
data "azurerm_cognitive_account" "example" {
  name                = "example-account"
  resource_group_name = "example-resources"
}

list "cognitive_account_connection_custom_keys" "example" {
  provider = azurerm
  config {
    cognitive_account_id = data.azurerm_cognitive_account.example.id
  }
}
```
## Argument Reference
This list resource supports the following arguments:
* `cognitive_account_id` - (Required) The ID of the Cognitive Account to query.
