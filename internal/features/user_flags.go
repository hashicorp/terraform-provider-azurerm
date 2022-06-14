package features

type UserFeatures struct {
	ApiManagement          ApiManagementFeatures
	ApplicationInsights    ApplicationInsightFeatures
	CognitiveAccount       CognitiveAccountFeatures
	VirtualMachine         VirtualMachineFeatures
	VirtualMachineScaleSet VirtualMachineScaleSetFeatures
	KeyVault               KeyVaultFeatures
	TemplateDeployment     TemplateDeploymentFeatures
	LogAnalyticsWorkspace  LogAnalyticsWorkspaceFeatures
	ResourceGroup          ResourceGroupFeatures
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
	ScaleToZeroOnDelete       bool
}

type KeyVaultFeatures struct {
	PurgeSoftDeleteOnDestroy         bool
	PurgeSoftDeletedKeysOnDestroy    bool
	PurgeSoftDeletedCertsOnDestroy   bool
	PurgeSoftDeletedSecretsOnDestroy bool
	PurgeSoftDeletedHSMsOnDestroy    bool
	RecoverSoftDeletedKeyVaults      bool
	RecoverSoftDeletedKeys           bool
	RecoverSoftDeletedCerts          bool
	RecoverSoftDeletedSecrets        bool
}

type TemplateDeploymentFeatures struct {
	DeleteNestedItemsDuringDeletion bool
}

type LogAnalyticsWorkspaceFeatures struct {
	PermanentlyDeleteOnDestroy bool
}

type ResourceGroupFeatures struct {
	PreventDeletionIfContainsResources bool
}

type ApiManagementFeatures struct {
	PurgeSoftDeleteOnDestroy bool
	RecoverSoftDeleted       bool
}

type ApplicationInsightFeatures struct {
	DisableGeneratedRule bool
}
