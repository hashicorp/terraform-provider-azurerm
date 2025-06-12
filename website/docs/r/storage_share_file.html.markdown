---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_share_file"
description: |-
  Manages a File within an Azure Storage File Share.
---

# azurerm_storage_share_file

Manages a File within an Azure Storage File Share.

-> **Note:** When using Azure Active Directory Authentication (i.e. setting the provider property `storage_use_azuread = true`), the principal running Terraform must have the *Storage File Data Privileged Contributor* IAM role assigned. The *Storage File Data SMB Share Contributor* does not have sufficient permissions to create files. Refer to [official documentation](https://learn.microsoft.com/en-us/rest/api/storageservices/authorize-with-azure-active-directory#permissions-for-file-service-operations) for more details.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
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
  name             = "my-awesome-content.zip"
  storage_share_id = azurerm_storage_share.example.id
  source           = "some-local-file.zip"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name (or path) of the File that should be created within this File Share. Changing this forces a new resource to be created.

* `storage_share_id` - (Required) The Storage Share ID in which this file will be placed into. Changing this forces a new resource to be created.

* `path` - (Optional) The storage share directory that you would like the file placed into. Changing this forces a new resource to be created. Defaults to `""`.

* `source` - (Optional) An absolute path to a file on the local system. Changing this forces a new resource to be created.

~> **Note:** The file specified with `source` can not be empty.

* `content_type` - (Optional) The content type of the share file. Defaults to `application/octet-stream`.

* `content_md5` - (Optional) The MD5 sum of the file contents. Changing this forces a new resource to be created.

~> **Note:** This property is intended to be used with the Terraform internal [filemd5](https://www.terraform.io/docs/configuration/functions/filemd5.html) and [md5](https://www.terraform.io/docs/configuration/functions/md5.html) functions when `source` is defined.

* `content_encoding` - (Optional) Specifies which content encodings have been applied to the file.

* `content_disposition` - (Optional) Sets the file’s Content-Disposition header.

* `metadata` - (Optional) A mapping of metadata to assign to this file.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the file within the File Share.
* `content_length` - The length in bytes of the file content

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Share File.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Share File.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Share File.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Share File.

## Import

Directories within an Azure Storage File Share can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_share_file.example https://account1.file.core.windows.net/share1/file1
```
