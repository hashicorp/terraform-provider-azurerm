---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_blob_inventory_policy"
description: |-
  Manages a Storage Blob Inventory Policy.
---

# azurerm_storage_blob_inventory_policy

Manages a Storage Blob Inventory Policy.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestoracc"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  blob_properties {
    versioning_enabled = true
  }
}

resource "azurerm_storage_container" "example" {
  name                  = "examplecontainer"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

resource "azurerm_storage_blob_inventory_policy" "example" {
  storage_account_id = azurerm_storage_account.example.id
  rules {
    name                   = "rule1"
    storage_container_name = azurerm_storage_container.example.name
    format                 = "Csv"
    schedule               = "Daily"
    scope                  = "Container"
    schema_fields = [
      "Name",
      "Last-Modified",
    ]
  }
}

```

## Arguments Reference

The following arguments are supported:

* `storage_account_id` - (Required) The ID of the storage account to apply this Blob Inventory Policy to. Changing this forces a new Storage Blob Inventory Policy to be created.

* `rules` - (Required) One or more `rules` blocks as defined below.

---

A `filter` block supports the following:

* `blob_types` - (Required) A set of blob types. Possible values are `blockBlob`, `appendBlob`, and `pageBlob`. The storage account with `is_hns_enabled` is `true` doesn't support `pageBlob`.

~> **Note:** The `rules.*.schema_fields` for this rule has to include `BlobType` so that you can specify the `blob_types`.

* `include_blob_versions` - (Optional) Includes blob versions in blob inventory or not? Defaults to `false`.

~> **Note:** The `rules.*.schema_fields` for this rule has to include `IsCurrentVersion` and `VersionId` so that you can specify the `include_blob_versions`.

* `include_deleted` - (Optional) Includes deleted blobs in blob inventory or not? Defaults to `false`.

~> **Note:** If `rules.*.scope` is `Container`, the `rules.*.schema_fields` for this rule must include `Deleted`, `Version`, `DeletedTime`, and `RemainingRetentionDays` so that you can specify the `include_deleted`. If `rules.*.scope` is `Blob`, the `rules.*.schema_fields` must include `Deleted` and `RemainingRetentionDays` so that you can specify the `include_deleted`. If `rules.*.scope` is `Blob` and the storage account specified by `storage_account_id` has hierarchical namespaces enabled (`is_hns_enabled` is `true` on the storage account), the `rules.*.schema_fields` for this rule must include `Deleted`, `Version`, `DeletedTime`, and `RemainingRetentionDays` so that you can specify the `include_deleted`.

* `include_snapshots` - (Optional) Includes blob snapshots in blob inventory or not? Defaults to `false`.

~> **Note:** The `rules.*.schema_fields` for this rule has to include `Snapshot` so that you can specify the `include_snapshots`.

* `prefix_match` - (Optional) A set of strings for blob prefixes to be matched. Maximum of 10 blob prefixes.

* `exclude_prefixes` - (Optional) A set of strings for blob prefixes to be excluded. Maximum of 10 blob prefixes.

---

A `rules` block supports the following:

* `name` - (Required) The name which should be used for this Blob Inventory Policy Rule.

* `storage_container_name` - (Required) The storage container name to store the blob inventory files for this rule.

* `format` - (Required) The format of the inventory files. Possible values are `Csv` and `Parquet`.

* `schedule` - (Required) The inventory schedule applied by this rule. Possible values are `Daily` and `Weekly`.

* `scope` - (Required) The scope of the inventory for this rule. Possible values are `Blob` and `Container`.

* `schema_fields` - (Required) A list of fields to be included in the inventory. See the [Azure API reference](https://docs.microsoft.com/rest/api/storagerp/blob-inventory-policies/create-or-update#blobinventorypolicydefinition) for all the supported fields.

* `filter` - (Optional) A `filter` block as defined above.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Storage Blob Inventory Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Blob Inventory Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Blob Inventory Policy.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Blob Inventory Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Blob Inventory Policy.

## Import

Storage Blob Inventory Policies can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_blob_inventory_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Storage/storageAccounts/storageAccount1
```
