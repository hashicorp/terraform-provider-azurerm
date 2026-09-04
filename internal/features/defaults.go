// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package features

func Default() UserFeatures {
	return UserFeatures{
		// NOTE: ensure all nested objects are fully populated
		ApiManagement: ApiManagementFeatures{
			PurgeSoftDeleteOnDestroy: true,
			RecoverSoftDeleted:       true,
		},
		EnhancedValidation: EnhancedValidationFeatures{},
		AppConfiguration: AppConfigurationFeatures{
			PurgeSoftDeleteOnDestroy: true,
			RecoverSoftDeleted:       true,
		},
		ApplicationInsights: ApplicationInsightFeatures{},
		CognitiveAccount: CognitiveAccountFeatures{
			PurgeSoftDeleteOnDestroy: true,
		},
		KeyVault: KeyVaultFeatures{
			// Standard
			PurgeSoftDeleteOnDestroy:         true,
			PurgeSoftDeletedKeysOnDestroy:    true,
			PurgeSoftDeletedCertsOnDestroy:   true,
			PurgeSoftDeletedSecretsOnDestroy: true,
			RecoverSoftDeletedKeyVaults:      true,
			RecoverSoftDeletedKeys:           true,
			RecoverSoftDeletedCerts:          true,
			RecoverSoftDeletedSecrets:        true,

			// HSM
			PurgeSoftDeletedHSMsOnDestroy:    true,
			PurgeSoftDeletedHSMKeysOnDestroy: true,
			RecoverSoftDeletedHSMKeys:        true,
		},
		LogAnalyticsWorkspace: LogAnalyticsWorkspaceFeatures{},
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
			DeleteOSDiskOnDeletion: true,
		},
		VirtualMachineScaleSet: VirtualMachineScaleSetFeatures{
			ReimageOnManualUpgrade:    true,
			RollInstancesWhenRequired: true,
			ScaleToZeroOnDelete:       true,
		},
		Storage: StorageFeatures{
			DataPlaneAvailable: true,
		},
		Subscription: SubscriptionFeatures{},
		PostgresqlFlexibleServer: PostgresqlFlexibleServerFeatures{
			RestartServerOnConfigurationValueChange: true,
		},
		MachineLearning: MachineLearningFeatures{},
		RecoveryService: RecoveryServiceFeatures{},
		NetApp: NetAppFeatures{
			PreventVolumeDestruction: true,
		},
		DatabricksWorkspace: DatabricksWorkspaceFeatures{},
		ServiceBus:          ServiceBusFeatures{},
	}
}
