---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_blob"
description: |-
  Manages a Blob within a Storage Container.
---

# azurerm_storage_blob

Manages a Blob within a Storage Container.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestoracc"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name                  = "content"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

resource "azurerm_storage_blob" "example" {
  name                   = "my-awesome-content.zip"
  storage_account_name   = azurerm_storage_account.example.name
  storage_container_name = azurerm_storage_container.example.name
  type                   = "Block"
  source                 = "some-local-file.zip"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the storage blob. Must be unique within the storage container the blob is located.

* `storage_account_name` - (Required) Specifies the storage account in which to create the storage container.
 Changing this forces a new resource to be created.

* `storage_container_name` - (Required) The name of the storage container in which this blob should be created.

* `type` - (Required) The type of the storage blob to be created. Possible values are `Append`, `Block` or `Page`. Changing this forces a new resource to be created.

* `size` - (Optional) Used only for `page` blobs to specify the size in bytes of the blob to be created. Must be a multiple of 512. Defaults to 0.

~> **Note:** `size` is required if `source_uri` is not set.

* `access_tier` - (Optional) The access tier of the storage blob. Possible values are `Archive`, `Cool` and `Hot`.

* `content_type` - (Optional) The content type of the storage blob. Cannot be defined if `source_uri` is defined. Defaults to `application/octet-stream`.

* `content_md5` - (Optional) Used to trigger blob updates. This field cannot be specified if `source_uri` is specified or if blob type is Append or Page. The only meaningful value is ${filemd5("path/to/file")} if source is defined or "${md5("Sample blob content")}" if source_content is defined. Changing this forces a new resource to be created.

* `source` - (Optional) An absolute path to a file on the local system. This field cannot be specified for Append blobs and cannot be specified if `source_content` or `source_uri` is specified.

* `source_content` - (Optional) The content for this blob which should be defined inline. This field can only be specified for Block blobs and cannot be specified if `source` or `source_uri` is specified.

* `source_uri` - (Optional) The URI of an existing blob, or a file in the Azure File service, to use as the source contents
    for the blob to be created. Changing this forces a new resource to be created. This field cannot be specified for Append blobs and cannot be specified if `source` or `source_content` is specified.

* `parallelism` - (Optional) The number of workers per CPU core to run for concurrent uploads. Defaults to `8`.

~> **NOTE:** `parallelism` is only applicable for Page blobs - support for [Block Blobs is blocked on the upstream issue](https://github.com/tombuildsstuff/giovanni/issues/15).

* `metadata` - (Optional) A map of custom blob metadata.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the Storage Blob.
* `url` - The URL of the blob

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Blob.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Blob.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Blob.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Blob.

## Import

Storage Blob's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_blob.blob1 https://example.blob.core.windows.net/container/blob.vhd
```
