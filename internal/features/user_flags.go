// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package features

type UserFeatures struct {
	ApiManagement            ApiManagementFeatures
	AppConfiguration         AppConfigurationFeatures
	ApplicationInsights      ApplicationInsightFeatures
	CognitiveAccount         CognitiveAccountFeatures
	VirtualMachine           VirtualMachineFeatures
	VirtualMachineScaleSet   VirtualMachineScaleSetFeatures
	KeyVault                 KeyVaultFeatures
	TemplateDeployment       TemplateDeploymentFeatures
	LogAnalyticsWorkspace    LogAnalyticsWorkspaceFeatures
	ResourceGroup            ResourceGroupFeatures
	RecoveryServicesVault    RecoveryServicesVault
	ManagedDisk              ManagedDiskFeatures
	Storage                  StorageFeatures
	Subscription             SubscriptionFeatures
	PostgresqlFlexibleServer PostgresqlFlexibleServerFeatures
	MachineLearning          MachineLearningFeatures
	RecoveryService          RecoveryServiceFeatures
	NetApp                   NetAppFeatures
	DatabricksWorkspace      DatabricksWorkspaceFeatures
}

type CognitiveAccountFeatures struct {
	PurgeSoftDeleteOnDestroy bool
}

type VirtualMachineFeatures struct {
	DetachImplicitDataDiskOnDeletion bool
	DeleteOSDiskOnDeletion           bool
	GracefulShutdown                 bool // TODO: Remove in 5.0 - Currently not possible to deprecate feature block struct items via feature flagging. Feature made redundant/ineffective by a breaking API change.
	SkipShutdownAndForceDelete       bool
}

type VirtualMachineScaleSetFeatures struct {
	ForceDelete               bool
	ReimageOnManualUpgrade    bool
	RollInstancesWhenRequired bool
	ScaleToZeroOnDelete       bool
}

type KeyVaultFeatures struct {
	PurgeSoftDeleteOnDestroy         bool
	PurgeSoftDeletedKeysOnDestroy    bool
	PurgeSoftDeletedCertsOnDestroy   bool
	PurgeSoftDeletedSecretsOnDestroy bool
	PurgeSoftDeletedHSMsOnDestroy    bool
	PurgeSoftDeletedHSMKeysOnDestroy bool
	RecoverSoftDeletedKeyVaults      bool
	RecoverSoftDeletedKeys           bool
	RecoverSoftDeletedCerts          bool
	RecoverSoftDeletedSecrets        bool
	RecoverSoftDeletedHSMKeys        bool
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

type ManagedDiskFeatures struct {
	ExpandWithoutDowntime bool
}

type AppConfigurationFeatures struct {
	PurgeSoftDeleteOnDestroy bool
	RecoverSoftDeleted       bool
}

type StorageFeatures struct {
	DataPlaneAvailable bool
}

type SubscriptionFeatures struct {
	PreventCancellationOnDestroy bool
}

type RecoveryServicesVault struct {
	RecoverSoftDeletedBackupProtectedVM bool
}

type PostgresqlFlexibleServerFeatures struct {
	RestartServerOnConfigurationValueChange bool
}

type MachineLearningFeatures struct {
	PurgeSoftDeletedWorkspaceOnDestroy bool
}

type RecoveryServiceFeatures struct {
	VMBackupStopProtectionAndRetainDataOnDestroy    bool
	VMBackupSuspendProtectionAndRetainDataOnDestroy bool
	PurgeProtectedItemsFromVaultOnDestroy           bool
}

type NetAppFeatures struct {
	DeleteBackupsOnBackupVaultDestroy bool
	PreventVolumeDestruction          bool
}

type DatabricksWorkspaceFeatures struct {
	ForceDelete bool
}
