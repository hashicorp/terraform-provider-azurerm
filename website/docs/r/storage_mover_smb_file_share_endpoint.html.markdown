---
subcategory: "Storage Mover"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_mover_smb_file_share_endpoint"
description: |-
  Manages a Storage Mover SMB File Share Target Endpoint.
---

# azurerm_storage_mover_smb_file_share_endpoint

Manages a Storage Mover **SMB File Share Target Endpoint** (an endpoint that targets an Azure Files SMB file share as the migration destination in the Azure portal).

-> **Note:** For a target endpoint that uses an Azure Storage **blob container**, use [`azurerm_storage_mover_target_endpoint`](storage_mover_target_endpoint.html) instead. For **NFS** file shares, use [`azurerm_storage_mover_nfs_file_share_target_endpoint`](storage_mover_nfs_file_share_target_endpoint.html).

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
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_share" "example" {
  name               = "example-share"
  storage_account_id = azurerm_storage_account.example.id
  quota              = 50
}

resource "azurerm_storage_mover" "example" {
  name                = "example-ssm"
  resource_group_name = azurerm_resource_group.example.name
  location            = "West Europe"
}

resource "azurerm_storage_mover_smb_file_share_endpoint" "example" {
  name               = "example-smbfse"
  storage_mover_id   = azurerm_storage_mover.example.id
  storage_account_id = azurerm_storage_account.example.id
  file_share_name    = azurerm_storage_share.example.name
  description        = "Example SMB File Share Target Endpoint"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of this Storage Mover SMB File Share Target Endpoint. Changing this forces a new resource to be created.

* `storage_mover_id` - (Required) Specifies the ID of the Storage Mover. Changing this forces a new resource to be created.

* `storage_account_id` - (Required) Specifies the ID of the Storage Account that hosts the SMB file share. Changing this forces a new resource to be created.

* `file_share_name` - (Required) The name of the Azure Files share (the same value as the `name` argument of the [`azurerm_storage_share`](storage_share.html) resource). Changing this forces a new resource to be created.

* `description` - (Optional) A description for this Target Endpoint.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Storage Mover SMB File Share Target Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Mover SMB File Share Target Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Mover SMB File Share Target Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Mover SMB File Share Target Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Mover SMB File Share Target Endpoint.

## Import

Storage Mover SMB File Share Target Endpoints can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_mover_smb_file_share_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.StorageMover/storageMovers/storageMover1/endpoints/endpoint1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.StorageMover` - 2025-07-01
