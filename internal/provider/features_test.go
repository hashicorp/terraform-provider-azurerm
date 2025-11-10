// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
)

func TestExpandFeatures(t *testing.T) {
	testData := []struct {
		Name     string
		Input    []interface{}
		EnvVars  map[string]interface{}
		Expected features.UserFeatures
	}{
		{
			Name:  "Empty Block",
			Input: []interface{}{},
			Expected: features.UserFeatures{
				ApiManagement: features.ApiManagementFeatures{
					PurgeSoftDeleteOnDestroy: true,
					RecoverSoftDeleted:       true,
				},
				AppConfiguration: features.AppConfigurationFeatures{
					PurgeSoftDeleteOnDestroy: true,
					RecoverSoftDeleted:       true,
				},
				ApplicationInsights: features.ApplicationInsightFeatures{
					DisableGeneratedRule: false,
				},
				CognitiveAccount: features.CognitiveAccountFeatures{
					PurgeSoftDeleteOnDestroy: true,
				},
				KeyVault: features.KeyVaultFeatures{
					PurgeSoftDeletedCertsOnDestroy:   true,
					PurgeSoftDeletedKeysOnDestroy:    true,
					PurgeSoftDeletedSecretsOnDestroy: true,
					PurgeSoftDeleteOnDestroy:         true,
					PurgeSoftDeletedHSMsOnDestroy:    true,
					PurgeSoftDeletedHSMKeysOnDestroy: true,
					RecoverSoftDeletedCerts:          true,
					RecoverSoftDeletedKeys:           true,
					RecoverSoftDeletedKeyVaults:      true,
					RecoverSoftDeletedSecrets:        true,
					RecoverSoftDeletedHSMKeys:        true,
				},
				LogAnalyticsWorkspace: features.LogAnalyticsWorkspaceFeatures{
					PermanentlyDeleteOnDestroy: false,
				},
				ManagedDisk: features.ManagedDiskFeatures{
					ExpandWithoutDowntime: true,
				},
				TemplateDeployment: features.TemplateDeploymentFeatures{
					DeleteNestedItemsDuringDeletion: true,
				},
				VirtualMachine: features.VirtualMachineFeatures{
					DetachImplicitDataDiskOnDeletion: false,
					DeleteOSDiskOnDeletion:           true,
					SkipShutdownAndForceDelete:       false,
				},
				VirtualMachineScaleSet: features.VirtualMachineScaleSetFeatures{
					ForceDelete:               false,
					ReimageOnManualUpgrade:    true,
					RollInstancesWhenRequired: true,
					ScaleToZeroOnDelete:       true,
				},
				ResourceGroup: features.ResourceGroupFeatures{
					PreventDeletionIfContainsResources: true,
				},
				RecoveryServicesVault: features.RecoveryServicesVault{
					RecoverSoftDeletedBackupProtectedVM: true,
				},
				Storage: features.StorageFeatures{
					DataPlaneAvailable: true,
				},
				Subscription: features.SubscriptionFeatures{
					PreventCancellationOnDestroy: false,
				},
				PostgresqlFlexibleServer: features.PostgresqlFlexibleServerFeatures{
					RestartServerOnConfigurationValueChange: true,
				},
				MachineLearning: features.MachineLearningFeatures{
					PurgeSoftDeletedWorkspaceOnDestroy: false,
				},
				RecoveryService: features.RecoveryServiceFeatures{
					VMBackupStopProtectionAndRetainDataOnDestroy:    false,
					VMBackupSuspendProtectionAndRetainDataOnDestroy: false,
					PurgeProtectedItemsFromVaultOnDestroy:           false,
				},
				NetApp: features.NetAppFeatures{
					DeleteBackupsOnBackupVaultDestroy: false,
					PreventVolumeDestruction:          true,
				},
				DatabricksWorkspace: features.DatabricksWorkspaceFeatures{
					ForceDelete: false,
				},
			},
		},
		{
			Name: "Complete Enabled",
			Input: []interface{}{
				map[string]interface{}{
					"api_management": []interface{}{
						map[string]interface{}{
							"purge_soft_delete_on_destroy": true,
							"recover_soft_deleted":         true,
						},
					},
					"app_configuration": []interface{}{
						map[string]interface{}{
							"purge_soft_delete_on_destroy": true,
							"recover_soft_deleted":         true,
						},
					},
					"application_insights": []interface{}{
						map[string]interface{}{
							"disable_generated_rule": true,
						},
					},
					"cognitive_account": []interface{}{
						map[string]interface{}{
							"purge_soft_delete_on_destroy": true,
						},
					},
					"key_vault": []interface{}{
						map[string]interface{}{
							"purge_soft_deleted_certificates_on_destroy":                  true,
							"purge_soft_deleted_keys_on_destroy":                          true,
							"purge_soft_deleted_secrets_on_destroy":                       true,
							"purge_soft_deleted_hardware_security_modules_on_destroy":     true,
							"purge_soft_deleted_hardware_security_module_keys_on_destroy": true,
							"purge_soft_delete_on_destroy":                                true,
							"recover_soft_deleted_certificates":                           true,
							"recover_soft_deleted_keys":                                   true,
							"recover_soft_deleted_key_vaults":                             true,
							"recover_soft_deleted_secrets":                                true,
							"recover_soft_deleted_hardware_security_module_keys":          true,
						},
					},
					"log_analytics_workspace": []interface{}{
						map[string]interface{}{
							"permanently_delete_on_destroy": true,
						},
					},
					"managed_disk": []interface{}{
						map[string]interface{}{
							"expand_without_downtime": true,
						},
					},
					"postgresql_flexible_server": []interface{}{
						map[string]interface{}{
							"restart_server_on_configuration_value_change": true,
						},
					},
					"resource_group": []interface{}{
						map[string]interface{}{
							"prevent_deletion_if_contains_resources": true,
						},
					},
					"recovery_services_vaults": []interface{}{
						map[string]interface{}{
							"recover_soft_deleted_backup_protected_vm": true,
						},
					},
					"storage": []interface{}{
						map[string]interface{}{
							"data_plane_available": true,
						},
					},
					"subscription": []interface{}{
						map[string]interface{}{
							"prevent_cancellation_on_destroy": true,
						},
					},
					"template_deployment": []interface{}{
						map[string]interface{}{
							"delete_nested_items_during_deletion": true,
						},
					},
					"virtual_machine": []interface{}{
						map[string]interface{}{
							"detach_implicit_data_disk_on_deletion": true,
							"delete_os_disk_on_deletion":            true,
							"skip_shutdown_and_force_delete":        true,
						},
					},
					"virtual_machine_scale_set": []interface{}{
						map[string]interface{}{
							"reimage_on_manual_upgrade":     true,
							"roll_instances_when_required":  true,
							"force_delete":                  true,
							"scale_to_zero_before_deletion": true,
						},
					},
					"machine_learning": []interface{}{
						map[string]interface{}{
							"purge_soft_deleted_workspace_on_destroy": true,
						},
					},
					"recovery_service": []interface{}{
						map[string]interface{}{
							"vm_backup_stop_protection_and_retain_data_on_destroy":    true,
							"vm_backup_suspend_protection_and_retain_data_on_destroy": true,
							"purge_protected_items_from_vault_on_destroy":             true,
						},
					},
					"netapp": []interface{}{
						map[string]interface{}{
							"delete_backups_on_backup_vault_destroy": true,
							"prevent_volume_destruction":             true,
						},
					},
					"databricks_workspace": []interface{}{
						map[string]interface{}{
							"force_delete": true,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				ApiManagement: features.ApiManagementFeatures{
					PurgeSoftDeleteOnDestroy: true,
					RecoverSoftDeleted:       true,
				},
				AppConfiguration: features.AppConfigurationFeatures{
					PurgeSoftDeleteOnDestroy: true,
					RecoverSoftDeleted:       true,
				},
				ApplicationInsights: features.ApplicationInsightFeatures{
					DisableGeneratedRule: true,
				},
				CognitiveAccount: features.CognitiveAccountFeatures{
					PurgeSoftDeleteOnDestroy: true,
				},
				KeyVault: features.KeyVaultFeatures{
					PurgeSoftDeletedCertsOnDestroy:   true,
					PurgeSoftDeletedKeysOnDestroy:    true,
					PurgeSoftDeletedSecretsOnDestroy: true,
					PurgeSoftDeleteOnDestroy:         true,
					PurgeSoftDeletedHSMsOnDestroy:    true,
					PurgeSoftDeletedHSMKeysOnDestroy: true,
					RecoverSoftDeletedCerts:          true,
					RecoverSoftDeletedKeys:           true,
					RecoverSoftDeletedKeyVaults:      true,
					RecoverSoftDeletedSecrets:        true,
					RecoverSoftDeletedHSMKeys:        true,
				},
				LogAnalyticsWorkspace: features.LogAnalyticsWorkspaceFeatures{
					PermanentlyDeleteOnDestroy: true,
				},
				ManagedDisk: features.ManagedDiskFeatures{
					ExpandWithoutDowntime: true,
				},
				ResourceGroup: features.ResourceGroupFeatures{
					PreventDeletionIfContainsResources: true,
				},
				RecoveryServicesVault: features.RecoveryServicesVault{
					RecoverSoftDeletedBackupProtectedVM: true,
				},
				Storage: features.StorageFeatures{
					DataPlaneAvailable: true,
				},
				Subscription: features.SubscriptionFeatures{
					PreventCancellationOnDestroy: true,
				},
				TemplateDeployment: features.TemplateDeploymentFeatures{
					DeleteNestedItemsDuringDeletion: true,
				},
				VirtualMachine: features.VirtualMachineFeatures{
					DetachImplicitDataDiskOnDeletion: true,
					DeleteOSDiskOnDeletion:           true,
					SkipShutdownAndForceDelete:       true,
				},
				VirtualMachineScaleSet: features.VirtualMachineScaleSetFeatures{
					ReimageOnManualUpgrade:    true,
					RollInstancesWhenRequired: true,
					ForceDelete:               true,
					ScaleToZeroOnDelete:       true,
				},
				PostgresqlFlexibleServer: features.PostgresqlFlexibleServerFeatures{
					RestartServerOnConfigurationValueChange: true,
				},
				MachineLearning: features.MachineLearningFeatures{
					PurgeSoftDeletedWorkspaceOnDestroy: true,
				},
				RecoveryService: features.RecoveryServiceFeatures{
					VMBackupStopProtectionAndRetainDataOnDestroy:    true,
					VMBackupSuspendProtectionAndRetainDataOnDestroy: true,
					PurgeProtectedItemsFromVaultOnDestroy:           true,
				},
				NetApp: features.NetAppFeatures{
					DeleteBackupsOnBackupVaultDestroy: true,
					PreventVolumeDestruction:          true,
				},
				DatabricksWorkspace: features.DatabricksWorkspaceFeatures{
					ForceDelete: true,
				},
			},
		},
		{
			Name: "Complete Disabled",
			Input: []interface{}{
				map[string]interface{}{
					"api_management": []interface{}{
						map[string]interface{}{
							"purge_soft_delete_on_destroy": false,
							"recover_soft_deleted":         false,
						},
					},
					"app_configuration": []interface{}{
						map[string]interface{}{
							"purge_soft_delete_on_destroy": false,
							"recover_soft_deleted":         false,
						},
					},
					"application_insights": []interface{}{
						map[string]interface{}{
							"disable_generated_rule": false,
						},
					},
					"cognitive_account": []interface{}{
						map[string]interface{}{
							"purge_soft_delete_on_destroy": false,
						},
					},
					"key_vault": []interface{}{
						map[string]interface{}{
							"purge_soft_deleted_certificates_on_destroy":                  false,
							"purge_soft_deleted_keys_on_destroy":                          false,
							"purge_soft_deleted_secrets_on_destroy":                       false,
							"purge_soft_deleted_hardware_security_modules_on_destroy":     false,
							"purge_soft_deleted_hardware_security_module_keys_on_destroy": false,
							"purge_soft_delete_on_destroy":                                false,
							"recover_soft_deleted_certificates":                           false,
							"recover_soft_deleted_keys":                                   false,
							"recover_soft_deleted_key_vaults":                             false,
							"recover_soft_deleted_secrets":                                false,
							"recover_soft_deleted_hardware_security_module_keys":          false,
						},
					},
					"log_analytics_workspace": []interface{}{
						map[string]interface{}{
							"permanently_delete_on_destroy": false,
						},
					},
					"managed_disk": []interface{}{
						map[string]interface{}{
							"expand_without_downtime": false,
						},
					},
					"postgresql_flexible_server": []interface{}{
						map[string]interface{}{
							"restart_server_on_configuration_value_change": false,
						},
					},
					"resource_group": []interface{}{
						map[string]interface{}{
							"prevent_deletion_if_contains_resources": false,
						},
					},
					"recovery_services_vaults": []interface{}{
						map[string]interface{}{
							"recover_soft_deleted_backup_protected_vm": false,
						},
					},
					"storage": []interface{}{
						map[string]interface{}{
							"data_plane_available": false,
						},
					},
					"subscription": []interface{}{
						map[string]interface{}{
							"prevent_cancellation_on_destroy": false,
						},
					},
					"template_deployment": []interface{}{
						map[string]interface{}{
							"delete_nested_items_during_deletion": false,
						},
					},
					"virtual_machine": []interface{}{
						map[string]interface{}{
							"detach_implicit_data_disk_on_deletion": false,
							"delete_os_disk_on_deletion":            false,
							"skip_shutdown_and_force_delete":        false,
						},
					},
					"virtual_machine_scale_set": []interface{}{
						map[string]interface{}{
							"force_delete":                  false,
							"reimage_on_manual_upgrade":     false,
							"roll_instances_when_required":  false,
							"scale_to_zero_before_deletion": false,
						},
					},
					"machine_learning": []interface{}{
						map[string]interface{}{
							"purge_soft_deleted_workspace_on_destroy": false,
						},
					},
					"recovery_service": []interface{}{
						map[string]interface{}{
							"vm_backup_stop_protection_and_retain_data_on_destroy":    false,
							"vm_backup_suspend_protection_and_retain_data_on_destroy": false,
							"purge_protected_items_from_vault_on_destroy":             false,
						},
					},
					"netapp": []interface{}{
						map[string]interface{}{
							"delete_backups_on_backup_vault_destroy": false,
							"prevent_volume_destruction":             false,
						},
					},
					"databricks_workspace": []interface{}{
						map[string]interface{}{
							"force_delete": false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				ApiManagement: features.ApiManagementFeatures{
					PurgeSoftDeleteOnDestroy: false,
					RecoverSoftDeleted:       false,
				},
				AppConfiguration: features.AppConfigurationFeatures{
					PurgeSoftDeleteOnDestroy: false,
					RecoverSoftDeleted:       false,
				},
				ApplicationInsights: features.ApplicationInsightFeatures{
					DisableGeneratedRule: false,
				},
				CognitiveAccount: features.CognitiveAccountFeatures{
					PurgeSoftDeleteOnDestroy: false,
				},
				KeyVault: features.KeyVaultFeatures{
					PurgeSoftDeletedCertsOnDestroy:   false,
					PurgeSoftDeletedKeysOnDestroy:    false,
					PurgeSoftDeletedSecretsOnDestroy: false,
					PurgeSoftDeletedHSMsOnDestroy:    false,
					PurgeSoftDeletedHSMKeysOnDestroy: false,
					PurgeSoftDeleteOnDestroy:         false,
					RecoverSoftDeletedCerts:          false,
					RecoverSoftDeletedKeys:           false,
					RecoverSoftDeletedKeyVaults:      false,
					RecoverSoftDeletedSecrets:        false,
					RecoverSoftDeletedHSMKeys:        false,
				},
				LogAnalyticsWorkspace: features.LogAnalyticsWorkspaceFeatures{
					PermanentlyDeleteOnDestroy: false,
				},
				ManagedDisk: features.ManagedDiskFeatures{
					ExpandWithoutDowntime: false,
				},
				ResourceGroup: features.ResourceGroupFeatures{
					PreventDeletionIfContainsResources: false,
				},
				RecoveryServicesVault: features.RecoveryServicesVault{
					RecoverSoftDeletedBackupProtectedVM: false,
				},
				Storage: features.StorageFeatures{
					DataPlaneAvailable: false,
				},
				Subscription: features.SubscriptionFeatures{
					PreventCancellationOnDestroy: false,
				},
				TemplateDeployment: features.TemplateDeploymentFeatures{
					DeleteNestedItemsDuringDeletion: false,
				},
				VirtualMachine: features.VirtualMachineFeatures{
					DetachImplicitDataDiskOnDeletion: false,
					DeleteOSDiskOnDeletion:           false,
					SkipShutdownAndForceDelete:       false,
				},
				VirtualMachineScaleSet: features.VirtualMachineScaleSetFeatures{
					ForceDelete:               false,
					ReimageOnManualUpgrade:    false,
					RollInstancesWhenRequired: false,
					ScaleToZeroOnDelete:       false,
				},
				PostgresqlFlexibleServer: features.PostgresqlFlexibleServerFeatures{
					RestartServerOnConfigurationValueChange: false,
				},
				MachineLearning: features.MachineLearningFeatures{
					PurgeSoftDeletedWorkspaceOnDestroy: false,
				},
				RecoveryService: features.RecoveryServiceFeatures{
					VMBackupStopProtectionAndRetainDataOnDestroy:    false,
					VMBackupSuspendProtectionAndRetainDataOnDestroy: false,
					PurgeProtectedItemsFromVaultOnDestroy:           false,
				},
				NetApp: features.NetAppFeatures{
					DeleteBackupsOnBackupVaultDestroy: false,
					PreventVolumeDestruction:          false,
				},
				DatabricksWorkspace: features.DatabricksWorkspaceFeatures{
					ForceDelete: false,
				},
			},
		},
	}

	for _, testCase := range testData {
		t.Logf("[DEBUG] Test Case: %q", testCase.Name)
		result := expandFeatures(testCase.Input)
		if !reflect.DeepEqual(result, testCase.Expected) {
			t.Fatalf("Expected %+v but got %+v", result, testCase.Expected)
		}
	}
}

func TestExpandFeaturesApiManagement(t *testing.T) {
	testData := []struct {
		Name     string
		Input    []interface{}
		EnvVars  map[string]interface{}
		Expected features.UserFeatures
	}{
		{
			Name: "Empty Block",
			Input: []interface{}{
				map[string]interface{}{
					"api_management": []interface{}{},
				},
			},
			Expected: features.UserFeatures{
				ApiManagement: features.ApiManagementFeatures{
					PurgeSoftDeleteOnDestroy: true,
					RecoverSoftDeleted:       true,
				},
			},
		},
		{
			Name: "Purge Soft Delete On Destroy and Recover Soft Deleted Api Management Enabled",
			Input: []interface{}{
				map[string]interface{}{
					"api_management": []interface{}{
						map[string]interface{}{
							"purge_soft_delete_on_destroy": true,
							"recover_soft_deleted":         true,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				ApiManagement: features.ApiManagementFeatures{
					PurgeSoftDeleteOnDestroy: true,
					RecoverSoftDeleted:       true,
				},
			},
		},
		{
			Name: "Purge Soft Delete On Destroy and Recover Soft Deleted Api Management Disabled",
			Input: []interface{}{
				map[string]interface{}{
					"api_management": []interface{}{
						map[string]interface{}{
							"purge_soft_delete_on_destroy": false,
							"recover_soft_deleted":         false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				ApiManagement: features.ApiManagementFeatures{
					PurgeSoftDeleteOnDestroy: false,
					RecoverSoftDeleted:       false,
				},
			},
		},
	}

	for _, testCase := range testData {
		t.Logf("[DEBUG] Test Case: %q", testCase.Name)
		result := expandFeatures(testCase.Input)
		if !reflect.DeepEqual(result.ApiManagement, testCase.Expected.ApiManagement) {
			t.Fatalf("Expected %+v but got %+v", result.ApiManagement, testCase.Expected.ApiManagement)
		}
	}
}

func TestExpandFeaturesAppConfiguration(t *testing.T) {
	testData := []struct {
		Name     string
		Input    []interface{}
		EnvVars  map[string]interface{}
		Expected features.UserFeatures
	}{
		{
			Name: "Empty Block",
			Input: []interface{}{
				map[string]interface{}{
					"app_configuration": []interface{}{},
				},
			},
			Expected: features.UserFeatures{
				AppConfiguration: features.AppConfigurationFeatures{
					PurgeSoftDeleteOnDestroy: true,
					RecoverSoftDeleted:       true,
				},
			},
		},
		{
			Name: "Purge Soft Delete On Destroy and Recover Soft Deleted App Configuration Enabled",
			Input: []interface{}{
				map[string]interface{}{
					"app_configuration": []interface{}{
						map[string]interface{}{
							"purge_soft_delete_on_destroy": true,
							"recover_soft_deleted":         true,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				AppConfiguration: features.AppConfigurationFeatures{
					PurgeSoftDeleteOnDestroy: true,
					RecoverSoftDeleted:       true,
				},
			},
		},
		{
			Name: "Purge Soft Delete On Destroy and Recover Soft Deleted App Configuration Disabled",
			Input: []interface{}{
				map[string]interface{}{
					"app_configuration": []interface{}{
						map[string]interface{}{
							"purge_soft_delete_on_destroy": false,
							"recover_soft_deleted":         false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				AppConfiguration: features.AppConfigurationFeatures{
					PurgeSoftDeleteOnDestroy: false,
					RecoverSoftDeleted:       false,
				},
			},
		},
	}

	for _, testCase := range testData {
		t.Logf("[DEBUG] Test Case: %q", testCase.Name)
		result := expandFeatures(testCase.Input)
		if !reflect.DeepEqual(result.AppConfiguration, testCase.Expected.AppConfiguration) {
			t.Fatalf("Expected %+v but got %+v", result.AppConfiguration, testCase.Expected.AppConfiguration)
		}
	}
}

func TestExpandFeaturesApplicationInsights(t *testing.T) {
	testData := []struct {
		Name     string
		Input    []interface{}
		EnvVars  map[string]interface{}
		Expected features.UserFeatures
	}{
		{
			Name: "Empty Block",
			Input: []interface{}{
				map[string]interface{}{
					"application_insights": []interface{}{},
				},
			},
			Expected: features.UserFeatures{
				ApplicationInsights: features.ApplicationInsightFeatures{
					DisableGeneratedRule: false,
				},
			},
		},
		{
			Name: "Disable Generated Rule",
			Input: []interface{}{
				map[string]interface{}{
					"application_insights": []interface{}{
						map[string]interface{}{
							"disable_generated_rule": true,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				ApplicationInsights: features.ApplicationInsightFeatures{
					DisableGeneratedRule: true,
				},
			},
		},
		{
			Name: "Enable Generated Rule",
			Input: []interface{}{
				map[string]interface{}{
					"application_insights": []interface{}{
						map[string]interface{}{
							"disable_generated_rule": false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				ApplicationInsights: features.ApplicationInsightFeatures{
					DisableGeneratedRule: false,
				},
			},
		},
	}

	for _, testCase := range testData {
		t.Logf("[DEBUG] Test Case: %q", testCase.Name)
		result := expandFeatures(testCase.Input)
		if !reflect.DeepEqual(result.ApplicationInsights, testCase.Expected.ApplicationInsights) {
			t.Fatalf("Expected %+v but got %+v", result.ApplicationInsights, testCase.Expected.ApplicationInsights)
		}
	}
}

func TestExpandFeaturesCognitiveServices(t *testing.T) {
	testData := []struct {
		Name     string
		Input    []interface{}
		EnvVars  map[string]interface{}
		Expected features.UserFeatures
	}{
		{
			Name: "Empty Block",
			Input: []interface{}{
				map[string]interface{}{
					"cognitive_account": []interface{}{},
				},
			},
			Expected: features.UserFeatures{
				CognitiveAccount: features.CognitiveAccountFeatures{
					PurgeSoftDeleteOnDestroy: true,
				},
			},
		},
		{
			Name: "Purge on Destroy Enabled",
			Input: []interface{}{
				map[string]interface{}{
					"cognitive_account": []interface{}{
						map[string]interface{}{
							"purge_soft_delete_on_destroy": true,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				CognitiveAccount: features.CognitiveAccountFeatures{
					PurgeSoftDeleteOnDestroy: true,
				},
			},
		},
		{
			Name: "Purge on Destroy Disabled",
			Input: []interface{}{
				map[string]interface{}{
					"cognitive_account": []interface{}{
						map[string]interface{}{
							"purge_soft_delete_on_destroy": false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				CognitiveAccount: features.CognitiveAccountFeatures{
					PurgeSoftDeleteOnDestroy: false,
				},
			},
		},
	}

	for _, testCase := range testData {
		t.Logf("[DEBUG] Test Case: %q", testCase.Name)
		result := expandFeatures(testCase.Input)
		if !reflect.DeepEqual(result.CognitiveAccount, testCase.Expected.CognitiveAccount) {
			t.Fatalf("Expected %+v but got %+v", result.CognitiveAccount, testCase.Expected.CognitiveAccount)
		}
	}
}

func TestExpandFeaturesKeyVault(t *testing.T) {
	testData := []struct {
		Name     string
		Input    []interface{}
		EnvVars  map[string]interface{}
		Expected features.UserFeatures
	}{
		{
			Name: "Empty Block",
			Input: []interface{}{
				map[string]interface{}{
					"key_vault": []interface{}{},
				},
			},
			Expected: features.UserFeatures{
				KeyVault: features.KeyVaultFeatures{
					PurgeSoftDeletedCertsOnDestroy:   true,
					PurgeSoftDeletedKeysOnDestroy:    true,
					PurgeSoftDeletedSecretsOnDestroy: true,
					PurgeSoftDeleteOnDestroy:         true,
					PurgeSoftDeletedHSMsOnDestroy:    true,
					PurgeSoftDeletedHSMKeysOnDestroy: true,
					RecoverSoftDeletedCerts:          true,
					RecoverSoftDeletedKeys:           true,
					RecoverSoftDeletedKeyVaults:      true,
					RecoverSoftDeletedSecrets:        true,
					RecoverSoftDeletedHSMKeys:        true,
				},
			},
		},
		{
			Name: "Purge Soft Delete On Destroy and Recover Soft Deleted Key Vaults Enabled",
			Input: []interface{}{
				map[string]interface{}{
					"key_vault": []interface{}{
						map[string]interface{}{
							"purge_soft_deleted_certificates_on_destroy":                  true,
							"purge_soft_deleted_keys_on_destroy":                          true,
							"purge_soft_deleted_secrets_on_destroy":                       true,
							"purge_soft_deleted_hardware_security_modules_on_destroy":     true,
							"purge_soft_deleted_hardware_security_module_keys_on_destroy": true,
							"purge_soft_delete_on_destroy":                                true,
							"recover_soft_deleted_certificates":                           true,
							"recover_soft_deleted_keys":                                   true,
							"recover_soft_deleted_key_vaults":                             true,
							"recover_soft_deleted_secrets":                                true,
							"recover_soft_deleted_hardware_security_module_keys":          true,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				KeyVault: features.KeyVaultFeatures{
					PurgeSoftDeletedCertsOnDestroy:   true,
					PurgeSoftDeletedKeysOnDestroy:    true,
					PurgeSoftDeletedSecretsOnDestroy: true,
					PurgeSoftDeletedHSMsOnDestroy:    true,
					PurgeSoftDeletedHSMKeysOnDestroy: true,
					PurgeSoftDeleteOnDestroy:         true,
					RecoverSoftDeletedCerts:          true,
					RecoverSoftDeletedKeys:           true,
					RecoverSoftDeletedKeyVaults:      true,
					RecoverSoftDeletedSecrets:        true,
					RecoverSoftDeletedHSMKeys:        true,
				},
			},
		},
		{
			Name: "Purge Soft Delete On Destroy and Recover Soft Deleted Key Vaults Disabled",
			Input: []interface{}{
				map[string]interface{}{
					"key_vault": []interface{}{
						map[string]interface{}{
							"purge_soft_deleted_certificates_on_destroy":                  false,
							"purge_soft_deleted_keys_on_destroy":                          false,
							"purge_soft_deleted_secrets_on_destroy":                       false,
							"purge_soft_deleted_hardware_security_modules_on_destroy":     false,
							"purge_soft_deleted_hardware_security_module_keys_on_destroy": false,
							"purge_soft_delete_on_destroy":                                false,
							"recover_soft_deleted_certificates":                           false,
							"recover_soft_deleted_keys":                                   false,
							"recover_soft_deleted_key_vaults":                             false,
							"recover_soft_deleted_secrets":                                false,
							"recover_soft_deleted_hardware_security_module_keys":          false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				KeyVault: features.KeyVaultFeatures{
					PurgeSoftDeletedCertsOnDestroy:   false,
					PurgeSoftDeletedKeysOnDestroy:    false,
					PurgeSoftDeletedSecretsOnDestroy: false,
					PurgeSoftDeleteOnDestroy:         false,
					PurgeSoftDeletedHSMsOnDestroy:    false,
					PurgeSoftDeletedHSMKeysOnDestroy: false,
					RecoverSoftDeletedCerts:          false,
					RecoverSoftDeletedKeyVaults:      false,
					RecoverSoftDeletedKeys:           false,
					RecoverSoftDeletedSecrets:        false,
					RecoverSoftDeletedHSMKeys:        false,
				},
			},
		},
	}

	for _, testCase := range testData {
		t.Logf("[DEBUG] Test Case: %q", testCase.Name)
		result := expandFeatures(testCase.Input)
		if !reflect.DeepEqual(result.KeyVault, testCase.Expected.KeyVault) {
			t.Fatalf("Expected %+v but got %+v", result.KeyVault, testCase.Expected.KeyVault)
		}
	}
}

func TestExpandFeaturesTemplateDeployment(t *testing.T) {
	testData := []struct {
		Name     string
		Input    []interface{}
		EnvVars  map[string]interface{}
		Expected features.UserFeatures
	}{
		{
			Name: "Empty Block",
			Input: []interface{}{
				map[string]interface{}{
					"template_deployment": []interface{}{},
				},
			},
			Expected: features.UserFeatures{
				TemplateDeployment: features.TemplateDeploymentFeatures{
					DeleteNestedItemsDuringDeletion: true,
				},
			},
		},
		{
			Name: "Delete Nested Items During Deletion Enabled",
			Input: []interface{}{
				map[string]interface{}{
					"template_deployment": []interface{}{
						map[string]interface{}{
							"delete_nested_items_during_deletion": true,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				TemplateDeployment: features.TemplateDeploymentFeatures{
					DeleteNestedItemsDuringDeletion: true,
				},
			},
		},
		{
			Name: "Delete Nested Items During Deletion Disabled",
			Input: []interface{}{
				map[string]interface{}{
					"template_deployment": []interface{}{
						map[string]interface{}{
							"delete_nested_items_during_deletion": false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				TemplateDeployment: features.TemplateDeploymentFeatures{
					DeleteNestedItemsDuringDeletion: false,
				},
			},
		},
	}

	for _, testCase := range testData {
		t.Logf("[DEBUG] Test Case: %q", testCase.Name)
		result := expandFeatures(testCase.Input)
		if !reflect.DeepEqual(result.TemplateDeployment, testCase.Expected.TemplateDeployment) {
			t.Fatalf("Expected %+v but got %+v", result.TemplateDeployment, testCase.Expected.TemplateDeployment)
		}
	}
}

func TestExpandFeaturesVirtualMachine(t *testing.T) {
	testData := []struct {
		Name     string
		Input    []interface{}
		EnvVars  map[string]interface{}
		Expected features.UserFeatures
	}{
		{
			Name: "Empty Block",
			Input: []interface{}{
				map[string]interface{}{
					"virtual_machine": []interface{}{},
				},
			},
			Expected: features.UserFeatures{
				VirtualMachine: features.VirtualMachineFeatures{
					DetachImplicitDataDiskOnDeletion: false,
					DeleteOSDiskOnDeletion:           true,
					SkipShutdownAndForceDelete:       false,
				},
			},
		},
		{
			Name: "Detach implicit Data Disk Enabled",
			Input: []interface{}{
				map[string]interface{}{
					"virtual_machine": []interface{}{
						map[string]interface{}{
							"detach_implicit_data_disk_on_deletion": true,
							"delete_os_disk_on_deletion":            false,
							"force_delete":                          false,
							"shutdown_before_deletion":              false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				VirtualMachine: features.VirtualMachineFeatures{
					DetachImplicitDataDiskOnDeletion: true,
					DeleteOSDiskOnDeletion:           false,
					SkipShutdownAndForceDelete:       false,
				},
			},
		},
		{
			Name: "Delete OS Disk Enabled",
			Input: []interface{}{
				map[string]interface{}{
					"virtual_machine": []interface{}{
						map[string]interface{}{
							"detach_implicit_data_disk_on_deletion": false,
							"delete_os_disk_on_deletion":            true,
							"force_delete":                          false,
							"shutdown_before_deletion":              false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				VirtualMachine: features.VirtualMachineFeatures{
					DetachImplicitDataDiskOnDeletion: false,
					DeleteOSDiskOnDeletion:           true,
					SkipShutdownAndForceDelete:       false,
				},
			},
		},
		{
			Name: "Graceful Shutdown Enabled",
			Input: []interface{}{
				map[string]interface{}{
					"virtual_machine": []interface{}{
						map[string]interface{}{
							"detach_implicit_data_disk_on_deletion": false,
							"delete_os_disk_on_deletion":            false,
							"force_delete":                          false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				VirtualMachine: features.VirtualMachineFeatures{
					DetachImplicitDataDiskOnDeletion: false,
					DeleteOSDiskOnDeletion:           false,
					SkipShutdownAndForceDelete:       false,
				},
			},
		},
		{
			Name: "Skip Shutdown and Force Delete Enabled",
			Input: []interface{}{
				map[string]interface{}{
					"virtual_machine": []interface{}{
						map[string]interface{}{
							"detach_implicit_data_disk_on_deletion": false,
							"delete_os_disk_on_deletion":            false,
							"skip_shutdown_and_force_delete":        true,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				VirtualMachine: features.VirtualMachineFeatures{
					DetachImplicitDataDiskOnDeletion: false,
					DeleteOSDiskOnDeletion:           false,
					SkipShutdownAndForceDelete:       true,
				},
			},
		},
		{
			Name: "All Disabled",
			Input: []interface{}{
				map[string]interface{}{
					"virtual_machine": []interface{}{
						map[string]interface{}{
							"detach_implicit_data_disk_on_deletion": false,
							"delete_os_disk_on_deletion":            false,
							"skip_shutdown_and_force_delete":        false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				VirtualMachine: features.VirtualMachineFeatures{
					DetachImplicitDataDiskOnDeletion: false,
					DeleteOSDiskOnDeletion:           false,
					SkipShutdownAndForceDelete:       false,
				},
			},
		},
	}

	for _, testCase := range testData {
		t.Logf("[DEBUG] Test Case: %q", testCase.Name)
		result := expandFeatures(testCase.Input)
		if !reflect.DeepEqual(result.VirtualMachine, testCase.Expected.VirtualMachine) {
			t.Fatalf("Expected %+v but got %+v", result.VirtualMachine, testCase.Expected.VirtualMachine)
		}
	}
}

func TestExpandFeaturesVirtualMachineScaleSet(t *testing.T) {
	testData := []struct {
		Name     string
		Input    []interface{}
		EnvVars  map[string]interface{}
		Expected features.UserFeatures
	}{
		{
			Name: "Empty Block",
			Input: []interface{}{
				map[string]interface{}{
					"virtual_machine_scale_set": []interface{}{},
				},
			},
			Expected: features.UserFeatures{
				VirtualMachineScaleSet: features.VirtualMachineScaleSetFeatures{
					ReimageOnManualUpgrade:    true,
					RollInstancesWhenRequired: true,
					ScaleToZeroOnDelete:       true,
				},
			},
		},
		{
			Name: "Force Delete Enabled",
			Input: []interface{}{
				map[string]interface{}{
					"virtual_machine_scale_set": []interface{}{
						map[string]interface{}{
							"force_delete":                 true,
							"roll_instances_when_required": false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				VirtualMachineScaleSet: features.VirtualMachineScaleSetFeatures{
					ForceDelete:               true,
					ReimageOnManualUpgrade:    true,
					RollInstancesWhenRequired: false,
					ScaleToZeroOnDelete:       true,
				},
			},
		},
		{
			Name: "Roll Instances Enabled",
			Input: []interface{}{
				map[string]interface{}{
					"virtual_machine_scale_set": []interface{}{
						map[string]interface{}{
							"force_delete":                 false,
							"roll_instances_when_required": true,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				VirtualMachineScaleSet: features.VirtualMachineScaleSetFeatures{
					ForceDelete:               false,
					ReimageOnManualUpgrade:    true,
					RollInstancesWhenRequired: true,
					ScaleToZeroOnDelete:       true,
				},
			},
		},
		{
			Name: "Scale In On Delete Disabled",
			Input: []interface{}{
				map[string]interface{}{
					"virtual_machine_scale_set": []interface{}{
						map[string]interface{}{
							"force_delete":                  false,
							"roll_instances_when_required":  true,
							"scale_to_zero_before_deletion": false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				VirtualMachineScaleSet: features.VirtualMachineScaleSetFeatures{
					ForceDelete:               false,
					ReimageOnManualUpgrade:    true,
					RollInstancesWhenRequired: true,
					ScaleToZeroOnDelete:       false,
				},
			},
		},
		{
			Name: "All Fields Disabled",
			Input: []interface{}{
				map[string]interface{}{
					"virtual_machine_scale_set": []interface{}{
						map[string]interface{}{
							"force_delete":                  false,
							"reimage_on_manual_upgrade":     false,
							"roll_instances_when_required":  false,
							"scale_to_zero_before_deletion": false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				VirtualMachineScaleSet: features.VirtualMachineScaleSetFeatures{
					ForceDelete:               false,
					ReimageOnManualUpgrade:    false,
					RollInstancesWhenRequired: false,
					ScaleToZeroOnDelete:       false,
				},
			},
		},
	}

	for _, testCase := range testData {
		t.Logf("[DEBUG] Test Case: %q", testCase.Name)
		result := expandFeatures(testCase.Input)
		if !reflect.DeepEqual(result.VirtualMachineScaleSet, testCase.Expected.VirtualMachineScaleSet) {
			t.Fatalf("Expected %+v but got %+v", testCase.Expected.VirtualMachineScaleSet, result.VirtualMachineScaleSet)
		}
	}
}

func TestExpandFeaturesLogAnalyticsWorkspace(t *testing.T) {
	testData := []struct {
		Name     string
		Input    []interface{}
		EnvVars  map[string]interface{}
		Expected features.UserFeatures
	}{
		{
			Name: "Empty Block",
			Input: []interface{}{
				map[string]interface{}{
					"log_analytics_workspace": []interface{}{},
				},
			},
			Expected: features.UserFeatures{
				LogAnalyticsWorkspace: features.LogAnalyticsWorkspaceFeatures{
					PermanentlyDeleteOnDestroy: false,
				},
			},
		},
		{
			Name: "Permanent Delete Enabled",
			Input: []interface{}{
				map[string]interface{}{
					"log_analytics_workspace": []interface{}{
						map[string]interface{}{
							"permanently_delete_on_destroy": true,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				LogAnalyticsWorkspace: features.LogAnalyticsWorkspaceFeatures{
					PermanentlyDeleteOnDestroy: true,
				},
			},
		},
		{
			Name: "Permanent Delete Disabled",
			Input: []interface{}{
				map[string]interface{}{
					"log_analytics_workspace": []interface{}{
						map[string]interface{}{
							"permanently_delete_on_destroy": false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				LogAnalyticsWorkspace: features.LogAnalyticsWorkspaceFeatures{
					PermanentlyDeleteOnDestroy: false,
				},
			},
		},
	}
	for _, testCase := range testData {
		t.Logf("[DEBUG] Test Case: %q", testCase.Name)
		result := expandFeatures(testCase.Input)
		if !reflect.DeepEqual(result.LogAnalyticsWorkspace, testCase.Expected.LogAnalyticsWorkspace) {
			t.Fatalf("Expected %+v but got %+v", result.LogAnalyticsWorkspace, testCase.Expected.LogAnalyticsWorkspace)
		}
	}
}

func TestExpandFeaturesResourceGroup(t *testing.T) {
	testData := []struct {
		Name     string
		Input    []interface{}
		EnvVars  map[string]interface{}
		Expected features.UserFeatures
	}{
		{
			Name: "Empty Block",
			Input: []interface{}{
				map[string]interface{}{
					"resource_group": []interface{}{},
				},
			},
			Expected: features.UserFeatures{
				ResourceGroup: features.ResourceGroupFeatures{
					PreventDeletionIfContainsResources: true,
				},
			},
		},
		{
			Name: "Prevent Deletion If Contains Resources Enabled",
			Input: []interface{}{
				map[string]interface{}{
					"resource_group": []interface{}{
						map[string]interface{}{
							"prevent_deletion_if_contains_resources": true,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				ResourceGroup: features.ResourceGroupFeatures{
					PreventDeletionIfContainsResources: true,
				},
			},
		},
		{
			Name: "Prevent Deletion If Contains Resources Disabled",
			Input: []interface{}{
				map[string]interface{}{
					"resource_group": []interface{}{
						map[string]interface{}{
							"prevent_deletion_if_contains_resources": false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				ResourceGroup: features.ResourceGroupFeatures{
					PreventDeletionIfContainsResources: false,
				},
			},
		},
	}

	for _, testCase := range testData {
		t.Logf("[DEBUG] Test Case: %q", testCase.Name)
		result := expandFeatures(testCase.Input)
		if !reflect.DeepEqual(result.ResourceGroup, testCase.Expected.ResourceGroup) {
			t.Fatalf("Expected %+v but got %+v", result.ResourceGroup, testCase.Expected.ResourceGroup)
		}
	}
}

func TestExpandFeaturesRecoveryServicesVault(t *testing.T) {
	testData := []struct {
		Name     string
		Input    []interface{}
		EnvVars  map[string]interface{}
		Expected features.UserFeatures
	}{
		{
			Name: "Empty Block",
			Input: []interface{}{
				map[string]interface{}{
					"recovery_services_vaults": []interface{}{},
				},
			},
			Expected: features.UserFeatures{
				RecoveryServicesVault: features.RecoveryServicesVault{
					RecoverSoftDeletedBackupProtectedVM: true,
				},
			},
		},
		{
			Name: "Recover Soft Deleted Protected VM Enabled",
			Input: []interface{}{
				map[string]interface{}{
					"recovery_services_vaults": []interface{}{
						map[string]interface{}{
							"recover_soft_deleted_backup_protected_vm": true,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				RecoveryServicesVault: features.RecoveryServicesVault{
					RecoverSoftDeletedBackupProtectedVM: true,
				},
			},
		},
		{
			Name: "Recover Soft Deleted Protected VM Disabled",
			Input: []interface{}{
				map[string]interface{}{
					"recovery_services_vaults": []interface{}{
						map[string]interface{}{
							"recover_soft_deleted_backup_protected_vm": false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				RecoveryServicesVault: features.RecoveryServicesVault{
					RecoverSoftDeletedBackupProtectedVM: false,
				},
			},
		},
	}

	for _, testCase := range testData {
		t.Logf("[DEBUG] Test Case: %q", testCase.Name)
		result := expandFeatures(testCase.Input)
		if !reflect.DeepEqual(result.RecoveryServicesVault, testCase.Expected.RecoveryServicesVault) {
			t.Fatalf("Expected %+v but got %+v", testCase.Expected.RecoveryServicesVault, result.RecoveryServicesVault)
		}
	}
}

func TestExpandFeaturesManagedDisk(t *testing.T) {
	testData := []struct {
		Name     string
		Input    []interface{}
		EnvVars  map[string]interface{}
		Expected features.UserFeatures
	}{
		{
			Name: "Empty Block",
			Input: []interface{}{
				map[string]interface{}{
					"managed_disk": []interface{}{},
				},
			},
			Expected: features.UserFeatures{
				ManagedDisk: features.ManagedDiskFeatures{
					ExpandWithoutDowntime: true,
				},
			},
		},
		{
			Name: "No Downtime Resize Enabled",
			Input: []interface{}{
				map[string]interface{}{
					"managed_disk": []interface{}{
						map[string]interface{}{
							"expand_without_downtime": true,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				ManagedDisk: features.ManagedDiskFeatures{
					ExpandWithoutDowntime: true,
				},
			},
		},
		{
			Name: "No Downtime Resize Disabled",
			Input: []interface{}{
				map[string]interface{}{
					"managed_disk": []interface{}{
						map[string]interface{}{
							"expand_without_downtime": false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				ManagedDisk: features.ManagedDiskFeatures{
					ExpandWithoutDowntime: false,
				},
			},
		},
	}

	for _, testCase := range testData {
		t.Logf("[DEBUG] Test Case: %q", testCase.Name)
		result := expandFeatures(testCase.Input)
		if !reflect.DeepEqual(result.ManagedDisk, testCase.Expected.ManagedDisk) {
			t.Fatalf("Expected %+v but got %+v", result.ManagedDisk, testCase.Expected.ManagedDisk)
		}
	}
}

func TestExpandFeaturesStorage(t *testing.T) {
	testData := []struct {
		Name     string
		Input    []interface{}
		EnvVars  map[string]interface{}
		Expected features.UserFeatures
	}{
		{
			Name: "Empty Block",
			Input: []interface{}{
				map[string]interface{}{
					"storage": []interface{}{},
				},
			},
			Expected: features.UserFeatures{
				Storage: features.StorageFeatures{
					DataPlaneAvailable: true,
				},
			},
		},
		{
			Name: "Storage Data Plane on Create is Disabled",
			Input: []interface{}{
				map[string]interface{}{
					"storage": []interface{}{
						map[string]interface{}{
							"data_plane_available": false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				Storage: features.StorageFeatures{
					DataPlaneAvailable: false,
				},
			},
		},
	}

	for _, testCase := range testData {
		t.Logf("[DEBUG] Test Case: %q", testCase.Name)
		result := expandFeatures(testCase.Input)
		if !reflect.DeepEqual(result.Storage, testCase.Expected.Storage) {
			t.Fatalf("Expected %+v but got %+v", result.Storage, testCase.Expected.Storage)
		}
	}
}

func TestExpandFeaturesSubscription(t *testing.T) {
	testData := []struct {
		Name     string
		Input    []interface{}
		EnvVars  map[string]interface{}
		Expected features.UserFeatures
	}{
		{
			Name: "Empty Block",
			Input: []interface{}{
				map[string]interface{}{
					"subscription": []interface{}{},
				},
			},
			Expected: features.UserFeatures{
				Subscription: features.SubscriptionFeatures{
					PreventCancellationOnDestroy: false,
				},
			},
		},
		{
			Name: "Subscription cancellation on destroy is Disabled",
			Input: []interface{}{
				map[string]interface{}{
					"subscription": []interface{}{
						map[string]interface{}{
							"prevent_cancellation_on_destroy": true,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				Subscription: features.SubscriptionFeatures{
					PreventCancellationOnDestroy: true,
				},
			},
		},
	}

	for _, testCase := range testData {
		t.Logf("[DEBUG] Test Case: %q", testCase.Name)
		result := expandFeatures(testCase.Input)
		if !reflect.DeepEqual(result.Subscription, testCase.Expected.Subscription) {
			t.Fatalf("Expected %+v but got %+v", result.Subscription, testCase.Expected.Subscription)
		}
	}
}

func TestExpandFeaturesPosgresqlFlexibleServer(t *testing.T) {
	testData := []struct {
		Name     string
		Input    []interface{}
		EnvVars  map[string]interface{}
		Expected features.UserFeatures
	}{
		{
			Name: "Empty Block",
			Input: []interface{}{
				map[string]interface{}{
					"postgresql_flexible_server": []interface{}{},
				},
			},
			Expected: features.UserFeatures{
				PostgresqlFlexibleServer: features.PostgresqlFlexibleServerFeatures{
					RestartServerOnConfigurationValueChange: true,
				},
			},
		},
		{
			Name: "Postgresql Flexible Server restarts on configuration change is Enabled",
			Input: []interface{}{
				map[string]interface{}{
					"postgresql_flexible_server": []interface{}{
						map[string]interface{}{
							"restart_server_on_configuration_value_change": true,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				PostgresqlFlexibleServer: features.PostgresqlFlexibleServerFeatures{
					RestartServerOnConfigurationValueChange: true,
				},
			},
		},
		{
			Name: "Postgresql Flexible Server restarts on configuration change is Disabled",
			Input: []interface{}{
				map[string]interface{}{
					"postgresql_flexible_server": []interface{}{
						map[string]interface{}{
							"restart_server_on_configuration_value_change": false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				PostgresqlFlexibleServer: features.PostgresqlFlexibleServerFeatures{
					RestartServerOnConfigurationValueChange: false,
				},
			},
		},
	}

	for _, testCase := range testData {
		t.Logf("[DEBUG] Test Case: %q", testCase.Name)
		result := expandFeatures(testCase.Input)
		if !reflect.DeepEqual(result.Subscription, testCase.Expected.Subscription) {
			t.Fatalf("Expected %+v but got %+v", result.Subscription, testCase.Expected.Subscription)
		}
	}
}

func TestExpandFeaturesMachineLearning(t *testing.T) {
	testData := []struct {
		Name     string
		Input    []interface{}
		EnvVars  map[string]interface{}
		Expected features.UserFeatures
	}{
		{
			Name: "Empty Block",
			Input: []interface{}{
				map[string]interface{}{
					"machine_learning": []interface{}{},
				},
			},
			Expected: features.UserFeatures{
				MachineLearning: features.MachineLearningFeatures{
					PurgeSoftDeletedWorkspaceOnDestroy: false,
				},
			},
		},
		{
			Name: "MachineLearning Workspace purge soft delete on destroy",
			Input: []interface{}{
				map[string]interface{}{
					"machine_learning": []interface{}{
						map[string]interface{}{
							"purge_soft_deleted_workspace_on_destroy": true,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				MachineLearning: features.MachineLearningFeatures{
					PurgeSoftDeletedWorkspaceOnDestroy: true,
				},
			},
		},
		{
			Name: "MachineLearning Workspace does not purge soft delete on destroy",
			Input: []interface{}{
				map[string]interface{}{
					"machine_learning": []interface{}{
						map[string]interface{}{
							"purge_soft_deleted_workspace_on_destroy": false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				MachineLearning: features.MachineLearningFeatures{
					PurgeSoftDeletedWorkspaceOnDestroy: false,
				},
			},
		},
	}

	for _, testCase := range testData {
		t.Logf("[DEBUG] Test Case: %q", testCase.Name)
		result := expandFeatures(testCase.Input)
		if !reflect.DeepEqual(result.Subscription, testCase.Expected.Subscription) {
			t.Fatalf("Expected %+v but got %+v", result.Subscription, testCase.Expected.Subscription)
		}
	}
}

func TestExpandFeaturesRecoveryService(t *testing.T) {
	testData := []struct {
		Name     string
		Input    []interface{}
		EnvVars  map[string]interface{}
		Expected features.UserFeatures
	}{
		{
			Name: "Empty Block",
			Input: []interface{}{
				map[string]interface{}{
					"recovery_service": []interface{}{},
				},
			},
			Expected: features.UserFeatures{
				RecoveryService: features.RecoveryServiceFeatures{
					VMBackupStopProtectionAndRetainDataOnDestroy:    false,
					VMBackupSuspendProtectionAndRetainDataOnDestroy: false,
					PurgeProtectedItemsFromVaultOnDestroy:           false,
				},
			},
		},
		{
			Name: "Recovery Service Features Enabled",
			Input: []interface{}{
				map[string]interface{}{
					"recovery_service": []interface{}{
						map[string]interface{}{
							"vm_backup_stop_protection_and_retain_data_on_destroy":    true,
							"vm_backup_suspend_protection_and_retain_data_on_destroy": true,
							"purge_protected_items_from_vault_on_destroy":             true,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				RecoveryService: features.RecoveryServiceFeatures{
					VMBackupStopProtectionAndRetainDataOnDestroy:    true,
					VMBackupSuspendProtectionAndRetainDataOnDestroy: true,
					PurgeProtectedItemsFromVaultOnDestroy:           true,
				},
			},
		},
		{
			Name: "Recovery Service Features Disabled",
			Input: []interface{}{
				map[string]interface{}{
					"recovery_service": []interface{}{
						map[string]interface{}{
							"vm_backup_stop_protection_and_retain_data_on_destroy":    false,
							"vm_backup_suspend_protection_and_retain_data_on_destroy": false,
							"purge_protected_items_from_vault_on_destroy":             false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				RecoveryService: features.RecoveryServiceFeatures{
					VMBackupStopProtectionAndRetainDataOnDestroy:    false,
					VMBackupSuspendProtectionAndRetainDataOnDestroy: false,
					PurgeProtectedItemsFromVaultOnDestroy:           false,
				},
			},
		},
	}

	for _, testCase := range testData {
		t.Logf("[DEBUG] Test Case: %q", testCase.Name)
		result := expandFeatures(testCase.Input)
		if !reflect.DeepEqual(result.Subscription, testCase.Expected.Subscription) {
			t.Fatalf("Expected %+v but got %+v", result.Subscription, testCase.Expected.Subscription)
		}
	}
}

func TestExpandFeaturesNetApp(t *testing.T) {
	testData := []struct {
		Name     string
		Input    []interface{}
		EnvVars  map[string]interface{}
		Expected features.UserFeatures
	}{
		{
			Name: "Empty Block",
			Input: []interface{}{
				map[string]interface{}{
					"netapp": []interface{}{},
				},
			},
			Expected: features.UserFeatures{
				NetApp: features.NetAppFeatures{
					DeleteBackupsOnBackupVaultDestroy: false,
					PreventVolumeDestruction:          true,
				},
			},
		},
		{
			Name: "NetApp Features Enabled",
			Input: []interface{}{
				map[string]interface{}{
					"netapp": []interface{}{
						map[string]interface{}{
							"delete_backups_on_backup_vault_destroy": true,
							"prevent_volume_destruction":             true,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				NetApp: features.NetAppFeatures{
					DeleteBackupsOnBackupVaultDestroy: false,
					PreventVolumeDestruction:          true,
				},
			},
		},
		{
			Name: "NetApp Features Disabled",
			Input: []interface{}{
				map[string]interface{}{
					"netapp": []interface{}{
						map[string]interface{}{
							"delete_backups_on_backup_vault_destroy": false,
							"prevent_volume_destruction":             false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				NetApp: features.NetAppFeatures{
					DeleteBackupsOnBackupVaultDestroy: false,
					PreventVolumeDestruction:          false,
				},
			},
		},
		{
			Name: "NetApp Features Reverse Default Values",
			Input: []interface{}{
				map[string]interface{}{
					"netapp": []interface{}{
						map[string]interface{}{
							"delete_backups_on_backup_vault_destroy": true,
							"prevent_volume_destruction":             false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				NetApp: features.NetAppFeatures{
					DeleteBackupsOnBackupVaultDestroy: true,
					PreventVolumeDestruction:          false,
				},
			},
		},
	}

	for _, testCase := range testData {
		t.Logf("[DEBUG] Test Case: %q", testCase.Name)
		result := expandFeatures(testCase.Input)
		if !reflect.DeepEqual(result.Subscription, testCase.Expected.Subscription) {
			t.Fatalf("Expected %+v but got %+v", result.Subscription, testCase.Expected.Subscription)
		}
	}
}

func TestExpandFeaturesDatabricksWorkspace(t *testing.T) {
	testData := []struct {
		Name     string
		Input    []interface{}
		EnvVars  map[string]interface{}
		Expected features.UserFeatures
	}{
		{
			Name: "Empty Block",
			Input: []interface{}{
				map[string]interface{}{
					"databricks_workspace": []interface{}{},
				},
			},
			Expected: features.UserFeatures{
				DatabricksWorkspace: features.DatabricksWorkspaceFeatures{
					ForceDelete: false,
				},
			},
		},
		{
			Name: "Databricks Workspace Features Enabled",
			Input: []interface{}{
				map[string]interface{}{
					"databricks_workspace": []interface{}{
						map[string]interface{}{
							"force_delete": true,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				DatabricksWorkspace: features.DatabricksWorkspaceFeatures{
					ForceDelete: true,
				},
			},
		},
		{
			Name: "Databricks Workspace Features Disabled",
			Input: []interface{}{
				map[string]interface{}{
					"databricks_workspace": []interface{}{
						map[string]interface{}{
							"force_delete": false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				DatabricksWorkspace: features.DatabricksWorkspaceFeatures{
					ForceDelete: false,
				},
			},
		},
	}

	for _, testCase := range testData {
		t.Logf("[DEBUG] Test Case: %q", testCase.Name)
		result := expandFeatures(testCase.Input)
		if !reflect.DeepEqual(result.DatabricksWorkspace, testCase.Expected.DatabricksWorkspace) {
			t.Fatalf("Expected %+v but got %+v", result.DatabricksWorkspace, testCase.Expected.DatabricksWorkspace)
		}
	}
}
