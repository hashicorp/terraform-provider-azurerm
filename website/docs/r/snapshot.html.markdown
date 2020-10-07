---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_snapshot"
description: |-
  Manages a Disk Snapshot.

---

# azurerm_snapshot

Manages a Disk Snapshot.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "snapshot-rg"
  location = "West Europe"
}

resource "azurerm_managed_disk" "example" {
  name                 = "managed-disk"
  location             = azurerm_resource_group.example.location
  resource_group_name  = azurerm_resource_group.example.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "10"
}

resource "azurerm_snapshot" "example" {
  name                = "snapshot"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  create_option       = "Copy"
  source_uri          = azurerm_managed_disk.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Snapshot resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Snapshot. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `create_option` - (Required) Indicates how the snapshot is to be created. Possible values are `Copy` or `Import`. Changing this forces a new resource to be created.

~> **Note:** One of `source_uri`, `source_resource_id` or `storage_account_id` must be specified.

* `source_uri` - (Optional) Specifies the URI to a Managed or Unmanaged Disk. Changing this forces a new resource to be created.

* `source_resource_id` - (Optional) Specifies a reference to an existing snapshot, when `create_option` is `Copy`. Changing this forces a new resource to be created.

* `storage_account_id` - (Optional) Specifies the ID of an storage account. Used with `source_uri` to allow authorization during import of unmanaged blobs from a different subscription. Changing this forces a new resource to be created.

* `disk_size_gb` - (Optional) The size of the Snapshotted Disk in GB.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The Snapshot ID.

* `disk_size_gb` - The Size of the Snapshotted Disk in GB.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Snapshot.
* `update` - (Defaults to 30 minutes) Used when updating the Snapshot.
* `read` - (Defaults to 5 minutes) Used when retrieving the Snapshot.
* `delete` - (Defaults to 30 minutes) Used when deleting the Snapshot.

## Import

Snapshots can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_snapshot.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/snapshots/snapshot1
```
