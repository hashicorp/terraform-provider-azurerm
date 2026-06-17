---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_sync_server_endpoint"
description: |-
  Lists Storage Sync Server Endpoint resources.
---

# List resource: azurerm_storage_sync_server_endpoint

Lists Storage Sync Server Endpoint resources.

## Example Usage

### List all Storage Sync Server Endpoints in a Storage Sync Group

```hcl
list "azurerm_storage_sync_server_endpoint" "example" {
  provider = azurerm
  config {
    storage_sync_group_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.StorageSync/storageSyncServices/service1/syncGroups/syncGroup1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `storage_sync_group_id` - (Required) The ID of the Storage Sync Group to query.
