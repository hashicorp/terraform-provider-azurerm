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
			RecoverSoftDeletedKeyVaults:      true,
			RecoverSoftDeletedKeys:           true,
			RecoverSoftDeletedCerts:          true,
			RecoverSoftDeletedSecrets:        true,

			// todo 4.0 move all HSM flags into their own features HSMFeatures block
			PurgeSoftDeletedHSMsOnDestroy:    true,
			PurgeSoftDeletedHSMKeysOnDestroy: true,
			RecoverSoftDeletedHSMKeys:        true,
		},
		LogAnalyticsWorkspace: LogAnalyticsWorkspaceFeatures{
			PermanentlyDeleteOnDestroy: false,
		},
		ManagedDisk: ManagedDiskFeatures{
			ExpandWithoutDowntime: true,
		},
		ResourceGroup: ResourceGroupFeatures{
			PreventDeletionIfContainsResources: true,
		},
		RecoveryServicesVault: RecoveryServicesVault{
			RecoverSoftDeletedBackupProtectedVM: true,
		},
		TemplateDeployment: TemplateDeploymentFeatures{
			DeleteNestedItemsDuringDeletion: true,
		},
		VirtualMachine: VirtualMachineFeatures{
			DetachImplicitDataDiskOnDeletion: false,
			DeleteOSDiskOnDeletion:           true,
			GracefulShutdown:                 false,
			SkipShutdownAndForceDelete:       false,
		},
		VirtualMachineScaleSet: VirtualMachineScaleSetFeatures{
			ForceDelete:               false,
			ReimageOnManualUpgrade:    true,
			RollInstancesWhenRequired: true,
			ScaleToZeroOnDelete:       true,
		},
		Subscription: SubscriptionFeatures{
			PreventCancellationOnDestroy: false,
		},
		PostgresqlFlexibleServer: PostgresqlFlexibleServerFeatures{
			RestartServerOnConfigurationValueChange: true,
		},
		MachineLearning: MachineLearningFeatures{
			PurgeSoftDeletedWorkspaceOnDestroy: false,
		},
		RecoveryService: RecoveryServiceFeatures{
			VMBackupStopProtectionAndRetainDataOnDestroy: false,
			PurgeProtectedItemsFromVaultOnDestroy:        false,
		},
	}
}
