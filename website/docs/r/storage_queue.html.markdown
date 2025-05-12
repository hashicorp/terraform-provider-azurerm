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

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Queue which should be created within the Storage Account. Must be unique within the storage account the queue is located. Changing this forces a new resource to be created.

* `storage_account_name` - (Required) Specifies the Storage Account in which the Storage Queue should exist. Changing this forces a new resource to be created.

* `metadata` - (Optional) A mapping of MetaData which should be assigned to this Storage Queue.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Storage Queue.

* `resource_manager_id` - The Resource Manager ID of this Storage Queue.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Queue.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Queue.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Queue.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Queue.

## Import

Storage Queue's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_queue.queue1 https://example.queue.core.windows.net/queue1
```
