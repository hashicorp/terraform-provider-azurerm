package dataprotection

const (
	GuardOperationStopProtectionWithRetainData = "Microsoft.RecoveryServices/vaults/backupFabrics/protectionContainers/protectedItems/write#stopProtectionWithRetainData"
	GuardOperationDeleteProtectedItem          = "Microsoft.RecoveryServices/vaults/backupFabrics/protectionContainers/protectedItems/delete"
	GuardOperationUpdateProtectedItem          = "Microsoft.RecoveryServices/vaults/backupFabrics/protectionContainers/protectedItems/write"
	GuardOperationUpdateBackupPolicy           = "Microsoft.RecoveryServices/vaults/backupPolicies/write"
	GuardOperationDeleteResourceGuardProxy     = "Microsoft.RecoveryServices/vaults/backupResourceGuardProxies/delete"
	GuardOperationGetSecurityPIN               = "Microsoft.RecoveryServices/vaults/backupSecurityPIN/action"
	GuardOperationDisableSoftDelete            = "Microsoft.RecoveryServices/vaults/backupconfig/write"
	GuardOperationReduceImmutabilityState      = "Microsoft.RecoveryServices/vaults/write#reduceImmutabilityState"
	GuardOperationTriggerResotre               = "Microsoft.RecoveryServices/vaults/backupFabrics/protectionContainers/protectedItems/recoveryPoints/restore/action"
)
