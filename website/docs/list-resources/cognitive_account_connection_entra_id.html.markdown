---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_account_connection_entra_id"
description: |-
  Lists Cognitive Services Account Connection with Entra ID authentication resources.
---

# List resource: azurerm_cognitive_account_connection_entra_id

Lists Cognitive Services Account Connection with Entra ID authentication resources.

## Example Usage
```hcl
data "azurerm_cognitive_account" "example" {
  name                = "example-account"
  resource_group_name = "example-resources"
}

list "cognitive_account_connection_entra_id" "example" {
  provider = azurerm
  config {
    cognitive_account_id = data.azurerm_cognitive_account.example.id
  }
}
```
## Argument Reference
This list resource supports the following arguments:
* `cognitive_account_id` - (Required) The ID of the Cognitive Account to query.
