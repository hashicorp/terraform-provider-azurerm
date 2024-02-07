// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package features

func Default() UserFeatures {
	return UserFeatures{
		// NOTE: ensure all nested objects are fully populated
		ApiManagement: ApiManagementFeatures{
			PurgeSoftDeleteOnDestroy: true,
			RecoverSoftDeleted:       true,
		},
		AppConfiguration: AppConfigurationFeatures{
			PurgeSoftDeleteOnDestroy: true,
			RecoverSoftDeleted:       true,
		},
		ApplicationInsights: ApplicationInsightFeatures{
			DisableGeneratedRule: false,
		},
		CognitiveAccount: CognitiveAccountFeatures{
			PurgeSoftDeleteOnDestroy: true,
		},
		KeyVault: KeyVaultFeatures{
			PurgeSoftDeleteOnDestroy:         true,
			PurgeSoftDeletedKeysOnDestroy:    true,
			PurgeSoftDeletedCertsOnDestroy:   true,
			PurgeSoftDeletedSecretsOnDestroy: true,
			PurgeSoftDeletedHSMsOnDestroy:    true,
			RecoverSoftDeletedKeyVaults:      true,
			RecoverSoftDeletedKeys:           true,
			RecoverSoftDeletedCerts:          true,
			RecoverSoftDeletedSecrets:        true,
		},
		LogAnalyticsWorkspace: LogAnalyticsWorkspaceFeatures{
			PermanentlyDeleteOnDestroy: true,
		},
		ManagedDisk: ManagedDiskFeatures{
			ExpandWithoutDowntime: true,
		},
		ResourceGroup: ResourceGroupFeatures{
			PreventDeletionIfContainsResources: true,
		},
		TemplateDeployment: TemplateDeploymentFeatures{
			DeleteNestedItemsDuringDeletion: true,
		},
		VirtualMachine: VirtualMachineFeatures{
			DeleteOSDiskOnDeletion:     true,
			GracefulShutdown:           false,
			SkipShutdownAndForceDelete: false,
		},
		VirtualMachineScaleSet: VirtualMachineScaleSetFeatures{
			ForceDelete:               false,
			RollInstancesWhenRequired: true,
			ScaleToZeroOnDelete:       true,
		},
		Subscription: SubscriptionFeatures{
			PreventCancellationOnDestroy: false,
		},
	}
}
