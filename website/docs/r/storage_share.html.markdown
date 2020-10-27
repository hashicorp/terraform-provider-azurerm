---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_share"
description: |-
  Manages a File Share within Azure Storage.
---

# azurerm_storage_share

Manages a File Share within Azure Storage.

~> **Note:** The storage share supports two storage tiers: premium and standard. Standard file shares are created in general purpose (GPv1 or GPv2) storage accounts and premium file shares are created in FileStorage storage accounts. For further information, refer to the section "What storage tiers are supported in Azure Files?" of [documentation](https://docs.microsoft.com/en-us/azure/storage/files/storage-files-faq#general).

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

  acl {
    id = "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI"

    access_policy {
      permissions = "rwdl"
      start       = "2019-07-02T09:38:21.0000000Z"
      expiry      = "2019-07-02T10:38:21.0000000Z"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the share. Must be unique within the storage account where the share is located.

* `storage_account_name` - (Required) Specifies the storage account in which to create the share.
 Changing this forces a new resource to be created.

* `acl` - (Optional) One or more `acl` blocks as defined below.

* `quota` - (Optional) The maximum size of the share, in gigabytes. For Standard storage accounts, this must be greater than 0 and less than 5120 GB (5 TB). For Premium FileStorage storage accounts, this must be greater than 100 GB and less than 102400 GB (100 TB). Default is 5120.

* `metadata` - (Optional) A mapping of MetaData for this File Share.

---

A `acl` block supports the following:

* `id` - (Required) The ID which should be used for this Shared Identifier.

* `access_policy` - (Required) An `access_policy` block as defined below.

---

A `access_policy` block supports the following:

* `permissions` - (Required) The permissions which should be associated with this Shared Identifier. Possible value is combination of `r` (read), `w` (write), `d` (delete), and `l` (list).

~> **Note:** Permission order is strict at the service side, and permissions need to be listed in the order above. 

* `start` - (Optional) The time at which this Access Policy should be valid from, in [ISO8601](https://en.wikipedia.org/wiki/ISO_8601) format.

* `expiry` - (Optional) The time at which this Access Policy should be valid until, in [ISO8601](https://en.wikipedia.org/wiki/ISO_8601) format.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the File Share.

* `resource_manager_id` - The Resource Manager ID of this File Share.

* `url` - The URL of the File Share

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Share.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Share.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Share.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Share.

## Import

Storage Shares can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_share.exampleShare https://account1.file.core.windows.net/share1
```
