---
subcategory: "Automation"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_automation_account"
description: |-
  Lists Automation Account resources.
---

# List resource: azurerm_automation_account

Lists Automation Account resources.

## Example Usage

### List all Automation Accounts in the subscription

```hcl
list "azurerm_automation_account" "example" {
  provider = azurerm
  config {}
}
```

### List all Automation Accounts in a specific resource group

```hcl
list "azurerm_automation_account" "example" {
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
