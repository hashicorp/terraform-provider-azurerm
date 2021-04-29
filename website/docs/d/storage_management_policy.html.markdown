---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_management_policy"
description: |-
  Gets information about an existing Storage Management Policy.
---

# Data Source: azurerm_storage_management_policy

Use this data source to access information about an existing Storage Management Policy.

## Example Usage

```hcl
data "azurerm_storage_account" "example" {
  name                = "storageaccountname"
  resource_group_name = "resourcegroupname"
}

data "azurerm_storage_management_policy" "example" {
  storage_account_id = azurerm_storage_account.example.id
}
```

## Argument Reference

The following arguments are supported:

* `storage_account_id` - Specifies the id of the storage account to retrieve the management policy for.

## Attributes Reference

* `id` - The ID of the Management Policy.
* `rule` - A `rule` block as documented below.

---

* `rule` supports the following:

* `name` - A rule name can contain any combination of alpha numeric characters. Rule name is case-sensitive. It must be unique within a policy.
* `enabled` -  Boolean to specify whether the rule is enabled.
* `filters` - A `filter` block as documented below.
* `actions` - An `actions` block as documented below.

---

`filters` supports the following:

* `prefix_match` - An array of strings for prefixes to be matched.
* `blob_types` - An array of predefined values. Valid options are `blockBlob` and `appendBlob`.
* `match_blob_index_tag` - A `match_blob_index_tag` block as defined below. The block defines the blob index tag based filtering for blob objects.
---

`actions` supports the following:

* `base_blob` - A `base_blob` block as documented below.
* `snapshot` - A `snapshot` block as documented below.
* `version` - A `version` block as documented below.

---

`base_blob` supports the following:

* `tier_to_cool_after_days_since_modification_greater_than` - The age in days after last modification to tier blobs to cool storage. Supports blob currently at Hot tier.
* `tier_to_archive_after_days_since_modification_greater_than` - The age in days after last modification to tier blobs to archive storage. Supports blob currently at Hot or Cool tier.
* `delete_after_days_since_modification_greater_than` - The age in days after last modification to delete the blob.

---

`snapshot` supports the following:

* `change_tier_to_archive_after_days_since_creation` - The age in days after creation to tier blob snapshot to archive storage.
* `change_tier_to_cool_after_days_since_creation` - The age in days after creation to tier blob snapshot to cool storage.
* `delete_after_days_since_creation_greater_than` - The age in days after creation to delete the blob snapshot.

---

`version` supports the following:

* `change_tier_to_archive_after_days_since_creation` - The age in days after creation to tier blob version to archive storage.
* `change_tier_to_cool_after_days_since_creation` - The age in days after creation to tier blob version to cool storage.
* `delete_after_days_since_creation` - The age in days after creation to delete the blob version.

---

`match_blob_index_tag` supports the following:

* `name` - The filter tag name used for tag based filtering for blob objects.
* `operation` - The comparison operator which is used for object comparison and filtering. Possible value is `==`. Defaults to `==`.
* `value` -  The filter tag value used for tag based filtering for blob objects.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Management Policy.
