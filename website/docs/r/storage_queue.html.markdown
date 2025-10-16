---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_queue"
description: |-
  Manages a Queue within an Azure Storage Account.
---

# azurerm_storage_queue

Manages a Queue within an Azure Storage Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestorageacc"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_queue" "example" {
  name                 = "mysamplequeue"
  storage_account_name = azurerm_storage_account.example.name
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Queue which should be created within the Storage Account. Must be unique within the storage account the queue is located. Changing this forces a new resource to be created.

* `storage_account_name` - (Optional) The name of the Storage Account where the Storage Queue should be created. Changing this forces a new resource to be created. This property is deprecated in favour of `storage_account_id`.

~> **Note:** Migrating from the deprecated `storage_account_name` to `storage_account_id` is supported without recreation. Any other change to either property will result in the resource being recreated.

* `storage_account_id` - (Optional) The name of the Storage Account where the Storage Queue should be created. Changing this forces a new resource to be created.

~> **Note:** One of `storage_account_name` or `storage_account_id` must be specified. When specifying `storage_account_id` the resource will use the Resource Manager API, rather than the Data Plane API.

* `metadata` - (Optional) A mapping of MetaData which should be assigned to this Storage Queue.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Storage Queue.

* `resource_manager_id` - The Resource Manager ID of this Storage Queue.

* `url` - The data plane URL of the Storage Queue in the format of `<storage queue endpoint>/<queue name>`. E.g. `https://example.queue.core.windows.net/queue1`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Queue.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Queue.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Queue.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Queue.

## Import

Storage Queue's can be imported using the `resource id`, e.g.

If `storage_account_name` is used:

```shell
terraform import azurerm_storage_queue.queue1 https://example.queue.core.windows.net/queue1
```

If `storage_account_id` is used:

```shell
terraform import azurerm_storage_queue.queue1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Storage/storageAccounts/myaccount/queueServices/default/queues/queue1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Storage` - 2023-05-01
