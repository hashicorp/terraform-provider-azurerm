---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_protection_backup_vault"
description: |-
  Manages a Backup Vault.
---

# azurerm_data_protection_backup_vault

Manages a Backup Vault.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_data_protection_backup_vault" "example" {
  name                = "example-backup-vault"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  datastore_type      = "VaultStore"
  redundancy          = "LocallyRedundant"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Backup Vault. Changing this forces a new Backup Vault to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Backup Vault should exist. Changing this forces a new Backup Vault to be created.

* `location` - (Required) The Azure Region where the Backup Vault should exist. Changing this forces a new Backup Vault to be created.

* `datastore_type` - (Required) Specifies the type of the data store. Possible values are `ArchiveStore`, `OperationalStore`, `SnapshotStore` and `VaultStore`. Changing this forces a new resource to be created.

-> **Note:** The `SnapshotStore` will be removed in version 4.0 as it has been replaced by `OperationalStore`.

* `redundancy` - (Required) Specifies the backup storage redundancy. Possible values are `GeoRedundant`, `LocallyRedundant` and `ZoneRedundant`. Changing this forces a new Backup Vault to be created.

* `cross_region_restore_enabled` - (Optional) Whether to enable cross-region restore for the Backup Vault.
 
-> **Note:** The `cross_region_restore_enabled` can only be specified when `redundancy` is specified for `GeoRedundant`. Once `cross_region_restore_enabled` is enabled, it cannot be disabled.

---

* `identity` - (Optional) An `identity` block as defined below.

* `retention_duration_in_days` - (Optional) The soft delete retention duration for this Backup Vault. Possible values are between `14` and `180`. Defaults to `14`.

-> **Note:** The `retention_duration_in_days` is the number of days for which deleted data is retained before being permanently deleted. Retention period till 14 days are free of cost, however, retention beyond 14 days may incur additional charges. The `retention_duration_in_days` is required when the `soft_delete` is set to `On`.

* `immutability` - (Optional) The state of immutability for this Backup Vault. Possible values are `Disabled`, `Locked`, and `Unlocked`. Defaults to `Disabled`. Changing this from `Locked` to anything else forces a new Backup Vault to be created.

* `soft_delete` - (Optional) The state of soft delete for this Backup Vault. Possible values are `AlwaysOn`, `Off`, and `On`. Defaults to `On`.

-> **Note:** Once the `soft_delete` is set to `AlwaysOn`, the setting cannot be changed.

* `tags` - (Optional) A mapping of tags which should be assigned to the Backup Vault.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Backup Vault. The only possible value is `SystemAssigned`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Backup Vault.

* `identity` - An `identity` block as defined below, which contains the Identity information for this Backup Vault.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Identity of this Backup Vault.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Identity of this Backup Vault.

-> **Note:** You can access the Principal ID via `${azurerm_data_protection_backup_vault.example.identity[0].principal_id}` and the Tenant ID via `${azurerm_data_protection_backup_vault.example.identity[0].tenant_id}`

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Backup Vault.
* `read` - (Defaults to 5 minutes) Used when retrieving the Backup Vault.
* `update` - (Defaults to 30 minutes) Used when updating the Backup Vault.
* `delete` - (Defaults to 30 minutes) Used when deleting the Backup Vault.

## Import

Backup Vaults can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_protection_backup_vault.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DataProtection/backupVaults/vault1
```
