---
subcategory: "Storage Mover"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_mover_nfs_file_share_endpoint"
description: |-
  Manages a Storage Mover NFS File Share Endpoint.
---

# azurerm_storage_mover_nfs_file_share_endpoint

Manages a Storage Mover NFS File Share Endpoint.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestr"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Premium"
  account_replication_type = "LRS"
  account_kind             = "FileStorage"
}

resource "azurerm_storage_share" "example" {
  name                 = "example-share"
  storage_account_id   = azurerm_storage_account.example.id
  enabled_protocol     = "NFS"
  quota                = 100
}

resource "azurerm_storage_mover" "example" {
  name                = "example-ssm"
  resource_group_name = azurerm_resource_group.example.name
  location            = "West Europe"
}

resource "azurerm_storage_mover_nfs_file_share_endpoint" "example" {
  name               = "example-nfse"
  storage_mover_id   = azurerm_storage_mover.example.id
  storage_account_id = azurerm_storage_account.example.id
  file_share_name    = azurerm_storage_share.example.name
  description        = "Example NFS File Share Endpoint"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Storage Mover NFS File Share Endpoint. Changing this forces a new resource to be created.

* `storage_mover_id` - (Required) Specifies the ID of the Storage Mover for this NFS File Share Endpoint. Changing this forces a new resource to be created.

* `storage_account_id` - (Required) Specifies the ID of the Storage Account for this NFS File Share Endpoint. Changing this forces a new resource to be created.

* `file_share_name` - (Required) Specifies the name of the NFS File Share for this Storage Mover NFS File Share Endpoint. Changing this forces a new resource to be created.

* `description` - (Optional) Specifies a description for the Storage Mover NFS File Share Endpoint.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Storage Mover NFS File Share Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Mover NFS File Share Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Mover NFS File Share Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Mover NFS File Share Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Mover NFS File Share Endpoint.

## Import

Storage Mover NFS File Share Endpoint can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_mover_nfs_file_share_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.StorageMover/storageMovers/storageMover1/endpoints/endpoint1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.StorageMover` - 2025-07-01

