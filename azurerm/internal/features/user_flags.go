package features

type UserFeatures struct {
	VirtualMachine         VirtualMachineFeatures
	VirtualMachineScaleSet VirtualMachineScaleSetFeatures
	KeyVault               KeyVaultFeatures
	Network                NetworkFeatures
}

type VirtualMachineFeatures struct {
	DeleteOSDiskOnDeletion bool
}

type VirtualMachineScaleSetFeatures struct {
	RollInstancesWhenRequired bool
}

type KeyVaultFeatures struct {
	PurgeSoftDeleteOnDestroy    bool
	RecoverSoftDeletedKeyVaults bool
}

type NetworkFeatures struct {
	RelaxedLocking bool
}
