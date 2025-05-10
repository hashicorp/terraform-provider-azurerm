---
subcategory: "Recovery Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_recovery_services_vault"
description: |-
  Manages a Recovery Services Vault.
---

# azurerm_recovery_services_vault

Manages a Recovery Services Vault.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "tfex-recovery_vault"
  location = "West Europe"
}

resource "azurerm_recovery_services_vault" "vault" {
  name                = "example-recovery-vault"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Standard"

  soft_delete_enabled = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Recovery Services Vault. Recovery Service Vault name must be 2 - 50 characters long, start with a letter, contain only letters, numbers and hyphens. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Recovery Services Vault. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `identity` - (Optional) An `identity` block as defined below.

* `sku` - (Required) Sets the vault's SKU. Possible values include: `Standard`, `RS0`.

* `public_network_access_enabled` - (Optional) Is it enabled to access the vault from public networks. Defaults to `true`.

* `immutability` - (Optional) Immutability Settings of vault, possible values include: `Locked`, `Unlocked` and `Disabled`.

-> **Note:** Once `immutability` is set to `Locked`, changing it to other values forces a new Recovery Services Vault to be created.

* `storage_mode_type` - (Optional) The storage type of the Recovery Services Vault. Possible values are `GeoRedundant`, `LocallyRedundant` and `ZoneRedundant`. Defaults to `GeoRedundant`.

* `cross_region_restore_enabled` - (Optional) Is cross region restore enabled for this Vault? Only can be `true`, when `storage_mode_type` is `GeoRedundant`. Defaults to `false`.

-> **Note:** Once `cross_region_restore_enabled` is set to `true`, changing it back to `false` forces a new Recovery Service Vault to be created.

* `soft_delete_enabled` - (Optional) Is soft delete enable for this Vault? Defaults to `true`.

* `encryption` - (Optional) An `encryption` block as defined below. Required with `identity`.

!> **Note:** Once Encryption with your own key has been Enabled it's not possible to Disable it.

* `classic_vmware_replication_enabled` - (Optional) Whether to enable the Classic experience for VMware replication. If set to `false` VMware machines will be protected using the new stateless ASR replication appliance. Changing this forces a new resource to be created.

* `monitoring` - (Optional) A `monitoring` block as defined below.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Recovery Services Vault. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) A list of User Assigned Managed Identity IDs to be assigned to this App Configuration.

~> **Note:** `identity_ids` is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

An `encryption` block supports the following:

* `key_id` - (Required) The Key Vault key id used to encrypt this vault. Key managed by Vault Managed Hardware Security Module is also supported.

* `infrastructure_encryption_enabled` - (Required) Enabling/Disabling the Double Encryption state.

* `user_assigned_identity_id` - (Optional) Specifies the user assigned identity ID to be used.

* `use_system_assigned_identity` - (Optional) Indicate that system assigned identity should be used or not. Defaults to `true`. Must be set to `false` when `user_assigned_identity_id` is set.

!> **Note:** `use_system_assigned_identity` only be able to set to `false` for **new** vaults. Any vaults containing existing items registered or attempted to be registered to it are not supported. Details can be found in [the document](https://learn.microsoft.com/en-us/azure/backup/encryption-at-rest-with-cmk?tabs=portal#before-you-start)

!> **Note:** Once `infrastructure_encryption_enabled` has been set it's not possible to change it.

---

A `monitoring` block supports the following:

* `alerts_for_all_job_failures_enabled` - (Optional) Enabling/Disabling built-in Azure Monitor alerts for security scenarios and job failure scenarios. Defaults to `true`.

* `alerts_for_critical_operation_failures_enabled` - (Optional) Enabling/Disabling alerts from the older (classic alerts) solution. Defaults to `true`. More details could be found [here](https://learn.microsoft.com/en-us/azure/backup/monitoring-and-alerts-overview).

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Recovery Services Vault.

* `identity` - An `identity` block as defined below.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 hours) Used when creating the Recovery Services Vault.
* `read` - (Defaults to 5 minutes) Used when retrieving the Recovery Services Vault.
* `update` - (Defaults to 1 hour) Used when updating the Recovery Services Vault.
* `delete` - (Defaults to 30 minutes) Used when deleting the Recovery Services Vault.

## Import

Recovery Services Vaults can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_recovery_services_vault.vault1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.RecoveryServices/vaults/vault1
```
