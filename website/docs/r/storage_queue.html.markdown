---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_queue"
sidebar_current: "docs-azurerm-resource-storage-queue"
description: |-
  Manages a Queue within an Azure Storage Account.
---

# azurerm_storage_queue

Manages a Queue within an Azure Storage Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "test" {
  name                     = "examplestorageacc"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_queue" "test" {
  name                 = "mysamplequeue"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Queue which should be created within the Storage Account. Must be unique within the storage account the queue is located.

* `storage_account_name` - (Required) Specifies the Storage Account in which the Storage Queue should exist. Changing this forces a new resource to be created.

* `resource_group_name` - (Optional / **Deprecated**) The name of the resource group in which to create the storage queue.

* `metadata` - (Optional) A mapping of MetaData which should be assigned to this Storage Queue.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the Storage Queue.

## Import

Storage Queue's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_queue.queue1 https://example.queue.core.windows.net/queue1
```
