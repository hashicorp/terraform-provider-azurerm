---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_management_policy"
description: |-
  Manages an Azure Storage Account Management Policy.
---

# azurerm_storage_management_policy

Manages an Azure Storage Account Management Policy.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "resourceGroupName"
  location = "westus"
}

resource "azurerm_storage_account" "example" {
  name                = "storageaccountname"
  resource_group_name = azurerm_resource_group.example.name

  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "BlobStorage"
}

resource "azurerm_storage_management_policy" "example" {
  storage_account_id = azurerm_storage_account.example.id

  rule {
    name    = "rule1"
    enabled = true
    filters {
      prefix_match = ["container1/prefix1"]
      blob_types   = ["blockBlob"]
      blob_index_match_tag {
        name      = "tag1"
        operation = "=="
        value     = "val1"
      }
    }
    actions {
      base_blob {
        tier_to_cool_after_days_since_modification_greater_than    = 10
        tier_to_archive_after_days_since_modification_greater_than = 50
        delete_after_days_since_modification_greater_than          = 100
      }
      snapshot {
        delete_after_days_since_creation_greater_than = 30
      }
    }
  }
  rule {
    name    = "rule2"
    enabled = false
    filters {
      prefix_match = ["container2/prefix1", "container2/prefix2"]
      blob_types   = ["blockBlob"]
    }
    actions {
      base_blob {
        tier_to_cool_after_days_since_modification_greater_than    = 11
        tier_to_archive_after_days_since_modification_greater_than = 51
        delete_after_days_since_modification_greater_than          = 101
      }
      snapshot {
        delete_after_days_since_creation_greater_than = 31
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `storage_account_id` - (Required) Specifies the id of the storage account to apply the management policy to.

* `rule` - (Optional) A `rule` block as documented below.

---

* `rule` supports the following:

* `name` - (Required) A rule name can contain any combination of alpha numeric characters. Rule name is case-sensitive. It must be unique within a policy.
* `enabled` - (Required)  Boolean to specify whether the rule is enabled.
* `filters` - A `filter` block as documented below.
* `actions` - An `actions` block as documented below.

---

`filters` supports the following:

* `prefix_match` - An array of strings for prefixes to be matched.
* `blob_types` - An array of predefined values. Valid options are `blockBlob` and `appendBlob`.
* `blob_index_match_tag` - A `blob_index_match_tag` block as defined below. The block defines the blob index tag based filtering for blob objects.
~> **NOTE:** This property requires enabling the `blobIndex` feature with [PSH or Cli commands](https://azure.microsoft.com/en-us/blog/manage-and-find-data-with-blob-index-for-azure-storage-now-in-preview/) before setting the block `blob_index_match_tag`. 
---

`actions` supports the following:

* `base_blob` - A `base_blob` block as documented below.
* `snapshot` - A `snapshot` block as documented below.

---

`base_blob` supports the following:

* `tier_to_cool_after_days_since_modification_greater_than` - The age in days after last modification to tier blobs to cool storage. Supports blob currently at Hot tier. Must be at least 0.
* `tier_to_archive_after_days_since_modification_greater_than` - The age in days after last modification to tier blobs to archive storage. Supports blob currently at Hot or Cool tier. Must be at least 0.
* `delete_after_days_since_modification_greater_than` - The age in days after last modification to delete the blob. Must be at least 0.

---

`snapshot` supports the following:

* `delete_after_days_since_creation_greater_than` - The age in days after create to delete the snaphot. Must be at least 0.

---

`blob_index_match_tag` supports the following:

* `name` - The filter tag name used for tag based filtering for blob objects.
* `operation` - The comparison operator which is used for object comparison and filtering. Possible value is `==`. Defaults to `==`.
* `value` -  The filter tag value used for tag based filtering for blob objects.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Storage Account Management Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Account Management Policy.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Account Management Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Account Management Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Account Management Policy.

## Import

Storage Account Management Policies can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_management_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Storage/storageAccounts/myaccountname/managementPolicies/default
```
