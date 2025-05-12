---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_sync_server_endpoint"
description: |-
  Manages a Storage Sync Server Endpoint.
---

# azurerm_storage_sync_server_endpoint

Manages a Storage Sync Server Endpoint.

~> **Note:** The parent `azurerm_storage_sync_group` must have an `azurerm_storage_sync_cloud_endpoint` available before an `azurerm_storage_sync_server_endpoint` resource can be created.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_sync" "example" {
  name                = "example-storage-sync"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_storage_sync_group" "example" {
  name            = "example-storage-sync-group"
  storage_sync_id = azurerm_storage_sync.example.id
}

resource "azurerm_storage_account" "example" {
  name                     = "example-storage-account"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_share" "example" {
  name                 = "example-storage-share"
  storage_account_name = azurerm_storage_account.example.name
  quota                = 1

  acl {
    id = "GhostedRecall"
    access_policy {
      permissions = "r"
    }
  }
}

resource "azurerm_storage_sync_cloud_endpoint" "example" {
  name                  = "example-ss-ce"
  storage_sync_group_id = azurerm_storage_sync_group.example.id
  file_share_name       = azurerm_storage_share.example.name
  storage_account_id    = azurerm_storage_account.example.id
}

resource "azurerm_storage_sync_server_endpoint" "example" {
  name                  = "example-storage-sync-server-endpoint"
  storage_sync_group_id = azurerm_storage_sync_group.example.id
  registered_server_id  = azurerm_storage_sync.example.registered_servers[0]

  depends_on = [azurerm_storage_sync_cloud_endpoint.example]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Storage Sync. Changing this forces a new Storage Sync Server Endpoint to be created.

* `storage_sync_group_id` - (Required) The ID of the Storage Sync Group where the Storage Sync Server Endpoint should exist. Changing this forces a new Storage Sync Server Endpoint to be created.

* `registered_server_id` - (Required) The ID of the Registered Server that will be associate with the Storage Sync Server Endpoint. Changing this forces a new Storage Sync Server Endpoint to be created.

~> **Note:** The target server must already be registered with the parent `azurerm_storage_sync` prior to creating this endpoint. For more information on registering a server see the [Microsoft documentation](https://learn.microsoft.com/azure/storage/file-sync/file-sync-server-registration)

* `server_local_path` - (Required) The path on the Windows Server to be synced to the Azure file share. Changing this forces a new Storage Sync Server Endpoint to be created.

* `cloud_tiering_enabled` - (Optional)  Is Cloud Tiering Enabled? Defaults to `false`.

* `volume_free_space_percent` - (Optional) What percentage of free space on the volume should be preserved? Defaults to `20`.

* `tier_files_older_than_days` - (Optional) Files older than the specified age will be tiered to the cloud.

* `initial_download_policy` - (Optional)  Specifies how the server initially downloads the Azure file share data. Valid Values includes `NamespaceThenModifiedFiles`, `NamespaceOnly`, and `AvoidTieredFiles`. Defaults to `NamespaceThenModifiedFiles`.

* `local_cache_mode` - (Optional) Specifies how to handle the local cache. Valid Values include `UpdateLocallyCachedFiles` and `DownloadNewAndModifiedFiles`. Defaults to `UpdateLocallyCachedFiles`.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Storage Sync.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Sync Server Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Sync Server Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Sync Server Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Sync Server Endpoint.

## Import

Storage Sync Server Endpoints can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_sync_server_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.StorageSync/storageSyncServices/sync1/syncGroups/syncGroup1/serverEndpoints/endpoint1
```
