package features

type UserFeatures struct {
	CognitiveAccount       CognitiveAccountFeatures
	VirtualMachine         VirtualMachineFeatures
	VirtualMachineScaleSet VirtualMachineScaleSetFeatures
	KeyVault               KeyVaultFeatures
	Network                NetworkFeatures
	TemplateDeployment     TemplateDeploymentFeatures
	LogAnalyticsWorkspace  LogAnalyticsWorkspaceFeatures
}

type CognitiveAccountFeatures struct {
	PurgeSoftDeleteOnDestroy bool
}

type VirtualMachineFeatures struct {
	DeleteOSDiskOnDeletion     bool
	GracefulShutdown           bool
	SkipShutdownAndForceDelete bool
}

type VirtualMachineScaleSetFeatures struct {
	ForceDelete               bool
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

type LogAnalyticsWorkspaceFeatures struct {
	PermanentlyDeleteOnDestroy bool
}
