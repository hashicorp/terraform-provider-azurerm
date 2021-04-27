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
  storage_account_id     = azurerm_storage_account.example.id
  storage_container_name = azurerm_storage_container.example.name
  rules {
    name = "rule1"
    filter {
      blob_types            = ["blockBlob"]
      include_blob_versions = true
      include_snapshots     = true
      prefix_match          = ["*/example"]
    }
  }
}

```

## Arguments Reference

The following arguments are supported:

* `storage_account_id` - (Required) The ID of the storage account to apply this Blob Inventory Policy to. Changing this forces a new Storage Blob Inventory Policy to be created.

* `storage_container_name` - (Required) The storage container name to store the blob inventory files. Changing this forces a new Storage Blob Inventory Policy to be created.

* `rules` - (Required) One or more `rules` blocks as defined below.

---

A `filter` block supports the following:

* `blob_types` - (Required)  A set of blob types. Possible values are `blockBlob`, `appendBlob`, and `pageBlob`. The storage account with `is_hns_enabled` is `true` doesn't support `pageBlob`.

* `include_blob_versions` - (Optional) Includes blob versions in blob inventory or not? Defaults to `false`.

* `include_snapshots` - (Optional) Includes blob snapshots in blob inventory or not? Defaults to `false`.

* `prefix_match` - (Optional) A set of strings for blob prefixes to be matched.

---

A `rules` block supports the following:

* `filter` - (Required) A `filter` block as defined above.

* `name` - (Required) The name which should be used for this Blob Inventory Policy Rule.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Storage Blob Inventory Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Blob Inventory Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Blob Inventory Policy.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Blob Inventory Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Blob Inventory Policy.

## Import

Storage Blob Inventory Policys can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_blob_inventory_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Storage/storageAccounts/storageAccount1/inventoryPolicies/inventoryPolicy1
```
