---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_blob"
sidebar_current: "docs-azurerm-resource-storage-blob"
description: |-
  Manages a Blob within a Storage Container.
---

# azurerm_storage_blob

Manages a Blob within a Storage Container.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-d"
  location = "westus"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestaccs"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  storage_account_name  = "${azurerm_storage_account.test.name}"
  container_access_type = "private"
}

resource "azurerm_storage_blob" "testsb" {
  name = "sample.vhd"

  resource_group_name    = "${azurerm_resource_group.test.name}"
  storage_account_name   = "${azurerm_storage_account.test.name}"
  storage_container_name = "${azurerm_storage_container.test.name}"

  type = "page"
  size = 5120
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the storage blob. Must be unique within the storage container the blob is located.

* `storage_account_name` - (Required) Specifies the storage account in which to create the storage container.
 Changing this forces a new resource to be created.

* `storage_container_name` - (Required) The name of the storage container in which this blob should be created.

* `type` - (Optional) The type of the storage blob to be created. One of either `block` or `page`. When not copying from an existing blob,
    this becomes required.

* `size` - (Optional) Used only for `page` blobs to specify the size in bytes of the blob to be created. Must be a multiple of 512. Defaults to 0.

* `content_type` - (Optional) The content type of the storage blob. Cannot be defined if `source_uri` is defined. Defaults to `application/octet-stream`.

* `source` - (Optional) An absolute path to a file on the local system. Cannot be defined if `source_uri` is defined.

* `source_uri` - (Optional) The URI of an existing blob, or a file in the Azure File service, to use as the source contents
    for the blob to be created. Changing this forces a new resource to be created. Cannot be defined if `source` is defined.

* `parallelism` - (Optional) The number of workers per CPU core to run for concurrent uploads. Defaults to `8`.

* `metadata` - (Optional) A map of custom blob metadata.

* `attempts` - (Optional / **Deprecated**) The number of attempts to make per page or block when uploading. Defaults to `1`.

* `resource_group_name` - (Optional / **Deprecated**) The name of the resource group in which to create the storage container.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the Storage Blob.
* `url` - The URL of the blob

## Import

Storage Blob's can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_blob.blob1 https://example.blob.core.windows.net/container/blob.vhd
```
