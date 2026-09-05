// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

import "github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"

func ResourceGuardVaultCriticalOperationExclusionList() func(interface{}, string) ([]string, []error) {
	return validation.StringInSlice([]string{
		"Microsoft.DataProtection/backupVaults/backupInstances/delete",
		"Microsoft.DataProtection/backupVaults/backupInstances/restore/action",
		"Microsoft.DataProtection/backupVaults/backupInstances/stopProtection/action",
		"Microsoft.DataProtection/backupVaults/backupInstances/suspendBackups/action",
		"Microsoft.DataProtection/backupVaults/backupInstances/write",
		"Microsoft.DataProtection/backupVaults/write#modifyEncryptionSettings",
		"Microsoft.DataProtection/backupVaults/write#reduceImmutabilityState",
		"Microsoft.RecoveryServices/vaults/backupconfig/delete",
		"Microsoft.RecoveryServices/vaults/backupEncryptionConfigs/backupResourceEncryptionConfig/write",
		"Microsoft.RecoveryServices/vaults/backupFabrics/protectionContainers/delete",
		"Microsoft.RecoveryServices/vaults/backupFabrics/protectionContainers/protectedItems/delete",
		"Microsoft.RecoveryServices/vaults/backupFabrics/protectionContainers/protectedItems/recoveryPoints/restore/action",
		"Microsoft.RecoveryServices/vaults/backupFabrics/protectionContainers/protectedItems/write",
		"Microsoft.RecoveryServices/vaults/backupPolicies/write",
		"Microsoft.RecoveryServices/vaults/backupResourceGuardProxies/write",
		"Microsoft.RecoveryServices/vaults/backupSecurityPIN/action",
		"Microsoft.RecoveryServices/vaults/write#modifyEncryptionSettings",
		"Microsoft.RecoveryServices/vaults/write#reduceImmutabilityState",
	}, false)
}
