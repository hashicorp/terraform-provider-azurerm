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
		EnhancedValidation: EnhancedValidationFeatures{
			Locations:         false, // azignore:AZG007 - ensure all nested objects are fully populated
			ResourceProviders: false, // azignore:AZG007 - ensure all nested objects are fully populated
			PreflightEnabled:  false, // azignore:AZG007 - ensure all nested objects are fully populated
			LocationFallback:  nil, // azignore:AZG007 - ensure all nested objects are fully populated
		},
		AppConfiguration: AppConfigurationFeatures{
			PurgeSoftDeleteOnDestroy: true,
			RecoverSoftDeleted:       true,
		},
		ApplicationInsights: ApplicationInsightFeatures{
			DisableGeneratedRule: false, // azignore:AZG007 - ensure all nested objects are fully populated
		},
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
		LogAnalyticsWorkspace: LogAnalyticsWorkspaceFeatures{
			PermanentlyDeleteOnDestroy: false, // azignore:AZG007 - ensure all nested objects are fully populated
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
			DetachImplicitDataDiskOnDeletion: false, // azignore:AZG007 - ensure all nested objects are fully populated
			DeleteOSDiskOnDeletion:           true,
			SkipShutdownAndForceDelete:       false, // azignore:AZG007 - ensure all nested objects are fully populated
		},
		VirtualMachineScaleSet: VirtualMachineScaleSetFeatures{
			ForceDelete:               false, // azignore:AZG007 - ensure all nested objects are fully populated
			ReimageOnManualUpgrade:    true,
			RollInstancesWhenRequired: true,
			ScaleToZeroOnDelete:       true,
		},
		Storage: StorageFeatures{
			DataPlaneAvailable: true,
		},
		Subscription: SubscriptionFeatures{
			PreventCancellationOnDestroy: false, // azignore:AZG007 - ensure all nested objects are fully populated
		},
		PostgresqlFlexibleServer: PostgresqlFlexibleServerFeatures{
			RestartServerOnConfigurationValueChange: true,
		},
		MachineLearning: MachineLearningFeatures{
			PurgeSoftDeletedWorkspaceOnDestroy: false, // azignore:AZG007 - ensure all nested objects are fully populated
		},
		RecoveryService: RecoveryServiceFeatures{
			VMBackupStopProtectionAndRetainDataOnDestroy:    false, // azignore:AZG007 - ensure all nested objects are fully populated
			VMBackupSuspendProtectionAndRetainDataOnDestroy: false, // azignore:AZG007 - ensure all nested objects are fully populated
			PurgeProtectedItemsFromVaultOnDestroy:           false, // azignore:AZG007 - ensure all nested objects are fully populated
		},
		NetApp: NetAppFeatures{
			DeleteBackupsOnBackupVaultDestroy: false, // azignore:AZG007 - ensure all nested objects are fully populated
			PreventVolumeDestruction:          true,
		},
		DatabricksWorkspace: DatabricksWorkspaceFeatures{
			ForceDelete: false, // azignore:AZG007 - ensure all nested objects are fully populated
		},
		ServiceBus: ServiceBusFeatures{
			AutoDeleteSubscriptionDefaultRule: false, // azignore:AZG007 - ensure all nested objects are fully populated
		},
	}
}