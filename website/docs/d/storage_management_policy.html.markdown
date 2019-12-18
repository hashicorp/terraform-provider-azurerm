---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_management_policy"
sidebar_current: "docs-azurerm-datasource-storage-management-policy"
description: |-
  Gets information about an existing Storage Management Policy.
---

# Data Source: azurerm_storage_management_policy

Use this data source to access information about an existing Storage Management Policy.

## Example Usage

```terraform
data "azurerm_storage_account" "example" {
  name                = "storageaccountname"
  resource_group_name = "resourcegroupname"
}

data "azurerm_storage_management_policy" "example" {
  storage_account_id = "${azurerm_storage_account.example.id}"
}
```

## Argument Reference

The following arguments are supported:

* `storage_account_id` - (Required) Specifies the id of the storage account to retrieve the management policy for.

## Attributes Reference

* `id` - The ID of the Management Policy.
* `rule` - A `rule` block as documented below.

---

* `rule` supports the following:

* `name` - (Required) A rule name can contain any combination of alpha numeric characters. Rule name is case-sensitive. It must be unique within a policy.
* `enabled` - (Required)  Boolean to specify whether the rule is enabled.
* `filters` - A `filter` block as documented below.
* `actions` - An `actions` block as documented below.

---

`filters` supports the following:

* `prefix_match` - An array of strings for prefixes to be matched.
* `blob_types` - An array of predefined values. Only `blockBlob` is supported.

---

`actions` supports the following:

* `base_blob` - A `base_blob` block as documented below.
* `snapshot` - A `snapshot` block as documented below.

---

`base_blob` supports the following:

* `tier_to_cool_after_days_since_modification_greater_than` - The age in days after last modification to tier blobs to cool storage. Supports blob currently at Hot tier.
* `tier_to_archive_after_days_since_modification_greater_than` - The age in days after last modification to tier blobs to archive storage. Supports blob currently at Hot or Cool tier.
* `delete_after_days_since_modification_greater_than` - The age in days after last modification to delete the blob.

---

`snapshot` supports the following:

* `delete_after_days_since_creation_greater_than` - The age in days after create to delete the snaphot.
