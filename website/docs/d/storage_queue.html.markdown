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

* `storage_account_name` - (Optional) The name of the Storage Account where the Queue exists. This property is deprecated in favour of `storage_account_id`.

* `storage_account_id` - (Optional) The name of the Storage Account where the Queue exists. This property will become Required in version 5.0 of the Provider.

~> **Note:** One of `storage_account_name` or `storage_account_id` must be specified. When specifying `storage_account_id` the resource will use the Resource Manager API, rather than the Data Plane API.

## Attributes Reference

* `metadata` - A mapping of MetaData for this Queue.

* `resource_manager_id` - The Resource Manager ID of this Storage Queue.

* `url` - The data plane URL of the Storage Queue in the format of `<storage queue endpoint>/<queue name>`. E.g. `https://example.queue.core.windows.net/queue1`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Queue.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Storage` - 2023-05-01
