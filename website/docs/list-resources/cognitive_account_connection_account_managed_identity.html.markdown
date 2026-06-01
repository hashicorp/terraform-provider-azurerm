---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_account_connection_account_managed_identity"
description: |-
  Lists Cognitive Services Account Connection with Account Managed Identity authentication resources.
---

# List resource: azurerm_cognitive_account_connection_account_managed_identity

Lists Cognitive Services Account Connection with Account Managed Identity authentication resources.

## Example Usage

```hcl
list "azurerm_cognitive_account_connection_account_managed_identity" "example" {
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
