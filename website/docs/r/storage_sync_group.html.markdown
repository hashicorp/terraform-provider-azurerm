---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_sync_group"
description: |-
  Manages a Storage Sync Group.
---

# azurerm_storage_sync_group

Manages a Storage Sync Group.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_sync" "example" {
  name                = "example-ss"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_storage_sync_group" "example" {
  name            = "example-ss-group"
  storage_sync_id = azurerm_storage_sync.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Storage Sync Group. Changing this forces a new Storage Sync Group to be created.

* `storage_sync_id` - (Required) The resource ID of the Storage Sync where this Storage Sync Group is. Changing this forces a new Storage Sync Group to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Storage Sync Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Sync Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Sync Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Sync Group.

## Import

Storage Sync Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_sync_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.StorageSync/storageSyncServices/sync1/syncGroups/group1
```
