---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_queue"
sidebar_current: "docs-azurerm-resource-storage-queue"
description: |-
  Manages a Storage Queue within a Storage Account.
---

# azurerm_storage_queue

Manages a Storage Queue within a Storage Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  # ...
}

resource "azurerm_storage_account" "example" {
  # ...
}

resource "azurerm_storage_queue" "example" {
  name                 = "example-queue"
  resource_group_name  = "${azurerm_resource_group.example.name}"
  storage_account_name = "${azurerm_storage_account.example.name}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the storage queue. Changing this forces a new resource to be created.

-> **NOTE:** The `name` must be unique within the Storage Account where the Queue is located.

* `resource_group_name` - (Required) The name of the resource group in which to create the storage queue. Changing this forces a new resource to be created.

* `storage_account_name` - (Required) Specifies the storage account in which to create the storage queue. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the Storage Queue.
