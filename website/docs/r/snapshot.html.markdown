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

* `create_option` - (Required) Indicates how the snapshot is to be created. Possible values are `Copy` or `Import`. 

~> **Note:** One of `source_uri`, `source_resource_id` or `storage_account_id` must be specified.

* `source_uri` - (Optional) Specifies the URI to a Managed or Unmanaged Disk. Changing this forces a new resource to be created.

* `source_resource_id` - (Optional) Specifies a reference to an existing snapshot, when `create_option` is `Copy`. Changing this forces a new resource to be created.

* `storage_account_id` - (Optional) Specifies the ID of an storage account. Used with `source_uri` to allow authorization during import of unmanaged blobs from a different subscription. Changing this forces a new resource to be created.

* `disk_size_gb` - (Optional) The size of the Snapshotted Disk in GB.

* `encryption_settings` - (Optional) A `encryption_settings` block as defined below.

~> **Note:** Removing `encryption_settings` forces a new resource to be created.

* `incremental_enabled` - (Optional) Specifies if the Snapshot is incremental. Changing this forces a new resource to be created.

* `network_access_policy` - (Optional) Policy for accessing the disk via network. Possible values are `AllowAll`, `AllowPrivate`, or `DenyAll`. Defaults to `AllowAll`.

* `disk_access_id` - (Optional) Specifies the ID of the Disk Access which should be used for this Snapshot. This is used in conjunction with setting `network_access_policy` to `AllowPrivate`.

* `public_network_access_enabled` - (Optional) Policy for controlling export on the disk. Possible values are `true` or `false`. Defaults to `true`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `encryption_settings` block supports:

* `disk_encryption_key` - (Optional) A `disk_encryption_key` block as defined below.

* `key_encryption_key` - (Optional) A `key_encryption_key` block as defined below.

---

The `disk_encryption_key` block supports:

* `secret_url` - (Required) The URL to the Key Vault Secret used as the Disk Encryption Key. This can be found as `id` on the `azurerm_key_vault_secret` resource.

* `source_vault_id` - (Required) The ID of the source Key Vault. This can be found as `id` on the `azurerm_key_vault` resource.

---

The `key_encryption_key` block supports:

* `key_url` - (Required) The URL to the Key Vault Key used as the Key Encryption Key. This can be found as `id` on the `azurerm_key_vault_key` resource.

* `source_vault_id` - (Required) The ID of the source Key Vault. This can be found as `id` on the `azurerm_key_vault` resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Snapshot ID.

* `disk_size_gb` - The Size of the Snapshotted Disk in GB.

* `trusted_launch_enabled` - Whether Trusted Launch is enabled for the Snapshot.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Snapshot.
* `read` - (Defaults to 5 minutes) Used when retrieving the Snapshot.
* `update` - (Defaults to 30 minutes) Used when updating the Snapshot.
* `delete` - (Defaults to 30 minutes) Used when deleting the Snapshot.

## Import

Snapshots can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_snapshot.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/snapshots/snapshot1
```
