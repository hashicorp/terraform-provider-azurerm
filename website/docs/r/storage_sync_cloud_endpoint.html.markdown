---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_sync_cloud_endpoint"
description: |-
  Manages a Storage Sync Cloud Endpoint.
---

# azurerm_storage_sync_cloud_endpoint

Manages a Storage Sync Cloud Endpoint.

## Example Usage

```hcl
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

resource "azurerm_storage_account" "example" {
  name                     = "example-stracc"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_share" "example" {
  name                 = "example-share"
  storage_account_name = azurerm_storage_account.example.name
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
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Storage Sync Cloud Endpoint. Changing this forces a new Storage Sync Cloud Endpoint to be created.

* `storage_sync_group_id` - (Required) The ID of the Storage Sync Group where this Cloud Endpoint should be created. Changing this forces a new Storage Sync Cloud Endpoint to be created.

* `file_share_name` - (Required) The Storage Share name to be synchronized in this Storage Sync Cloud Endpoint. Changing this forces a new Storage Sync Cloud Endpoint to be created.

* `storage_account_id` - (Required) The ID of the Storage Account where the Storage Share exists. Changing this forces a new Storage Sync Cloud Endpoint to be created.



## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Storage Sync Cloud Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 45 minutes) Used when creating the Storage Sync Cloud Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Sync Cloud Endpoint.
* `delete` - (Defaults to 45 minutes) Used when deleting the Storage Sync Cloud Endpoint.

## Import

Storage Sync Cloud Endpoints can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_sync_cloud_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.StorageSync/storageSyncServices/sync1/syncGroups/syncgroup1/cloudEndpoints/cloudEndpoint1
```
