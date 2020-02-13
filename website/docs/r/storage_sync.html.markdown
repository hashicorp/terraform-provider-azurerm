---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_sync"
description: |-
  Manages an Azure Storage Sync.
---

# azurerm_storage_sync

Manages an Azure Storage Sync.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_sync" "test" {
  name                = "example-storage-sync"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  tags = {
    foo = "bar"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the storage sync. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the resource group the storage sync is located in. Changing this forces a new resource to be created.

* `location` - (Required) The Azure location where the storage sync exists. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags assigned to the storage sync.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the storage sync.

### Timeouts

~> **Note:** Custom Timeouts are available [as an opt-in Beta in version 1.43 of the Azure Provider](/docs/providers/azurerm/guides/2.0-beta.html) and will be enabled by default in version 2.0 of the Azure Provider.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Sync.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Sync.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Sync.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Sync.

## Import

Azure Storage Sync can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_sync.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.StorageSync/storageSyncServices/storagesync1
```
