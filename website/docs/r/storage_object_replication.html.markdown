---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_object_replication"
description: |-
  Copy Block Blobs between a source storage account and a destination account
---

# azurerm_storage_object_replication

Manages a Storage Object Replication.

## Example Usage

```hcl
resource "azurerm_resource_group" "src" {
  name     = "srcResourceGroupName"
  location = "West Europe"
}

resource "azurerm_storage_account" "src" {
  name                     = "srcstorageaccount"
  resource_group_name      = azurerm_resource_group.src.name
  location                 = azurerm_resource_group.src.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  blob_properties {
    versioning_enabled  = true
    change_feed_enabled = true
  }
}

resource "azurerm_storage_container" "src" {
  name                  = "srcstrcontainer"
  storage_account_name  = azurerm_storage_account.src.name
  container_access_type = "private"
}

resource "azurerm_resource_group" "dst" {
  name     = "dstResourceGroupName"
  location = "East US"
}

resource "azurerm_storage_account" "dst" {
  name                     = "dststorageaccount"
  resource_group_name      = azurerm_resource_group.dst.name
  location                 = azurerm_resource_group.dst.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  blob_properties {
    versioning_enabled  = true
    change_feed_enabled = true
  }
}

resource "azurerm_storage_container" "dst" {
  name                  = "dststrcontainer"
  storage_account_name  = azurerm_storage_account.dst.name
  container_access_type = "private"
}

resource "azurerm_storage_object_replication" "example" {
  source_storage_account_id      = azurerm_storage_account.src.id
  destination_storage_account_id = azurerm_storage_account.dst.id
  rules {
    source_container_name      = azurerm_storage_container.src.name
    destination_container_name = azurerm_storage_container.dst.name
  }
}
```

## Arguments Reference

The following arguments are supported:

* `source_storage_account_id` - (Required) The ID of the source storage account. Changing this forces a new Storage Object Replication to be created.

* `destination_storage_account_id` - (Required) The ID of the destination storage account. Changing this forces a new Storage Object Replication to be created.

* `rules` - (Required) One or more `rules` blocks as defined below.

---

A `rules` block supports the following:

* `source_container_name` - (Required) The source storage container name. Changing this forces a new Storage Object Replication to be created.

* `destination_container_name` - (Required) The destination storage container name. Changing this forces a new Storage Object Replication to be created.

* `copy_blobs_created_after` - (Optional) The time after which the Block Blobs created will be copies to the destination. Possible values are `OnlyNewObjects`, `Everything` and time in RFC3339 format: `2006-01-02T15:04:00Z`.

* `filter_out_blobs_with_prefix` - (Optional) Specifies a list of filters prefixes, the blobs whose names begin with which will be replicated.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Storage Object Replication in the destination storage account. It's composed as format `source_object_replication_id;destination_object_replication_id`. 

* `source_object_replication_id` - The ID of the Object Replication in the source storage account.

* `destination_object_replication_id` - The ID of the Object Replication in the destination storage account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Object Replication.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Object Replication.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Object Replication.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Object Replication.

## Import

Storage Object Replication Policys can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_object_replication.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Storage/storageAccounts/storageAccount1/objectReplicationPolicies/objectReplicationPolicy1;/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group2/providers/Microsoft.Storage/storageAccounts/storageAccount2/objectReplicationPolicies/objectReplicationPolicy2
```
