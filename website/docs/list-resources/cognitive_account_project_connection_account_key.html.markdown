---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_account_project_connection_account_key"
description: |-
  Lists Cognitive Services Project Connection with Account Key authentication resources.
---

# List resource: azurerm_cognitive_account_project_connection_account_key

Lists Cognitive Services Project Connection with Account Key authentication resources.

## Example Usage

```hcl
data "azurerm_cognitive_account_project" "example" {
  name                   = "existing-cognitive-account-project"
  cognitive_account_name = "existing-cognitive-account"
  resource_group_name    = "existing-resource-group"
}

list "azurerm_cognitive_account_project_connection_account_key" "example" {
  provider = azurerm
  config {
    cognitive_account_project_id = data.azurerm_cognitive_account_project.example.id
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `cognitive_account_project_id` - (Required) The ID of the Cognitive Services Project to query.
