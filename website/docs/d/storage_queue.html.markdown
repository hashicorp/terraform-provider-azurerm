---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_queue"
description: |-
  Gets information about an existing Storage Queue.
---

# Data Source: azurerm_storage_queue

Use this data source to access information about an existing Storage Queue.

## Example Usage

```hcl
data "azurerm_storage_queue" "example" {
  name                 = "example-queue-name"
  storage_account_name = "example-storage-account-name"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Queue.

* `storage_account_name` - (Required) The name of the Storage Account where the Queue exists.

## Attributes Reference

* `metadata` - A mapping of MetaData for this Queue.

* `resource_manager_id` - The Resource Manager ID of this Storage Queue.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Queue.
