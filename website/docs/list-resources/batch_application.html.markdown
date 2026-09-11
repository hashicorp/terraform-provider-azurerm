---
subcategory: "Batch"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_batch_application"
description: |-
    Lists Batch Application resources.
---

# List resource: azurerm_batch_application

Lists Batch Application resources.

## Example Usage

### List all Batch Applications in a Batch Account

```hcl
list "azurerm_batch_application" "example" {
  provider = azurerm
  config {
    batch_account_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Batch/batchAccounts/account1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `batch_account_id` - (Required) The ID of the Batch Account to query.
