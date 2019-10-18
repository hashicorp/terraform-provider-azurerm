---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_share_directory"
sidebar_current: "docs-azurerm-resource-storage-share-directory"
description: |-
  Manages a Directory within an Azure Storage File Share.
---

# azurerm_storage_share_directory

Manages a Directory within an Azure Storage File Share.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "azuretest"
  location = "West Europe"
}

resource "azurerm_storage_account" "test" {
  name                     = "azureteststorage"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_share" "test" {
  name                 = "sharename"
  storage_account_name = "${azurerm_storage_account.test.name}"
  quota                = 50
}

resource "azurerm_storage_share_directory" "test" {
  name                 = "example"
  share_name           = "${azurerm_storage_share.test.name}"
  storage_account_name = "${azurerm_storage_account.test.name}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name (or path) of the Directory that should be created within this File Share. Changing this forces a new resource to be created.

* `share_name` - (Required) The name of the File Share where this Directory should be created. Changing this forces a new resource to be created.

* `storage_account_name` - (Required) The name of the Storage Account within which the File Share is located. Changing this forces a new resource to be created.

* `metadata` - (Optional) A mapping of metadata to assign to this Directory.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the Directory within the File Share.

## Import

Directories within an Azure Storage File Share can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_share_directory.test https://tomdevsa20.file.core.windows.net/share1/directory1
```
