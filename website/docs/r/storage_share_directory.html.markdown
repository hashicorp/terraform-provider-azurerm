---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_share_directory"
description: |-
  Manages a Directory within an Azure Storage File Share.
---

# azurerm_storage_share_directory

Manages a Directory within an Azure Storage File Share.

-> **Note:** When using Azure Active Directory Authentication (i.e. setting the provider property `storage_use_azuread = true`), the principal running Terraform must have the *Storage File Data Privileged Contributor* IAM role assigned. The *Storage File Data SMB Share Contributor* does not have sufficient permissions to create directories. Refer to [official documentation](https://learn.microsoft.com/en-us/rest/api/storageservices/authorize-with-azure-active-directory#permissions-for-file-service-operations) for more details.

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

resource "azurerm_storage_share_directory" "example" {
  name             = "example"
  storage_share_id = azurerm_storage_share.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name (or path) of the Directory that should be created within this File Share. Changing this forces a new resource to be created.

* `storage_share_id` - (Required) The Storage Share ID in which this file will be placed into. Changing this forces a new resource to be created.

* `metadata` - (Optional) A mapping of metadata to assign to this Directory.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Directory within the File Share.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Share Directory.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Share Directory.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Share Directory.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Share Directory.

## Import

Directories within an Azure Storage File Share can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_share_directory.example https://tomdevsa20.file.core.windows.net/share1/directory1
```
