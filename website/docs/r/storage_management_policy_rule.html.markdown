---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_management_policy_rule"
description: |-
  Manages a Azure Storage Management Policy Rule.
---

# azurerm_storage_management_policy_rule

Manages a Azure Storage Management Policy Rule.

~> **NOTE on Storage Management Policy and Storage Management Policy Rules:** Terraform currently
provides both a standalone [Storage Management Policy Rule resource](storage_management_policy_rule.html), and allows for Storage Management Policy Rule to be defined in-line within the [Storage Management Policy](storage_management_policy.html).
At this time as there must be at least one `rule` set for creating a Storage Management Policy, if you want to manage Storage Management Policy Rules in standalone, you'll have to use `ignore_changes` lifecycle instructions on the `azurerm_storage_management_policy` for `rule` property, and be careful not to manage rules via `azurerm_storage_management_policy_rule` that are already managed by the `azurerm_storage_management_policy`. 

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "resourceGroupName"
  location = "West Europe"
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
    name = "rule1"
    filters {
      prefix_match = ["container1/prefix1"]
      blob_types   = ["blockBlob"]
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
  lifecycle {
    ignore_changes = [rule]
  }
}

resource "azurerm_storage_management_policy_rule" "example" {
  name                 = "rule2"
  management_policy_id = azurerm_storage_management_policy.example.id
  filter {
    prefix_match = ["container1/prefix1"]
    blob_types   = ["blockBlob"]
  }
  action {
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
```

## Arguments Reference

The following arguments are supported:

* `action` - (Required) A `action` block as defined below.

* `management_policy_id` - (Required) The ID of the Azure Storage Management Policy that this Azure Storage Management Policy Rule resides in. Changing this forces a new Azure Storage Management Policy Rule to be created.

* `name` - (Required) The name which should be used for this Azure Storage Management Policy Rule. Changing this forces a new Azure Storage Management Policy Rule to be created.

---

* `enabled` - (Optional) Should the Azure Storage Management Policy Rule be enabled? Defaults to `true`.

* `filter` - (Optional) A `filter` block as defined below.

---

A `action` block supports the following:

* `base_blob` - (Optional) A `base_blob` block as defined below.

* `snapshot` - (Optional) A `snapshot` block as defined below.

* `version` - (Optional) A `version` block as defined below.

---

A `base_blob` block supports the following:

* `delete_after_days_since_last_access_time_greater_than` - (Optional) The age in days after last access time to delete the blob. Must be between `0` and `99999`.

* `delete_after_days_since_modification_greater_than` - (Optional) The age in days after last modification to delete the blob. Must be between `0` and `99999`.

* `tier_to_archive_after_days_since_last_access_time_greater_than` - (Optional) The age in days after last access time to tier blobs to archive storage. Supports blob currently at Hot or Cool tier. Must be between `0` and `99999`.

* `tier_to_archive_after_days_since_last_tier_change_greater_than` - (Optional) The age in days after last tier change to the blobs to skip to be archved. Must be between `0` and `99999`.

* `tier_to_archive_after_days_since_modification_greater_than` - (Optional) The age in days after last modification to tier blobs to archive storage. Supports blob currently at Hot or Cool tier. Must be between `0` and `99999`.

* `tier_to_cool_after_days_since_modification_greater_than` - (Optional) The age in days after last modification to tier blobs to cool storage. Supports blob currently at Hot tier. Must be between `0` and `99999`.

* `tier_to_cool_after_days_since_last_access_time_greater_than` - (Optional) The age in days after last access time to tier blobs to cool storage. Supports blob currently at Hot tier. Must be between `0` and `99999`.

---

A `filter` block supports the following:

* `blob_types` - (Required) Specifies a list of blob types. Valid values are `blockBlob` and `appendBlob`.

* `match_blob_index_tag` - (Optional) One or more `match_blob_index_tag` blocks as defined below.

* `prefix_match` - (Optional) Specifies a list of prefixes to be matched.

~> **NOTE:** The `match_blob_index_tag` property requires enabling the `blobIndex` feature with [PSH or CLI commands](https://azure.microsoft.com/en-us/blog/manage-and-find-data-with-blob-index-for-azure-storage-now-in-preview/).

---

A `match_blob_index_tag` block supports the following:

* `name` - (Required) The name which should be used for this tag based filtering for blob objects.

* `value` - (Required) The filter tag value used for tag based filtering for blob objects.

* `operation` - (Optional) The comparison operator which is used for object comparison and filtering. Possible value is `==`. Defaults to `==`.

---

A `snapshot` block supports the following:

* `change_tier_to_archive_after_days_since_creation` - (Optional) The age in days after creation to tier blob snapshot to archive storage. Must be between `0` and `99999`.

* `change_tier_to_cool_after_days_since_creation` - (Optional) The age in days after creation to tier blob snapshot to cool storage. Must be between `0` and `99999`.

* `delete_after_days_since_creation_greater_than` - (Optional) The age in days after creation to delete the blob snapshot. Must be between `0` and `99999`.

* `tier_to_archive_after_days_since_last_tier_change_greater_than` - (Optional) The age in days after last tier change to the blobs to skip to be archved. Must be between `0` and `99999`.

---

A `version` block supports the following:

* `change_tier_to_archive_after_days_since_creation` - (Optional) The age in days after creation to tier blob version to archive storage. Must be between `0` and `99999`.

* `change_tier_to_cool_after_days_since_creation` - (Optional) The age in days creation create to tier blob version to cool storage. Must be between `0` and `99999`.

* `delete_after_days_since_creation` - (Optional) The age in days after creation to delete the blob version. Must be between `0` and `99999`.

* `tier_to_archive_after_days_since_last_tier_change_greater_than` - (Optional) The age in days after last tier change to the blobs to skip to be archved. Must be between `0` and `99999.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Azure Storage Management Policy Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure Storage Management Policy Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Storage Management Policy Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Azure Storage Management Policy Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Storage Management Policy Rule.

## Import

Azure Storage Management Policy Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_management_policy_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Storage/storageAccounts/myaccountname/managementPolicies/default/rules/rule1
```
