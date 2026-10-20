---
subcategory: "DataProtection"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_protection_resource_guard"
description: |-
  Manages a Resource Guard.
---

# azurerm_data_protection_resource_guard

Manages a Resource Guard.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_data_protection_resource_guard" "example" {
  name                = "example-resourceguard"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Resource Guard. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Resource Guard should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Resource Guard should exist. Changing this forces a new resource to be created.

* `vault_critical_operation_exclusion_list` - (Optional) A list of the critical operations which are not protected by this Resource Guard.

-> **Note:** Azure Backup documentation describes these operations by friendly names. See [Recovery Services vault critical operations](https://learn.microsoft.com/azure/backup/multi-user-authorization-concept?tabs=recovery-services-vault#critical-operations) and [Backup vault critical operations](https://learn.microsoft.com/azure/backup/multi-user-authorization-concept?tabs=backup-vault#critical-operations). The API expects the following literal values for `vault_critical_operation_exclusion_list`:

  | Vault type | Azure documentation operation | Valid API value |
  |---|---|---|
  | Recovery Services vault | `Delete protection` | `Microsoft.RecoveryServices/vaults/backupconfig/delete` |
  | Recovery Services vault | `Delete protection` | `Microsoft.RecoveryServices/vaults/backupFabrics/protectionContainers/protectedItems/delete` |
  | Recovery Services vault | `Modify protection` | `Microsoft.RecoveryServices/vaults/backupFabrics/protectionContainers/protectedItems/write` |
  | Recovery Services vault | `Modify protection` | `Microsoft.RecoveryServices/vaults/backupResourceGuardProxies/write` |
  | Recovery Services vault | `Modify policy` | `Microsoft.RecoveryServices/vaults/backupPolicies/write` |
  | Recovery Services vault | `Get backup security PIN` | `Microsoft.RecoveryServices/vaults/backupSecurityPIN/action` |
  | Recovery Services vault | `Modify encryption settings` | `Microsoft.RecoveryServices/vaults/backupEncryptionConfigs/backupResourceEncryptionConfig/write` |
  | Recovery Services vault | `Modify encryption settings` | `Microsoft.RecoveryServices/vaults/write#modifyEncryptionSettings` |
  | Recovery Services vault | `Disable immutability` | `Microsoft.RecoveryServices/vaults/write#reduceImmutabilityState` |
  | Recovery Services vault | `Restore` | `Microsoft.RecoveryServices/vaults/backupFabrics/protectionContainers/protectedItems/recoveryPoints/restore/action` |
  | Recovery Services vault | `Delete hybrid container` | `Microsoft.RecoveryServices/vaults/backupFabrics/protectionContainers/delete` |
  | Backup vault | `Delete Backup Instance` | `Microsoft.DataProtection/backupVaults/backupInstances/delete` |
  | Backup vault | `Disable immutability` | `Microsoft.DataProtection/backupVaults/write#reduceImmutabilityState` |
  | Backup vault | `Modify encryption settings` | `Microsoft.DataProtection/backupVaults/write#modifyEncryptionSettings` |
  | Backup vault | `Stop backup and retain forever` | `Microsoft.DataProtection/backupVaults/backupInstances/stopProtection/action` |
  | Backup vault | `Stop backup and retain as per policy` | `Microsoft.DataProtection/backupVaults/backupInstances/suspendBackups/action` |
  | Backup vault | `Change policy` | `Microsoft.DataProtection/backupVaults/backupInstances/write` |
  | Backup vault | `Restore` | `Microsoft.DataProtection/backupVaults/backupInstances/restore/action` |

* `tags` - (Optional) A mapping of tags which should be assigned to the Resource Guard.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Resource Guard.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Resource Guard.
* `read` - (Defaults to 5 minutes) Used when retrieving the Resource Guard.
* `update` - (Defaults to 30 minutes) Used when updating the Resource Guard.
* `delete` - (Defaults to 30 minutes) Used when deleting the Resource Guard.

## Import

Resource Guards can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_protection_resource_guard.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DataProtection/resourceGuards/resourceGuard1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.DataProtection` - 2025-07-01
