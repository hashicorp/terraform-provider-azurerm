---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_share"
description: |-
  Manages a File Share within Azure Storage.
---

# azurerm_storage_share

Manages a File Share within Azure Storage.

~> **Note:** The storage share supports two storage tiers: premium and standard. Standard file shares are created in general purpose (GPv1 or GPv2) storage accounts and premium file shares are created in FileStorage storage accounts. For further information, refer to the section "What storage tiers are supported in Azure Files?" of [documentation](https://docs.microsoft.com/azure/storage/files/storage-files-faq#general).

~> **Note:** Shared Key authentication will always be used for this resource, as AzureAD authentication is not supported by the Storage API for files.

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
  name               = "sharename"
  storage_account_id = azurerm_storage_account.example.id
  quota              = 50

  acl {
    id = "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"

    access_policy {
      permissions = "rwdl"
      start       = "2019-07-02T09:38:21Z"
      expiry      = "2019-07-02T10:38:21Z"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the share. Must be unique within the storage account where the share is located. Changing this forces a new resource to be created.

* `storage_account_name` - (Optional) Specifies the storage account in which to create the share. Changing this forces a new resource to be created. This property is deprecated in favour of `storage_account_id`.

~> **Note:** Migrating from the deprecated `storage_account_name` to `storage_account_id` is supported without recreation. Any other change to either property will result in the resource being recreated.

* `storage_account_id` - (Optional) Specifies the storage account in which to create the share. Changing this forces a new resource to be created.

~> **Note:** One of `storage_account_name` or `storage_account_id` must be specified. When specifying `storage_account_id` the resource will use the Resource Manager API, rather than the Data Plane API.

* `access_tier` - (Optional) The access tier of the File Share. Possible values are `Hot`, `Cool` and `TransactionOptimized`, `Premium`.

~> **Note:** The `FileStorage` `account_kind` of the `azurerm_storage_account` requires `Premium` `access_tier`.

* `acl` - (Optional) One or more `acl` blocks as defined below.

* `enabled_protocol` - (Optional) The protocol used for the share. Possible values are `SMB` and `NFS`. The `SMB` indicates the share can be accessed by SMBv3.0, SMBv2.1 and REST. The `NFS` indicates the share can be accessed by NFSv4.1. Defaults to `SMB`. Changing this forces a new resource to be created.

~> **Note:** The `FileStorage` `account_kind` of the `azurerm_storage_account` is required for the `NFS` protocol.

* `quota` - (Required) The maximum size of the share, in gigabytes.

~> **Note:** For Standard storage accounts, by default this must be `1` GB (or higher) and at most `5120` GB (`5` TB). This can be set to a value larger than `5120` GB if `large_file_share_enabled` is set to `true` in the parent `azurerm_storage_account`.

~> **Note:** For Premium FileStorage storage accounts, this must be greater than `100` GB and at most `102400` GB (`100` TB).

* `metadata` - (Optional) A mapping of MetaData for this File Share.

---

A `acl` block supports the following:

* `id` - (Required) The ID which should be used for this Shared Identifier.

* `access_policy` - (Optional) An `access_policy` block as defined below.

---

A `access_policy` block supports the following:

* `permissions` - (Required) The permissions which should be associated with this Shared Identifier. Possible value is combination of `r` (read), `w` (write), `d` (delete), and `l` (list).

~> **Note:** Permission order is strict at the service side, and permissions need to be listed in the order above.

* `start` - (Optional) The time at which this Access Policy should be valid from. When using `storage_account_id` this should be in RFC3339 format. If using the deprecated `storage_account_name` property, this uses the [ISO8601](https://en.wikipedia.org/wiki/ISO_8601) format.

* `expiry` - (Optional) The time at which this Access Policy should be valid untilWhen using `storage_account_id` this should be in RFC3339 format. If using the deprecated `storage_account_name` property, this uses the [ISO8601](https://en.wikipedia.org/wiki/ISO_8601) format.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the File Share.

* `resource_manager_id` - The Resource Manager ID of this File Share.

* `url` - The URL of the File Share

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Share.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Share.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Share.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Share.

## Import

Storage Shares can be imported using the `id`, e.g.

```shell
terraform import azurerm_storage_share.exampleShare /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Storage/storageAccounts/myAccount/fileServices/default/shares/exampleShare
```
