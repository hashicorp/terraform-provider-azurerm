---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_share_rm"
description: |-
  Manages a File Share within Azure Storage.
---

# azurerm_storage_share_rm

Manages a File Share within Azure Storage.

~> **Note** The storage share supports two storage tiers: premium and standard. Standard file shares are created in general purpose (GPv1 or GPv2) storage accounts and premium file shares are created in FileStorage storage accounts. For further information, refer to the section "What storage tiers are supported in Azure Files?" of [documentation](https://docs.microsoft.com/azure/storage/files/storage-files-faq#general).

~> **Note on Authentication** Shared Keys will not be used to authenticate the requests.

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

resource "azurerm_storage_share_rm" "example" {
  name                 = "sharename"
  resource_group_name  = azurerm_resource_group.example.name
  storage_account_name = azurerm_storage_account.example.name
  quota                = 50

  acl {
    id = "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"

    access_policy {
      permission  = "rwdl"
      start_time  = "2019-07-02T09:38:21.0000000Z"
      expiry_time = "2019-07-02T10:38:21.0000000Z"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the share. Must be unique within the storage account where the share is located. Changing this forces a new resource to be created.

* `storage_account_name` - (Required) Specifies the storage account in which to create the share. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the resource group the Storage Account is located in.

* `access_tier` - (Optional) The access tier of the File Share. Possible values are `Hot`, `Cool` and `TransactionOptimized`, `Premium`.

~>**NOTE:** The `FileStorage` `account_kind` of the `azurerm_storage_account` requires `Premium` `access_tier`.

* `acl` - (Optional) One or more `acl` blocks as defined below.

* `enabled_protocol` - (Optional) The protocol used for the share. Possible values are `SMB` and `NFS`. The `SMB` indicates the share can be accessed by SMBv3.0, SMBv2.1 and REST. The `NFS` indicates the share can be accessed by NFSv4.1. Defaults to `SMB`. Changing this forces a new resource to be created.

~>**NOTE:** The `FileStorage` `account_kind` of the `azurerm_storage_account` is required for the `NFS` protocol.

* `quota` - (Required) The maximum size of the share, in gigabytes.

~>**NOTE:** For Standard storage accounts, by default this must be `1` GB (or higher) and at most `5120` GB (`5` TB). This can be set to a value larger than `5120` GB if `large_file_share_enabled` is set to `true` in the parent `azurerm_storage_account`.

~>**NOTE:** For Premium FileStorage storage accounts, this must be greater than `100` GB and at most `102400` GB (`100` TB).

* `root_squash` - (Optional) The root squash of the File Share. Possible values are `AllSquash`, `RootSquash` and `NoRootSquash` the setting is only working for `NFS`. Defaults to  `NoRootSquash`.

* `metadata` - (Optional) A mapping of MetaData for this File Share.

---

A `acl` block supports the following:

* `id` - (Required) The ID which should be used for this Shared Identifier.

* `access_policy` - (Optional) An `access_policy` block as defined below.

---

A `access_policy` block supports the following:

* `permission` - (Required) The permission which should be associated with this Shared Identifier. Possible value is combination of `r` (read), `w` (write), `d` (delete), and `l` (list).

~> **Note:** Permission order is strict at the service side, and permissions need to be listed in the order above.

* `start_time` - (Optional) The time at which this Access Policy should be valid from, in [ISO8601](https://en.wikipedia.org/wiki/ISO_8601) format.

* `expiry_time` - (Optional) The time at which this Access Policy should be valid until, in [ISO8601](https://en.wikipedia.org/wiki/ISO_8601) format.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the File Share.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Share.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Share.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Share.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Share.

## Import

Storage Shares can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_share_rm.exampleShare /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.Storage/storageAccounts/examplestorageaccountname/fileServices/default/shares/eyamplesharename

```
