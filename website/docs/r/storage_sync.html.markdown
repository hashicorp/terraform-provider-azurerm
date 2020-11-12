---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_sync"
description: |-
  Manages a Storage Sync.
---

# azurerm_storage_sync

Manages a Storage Sync.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_sync" "test" {
  name                = "example-storage-sync"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  tags = {
    foo = "bar"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Storage Sync. Changing this forces a new Storage Sync to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Storage Sync should exist. Changing this forces a new Storage Sync to be created.

* `location` - (Required) The Azure Region where the Storage Sync should exist. Changing this forces a new Storage Sync to be created.

---

* `incoming_traffic_policy` - (Optional) Incoming traffic policy. Possible values are `AllowAllTraffic` and `AllowVirtualNetworksOnly`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Storage Sync.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Storage Sync.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Sync.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Sync.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Sync.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Sync.

## Import

Storage Syncs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_sync.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.StorageSync/storageSyncServices/sync1
```
