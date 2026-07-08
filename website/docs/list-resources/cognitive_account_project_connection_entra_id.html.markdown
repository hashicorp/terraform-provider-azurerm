---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_account_project_connection_entra_id"
description: |-
  Lists Cognitive Account Project Connection Entra ID resources.
---

# List resource: azurerm_cognitive_account_project_connection_entra_id

Lists Cognitive Account Project Connection Entra ID resources.

## Example Usage

```hcl
list "azurerm_cognitive_account_project_connection_entra_id" "example" {
  provider = azurerm
  config {
    cognitive_account_project_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.CognitiveServices/accounts/example-account/projects/example-project"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `cognitive_account_project_id` - (Required) The ID of the Cognitive Services Account Project to query.
