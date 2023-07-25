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
					RecoverSoftDeletedCerts:          true,
					RecoverSoftDeletedKeys:           true,
					RecoverSoftDeletedKeyVaults:      true,
					RecoverSoftDeletedSecrets:        true,
				},
				LogAnalyticsWorkspace: features.LogAnalyticsWorkspaceFeatures{
					PermanentlyDeleteOnDestroy: true,
				},
				ManagedDisk: features.ManagedDiskFeatures{
					ExpandWithoutDowntime: true,
				},
				TemplateDeployment: features.TemplateDeploymentFeatures{
					DeleteNestedItemsDuringDeletion: true,
				},
				VirtualMachine: features.VirtualMachineFeatures{
					DeleteOSDiskOnDeletion:     true,
					GracefulShutdown:           false,
					SkipShutdownAndForceDelete: false,
				},
				VirtualMachineScaleSet: features.VirtualMachineScaleSetFeatures{
					ForceDelete:               false,
					RollInstancesWhenRequired: true,
					ScaleToZeroOnDelete:       true,
				},
				ResourceGroup: features.ResourceGroupFeatures{
					PreventDeletionIfContainsResources: true,
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
							"purge_soft_deleted_certificates_on_destroy":              true,
							"purge_soft_deleted_keys_on_destroy":                      true,
							"purge_soft_deleted_secrets_on_destroy":                   true,
							"purge_soft_deleted_hardware_security_modules_on_destroy": true,
							"purge_soft_delete_on_destroy":                            true,
							"recover_soft_deleted_certificates":                       true,
							"recover_soft_deleted_keys":                               true,
							"recover_soft_deleted_key_vaults":                         true,
							"recover_soft_deleted_secrets":                            true,
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
					"resource_group": []interface{}{
						map[string]interface{}{
							"prevent_deletion_if_contains_resources": true,
						},
					},
					"template_deployment": []interface{}{
						map[string]interface{}{
							"delete_nested_items_during_deletion": true,
						},
					},
					"virtual_machine": []interface{}{
						map[string]interface{}{
							"delete_os_disk_on_deletion":     true,
							"graceful_shutdown":              true,
							"skip_shutdown_and_force_delete": true,
						},
					},
					"virtual_machine_scale_set": []interface{}{
						map[string]interface{}{
							"roll_instances_when_required":  true,
							"force_delete":                  true,
							"scale_to_zero_before_deletion": true,
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
					RecoverSoftDeletedCerts:          true,
					RecoverSoftDeletedKeys:           true,
					RecoverSoftDeletedKeyVaults:      true,
					RecoverSoftDeletedSecrets:        true,
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
				TemplateDeployment: features.TemplateDeploymentFeatures{
					DeleteNestedItemsDuringDeletion: true,
				},
				VirtualMachine: features.VirtualMachineFeatures{
					DeleteOSDiskOnDeletion:     true,
					GracefulShutdown:           true,
					SkipShutdownAndForceDelete: true,
				},
				VirtualMachineScaleSet: features.VirtualMachineScaleSetFeatures{
					RollInstancesWhenRequired: true,
					ForceDelete:               true,
					ScaleToZeroOnDelete:       true,
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
							"purge_soft_deleted_certificates_on_destroy":              false,
							"purge_soft_deleted_keys_on_destroy":                      false,
							"purge_soft_deleted_secrets_on_destroy":                   false,
							"purge_soft_deleted_hardware_security_modules_on_destroy": false,
							"purge_soft_delete_on_destroy":                            false,
							"recover_soft_deleted_certificates":                       false,
							"recover_soft_deleted_keys":                               false,
							"recover_soft_deleted_key_vaults":                         false,
							"recover_soft_deleted_secrets":                            false,
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
					"resource_group": []interface{}{
						map[string]interface{}{
							"prevent_deletion_if_contains_resources": false,
						},
					},
					"template_deployment": []interface{}{
						map[string]interface{}{
							"delete_nested_items_during_deletion": false,
						},
					},
					"virtual_machine": []interface{}{
						map[string]interface{}{
							"delete_os_disk_on_deletion":     false,
							"graceful_shutdown":              false,
							"skip_shutdown_and_force_delete": false,
						},
					},
					"virtual_machine_scale_set": []interface{}{
						map[string]interface{}{
							"force_delete":                  false,
							"roll_instances_when_required":  false,
							"scale_to_zero_before_deletion": false,
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
					PurgeSoftDeleteOnDestroy:         false,
					RecoverSoftDeletedCerts:          false,
					RecoverSoftDeletedKeys:           false,
					RecoverSoftDeletedKeyVaults:      false,
					RecoverSoftDeletedSecrets:        false,
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
				TemplateDeployment: features.TemplateDeploymentFeatures{
					DeleteNestedItemsDuringDeletion: false,
				},
				VirtualMachine: features.VirtualMachineFeatures{
					DeleteOSDiskOnDeletion:     false,
					GracefulShutdown:           false,
					SkipShutdownAndForceDelete: false,
				},
				VirtualMachineScaleSet: features.VirtualMachineScaleSetFeatures{
					ForceDelete:               false,
					RollInstancesWhenRequired: false,
					ScaleToZeroOnDelete:       false,
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
					RecoverSoftDeletedCerts:          true,
					RecoverSoftDeletedKeys:           true,
					RecoverSoftDeletedKeyVaults:      true,
					RecoverSoftDeletedSecrets:        true,
				},
			},
		},
		{
			Name: "Purge Soft Delete On Destroy and Recover Soft Deleted Key Vaults Enabled",
			Input: []interface{}{
				map[string]interface{}{
					"key_vault": []interface{}{
						map[string]interface{}{
							"purge_soft_deleted_certificates_on_destroy":              true,
							"purge_soft_deleted_keys_on_destroy":                      true,
							"purge_soft_deleted_secrets_on_destroy":                   true,
							"purge_soft_deleted_hardware_security_modules_on_destroy": true,
							"purge_soft_delete_on_destroy":                            true,
							"recover_soft_deleted_certificates":                       true,
							"recover_soft_deleted_keys":                               true,
							"recover_soft_deleted_key_vaults":                         true,
							"recover_soft_deleted_secrets":                            true,
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
					PurgeSoftDeleteOnDestroy:         true,
					RecoverSoftDeletedCerts:          true,
					RecoverSoftDeletedKeys:           true,
					RecoverSoftDeletedKeyVaults:      true,
					RecoverSoftDeletedSecrets:        true,
				},
			},
		},
		{
			Name: "Purge Soft Delete On Destroy and Recover Soft Deleted Key Vaults Disabled",
			Input: []interface{}{
				map[string]interface{}{
					"key_vault": []interface{}{
						map[string]interface{}{
							"purge_soft_deleted_certificates_on_destroy":              false,
							"purge_soft_deleted_keys_on_destroy":                      false,
							"purge_soft_deleted_secrets_on_destroy":                   false,
							"purge_soft_deleted_hardware_security_modules_on_destroy": false,
							"purge_soft_delete_on_destroy":                            false,
							"recover_soft_deleted_certificates":                       false,
							"recover_soft_deleted_keys":                               false,
							"recover_soft_deleted_key_vaults":                         false,
							"recover_soft_deleted_secrets":                            false,
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
					RecoverSoftDeletedCerts:          false,
					RecoverSoftDeletedKeyVaults:      false,
					RecoverSoftDeletedKeys:           false,
					RecoverSoftDeletedSecrets:        false,
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
					DeleteOSDiskOnDeletion:     true,
					GracefulShutdown:           false,
					SkipShutdownAndForceDelete: false,
				},
			},
		},
		{
			Name: "Delete OS Disk Enabled",
			Input: []interface{}{
				map[string]interface{}{
					"virtual_machine": []interface{}{
						map[string]interface{}{
							"delete_os_disk_on_deletion": true,
							"graceful_shutdown":          false,
							"force_delete":               false,
							"shutdown_before_deletion":   false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				VirtualMachine: features.VirtualMachineFeatures{
					DeleteOSDiskOnDeletion:     true,
					GracefulShutdown:           false,
					SkipShutdownAndForceDelete: false,
				},
			},
		},
		{
			Name: "Graceful Shutdown Enabled",
			Input: []interface{}{
				map[string]interface{}{
					"virtual_machine": []interface{}{
						map[string]interface{}{
							"delete_os_disk_on_deletion": false,
							"graceful_shutdown":          true,
							"force_delete":               false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				VirtualMachine: features.VirtualMachineFeatures{
					DeleteOSDiskOnDeletion:     false,
					GracefulShutdown:           true,
					SkipShutdownAndForceDelete: false,
				},
			},
		},
		{
			Name: "Skip Shutdown and Force Delete Enabled",
			Input: []interface{}{
				map[string]interface{}{
					"virtual_machine": []interface{}{
						map[string]interface{}{
							"delete_os_disk_on_deletion":     false,
							"graceful_shutdown":              false,
							"skip_shutdown_and_force_delete": true,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				VirtualMachine: features.VirtualMachineFeatures{
					DeleteOSDiskOnDeletion:     false,
					GracefulShutdown:           false,
					SkipShutdownAndForceDelete: true,
				},
			},
		},
		{
			Name: "All Disabled",
			Input: []interface{}{
				map[string]interface{}{
					"virtual_machine": []interface{}{
						map[string]interface{}{
							"delete_os_disk_on_deletion":     false,
							"graceful_shutdown":              false,
							"skip_shutdown_and_force_delete": false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				VirtualMachine: features.VirtualMachineFeatures{
					DeleteOSDiskOnDeletion:     false,
					GracefulShutdown:           false,
					SkipShutdownAndForceDelete: false,
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
							"roll_instances_when_required":  false,
							"scale_to_zero_before_deletion": false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				VirtualMachineScaleSet: features.VirtualMachineScaleSetFeatures{
					ForceDelete:               false,
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
					PermanentlyDeleteOnDestroy: !features.FourPointOhBeta(),
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
