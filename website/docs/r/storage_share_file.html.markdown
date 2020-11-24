---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_share_file"
description: |-
  Manages a File within an Azure Storage File Share.
---

# azurerm_storage_share_file

Manages a File within an Azure Storage File Share.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "azuretest"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "azureteststorage"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_share" "example" {
  name                 = "sharename"
  storage_account_name = azurerm_storage_account.example.name
  quota                = 50
}

resource "azurerm_storage_share_file" "example" {
  name                   = "my-awesome-content.zip"
  share_name             = azurerm_storage_share.example.name
  storage_account_name   = azurerm_storage_account.example.name
  source                 = "some-local-file.zip"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name (or path) of the File that should be created within this File Share. Changing this forces a new resource to be created.

* `share_name` - (Required) The name of the File Share where this File should be created. Changing this forces a new resource to be created.

* `storage_account_name` - (Required) The name of the Storage Account within which the File Share is located. Changing this forces a new resource to be created.

* `directory_name` - (Optional) The storage share directory that you would like the file placed into. Changing this forces a new resource to be created.

* `source` - (Optional) An absolute path to a file on the local system.

* `content_type` - (Optional) The content type of the share file. Defaults to `application/octet-stream`.

* `content_md5` - (Optional) The MD5 sum of the file contents. Changing this forces a new resource to be created.   

* `content_encoding` - (Optional) Specifies which content encodings have been applied to the file.

* `content_disposition` - (Optional) Sets the file’s Content-Disposition header.

* `parallelism` - (Optional)  The number of workers per CPU core to run for concurrent uploads. Defaults to `4`.

* `metadata` - (Optional) A mapping of metadata to assign to this file.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the file within the File Share.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Share File.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Share File.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Share File.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Share File.

## Import

Directories within an Azure Storage File Share can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_share_file.example https://account1.file.core.windows.net/share1/file1
```
