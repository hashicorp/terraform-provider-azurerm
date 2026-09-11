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
data "azurerm_storage_account" "example" {
  name                = "exampleaccount"
  resource_group_name = "examples"
}

data "azurerm_storage_queue" "example" {
  name               = "example-queue-name"
  storage_account_id = data.azurerm_storage_account.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Queue.

* `storage_account_id` - (Required) The ID of the Storage Account where the Queue exists.

## Attributes Reference

* `id` - The ID of this Storage Queue.

* `metadata` - A mapping of MetaData for this Queue.

* `url` - The data plane URL of the Storage Queue in the format of `<storage queue endpoint>/<queue name>`. E.g. `https://example.queue.core.windows.net/queue1`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Queue.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Storage` - 2025-08-01
