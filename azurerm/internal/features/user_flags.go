package features

type UserFeatures struct {
	Frontdoor              FrontdoorFeatures
	VirtualMachine         VirtualMachineFeatures
	VirtualMachineScaleSet VirtualMachineScaleSetFeatures
	KeyVault               KeyVaultFeatures
	Network                NetworkFeatures
	TemplateDeployment     TemplateDeploymentFeatures
}

type FrontdoorFeatures struct {
	IgnoreBackendPoolLimit bool
}

type VirtualMachineFeatures struct {
	DeleteOSDiskOnDeletion bool
	GracefulShutdown       bool
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

type TemplateDeploymentFeatures struct {
	DeleteNestedItemsDuringDeletion bool
}
