package features

func Default() UserFeatures {
	return UserFeatures{
		// NOTE: ensure all nested objects are fully populated
		KeyVault: KeyVaultFeatures{
			PurgeSoftDeleteOnDestroy:    true,
			RecoverSoftDeletedKeyVaults: true,
		},
		LogAnalyticsWorkspace: LogAnalyticsWorkspaceFeatures{
			PermanentlyDeleteOnDestroy: false,
		},
		Network: NetworkFeatures{
			RelaxedLocking: false,
		},
		TemplateDeployment: TemplateDeploymentFeatures{
			DeleteNestedItemsDuringDeletion: true,
		},
		VirtualMachine: VirtualMachineFeatures{
			DeleteOSDiskOnDeletion: true,
			GracefulShutdown:       false,
		},
		VirtualMachineScaleSet: VirtualMachineScaleSetFeatures{
			RollInstancesWhenRequired: true,
		},
	}
}
