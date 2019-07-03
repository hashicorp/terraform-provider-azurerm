---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_account_blob_settings"
sidebar_current: "docs-azurerm-resource-azurerm-storage-account-blob-settings"
description: |-
  Manages Azure Storage Account Blob Properties.
---

# azurerm_storage_account_blob_settings

Manages Azure Storage Account Blob Properties.

## Example Usage

```hcl
resource "azurerm_resource_group" "testrg" {
  name     = "resourceGroupName"
  location = "westus"
}

resource "azurerm_storage_account" "testsa" {
  name                     = "storageaccountname"
  resource_group_name      = "${azurerm_resource_group.testrg.name}"
  location                 = "westus"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_account_blob_settings" "testsabs" {
  resource_group_name        = "${azurerm_resource_group.testrg.name}"
  storage_account_name       = "${azurerm_storage_account.testsa.name}"
  enable_soft_delete         = true
  soft_delete_retention_days = 123
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group in which to create the storage account. Changing this forces a new resource to be created.

* `storage_account_name` - (Required) Specifies the name of the storage account the blob settings are applied to. 

* `enable_soft_delete` - (Optional) Boolean flag which controls whether soft delete is enabled or not. Defaults to `false`.

* `soft_delete_retention_days` - (Optional) Specifies the number of days that the blob should be retained. The minimum specified value can be 1 an the maximum value can be 365. Defaults to `7`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the Storage Account Blob Settings.

## Import

Storage Accounts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_account_blob_settings.storageAccBlobSettings /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Storage/storageAccounts/myaccount/blobServices/default
```
