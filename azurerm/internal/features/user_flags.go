package features

type UserFeatures struct {
	VirtualMachine         VirtualMachineFeatures
	VirtualMachineScaleSet VirtualMachineScaleSetFeatures
	KeyVault               KeyVaultFeatures
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
