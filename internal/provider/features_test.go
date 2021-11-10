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
					PurgeSoftDeleteOnDestroy: false,
				},
				CognitiveAccount: features.CognitiveAccountFeatures{
					PurgeSoftDeleteOnDestroy: true,
				},
				KeyVault: features.KeyVaultFeatures{
					PurgeSoftDeleteOnDestroy:    true,
					RecoverSoftDeletedKeyVaults: true,
				},
				LogAnalyticsWorkspace: features.LogAnalyticsWorkspaceFeatures{
					PermanentlyDeleteOnDestroy: false,
				},
				Network: features.NetworkFeatures{
					RelaxedLocking: false,
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
					PreventDeletionIfContainsResources: false,
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
						},
					},
					"cognitive_account": []interface{}{
						map[string]interface{}{
							"purge_soft_delete_on_destroy": true,
						},
					},
					"key_vault": []interface{}{
						map[string]interface{}{
							"purge_soft_delete_on_destroy":    true,
							"recover_soft_deleted_key_vaults": true,
						},
					},
					"log_analytics_workspace": []interface{}{
						map[string]interface{}{
							"permanently_delete_on_destroy": true,
						},
					},
					"network": []interface{}{
						map[string]interface{}{
							"relaxed_locking": true,
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
				},
				CognitiveAccount: features.CognitiveAccountFeatures{
					PurgeSoftDeleteOnDestroy: true,
				},
				KeyVault: features.KeyVaultFeatures{
					PurgeSoftDeleteOnDestroy:    true,
					RecoverSoftDeletedKeyVaults: true,
				},
				LogAnalyticsWorkspace: features.LogAnalyticsWorkspaceFeatures{
					PermanentlyDeleteOnDestroy: true,
				},
				Network: features.NetworkFeatures{
					RelaxedLocking: true,
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
						},
					},
					"cognitive_account": []interface{}{
						map[string]interface{}{
							"purge_soft_delete_on_destroy": false,
						},
					},
					"key_vault": []interface{}{
						map[string]interface{}{
							"purge_soft_delete_on_destroy":    false,
							"recover_soft_deleted_key_vaults": false,
						},
					},
					"log_analytics_workspace": []interface{}{
						map[string]interface{}{
							"permanently_delete_on_destroy": false,
						},
					},
					"network_locking": []interface{}{
						map[string]interface{}{
							"relaxed_locking": false,
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
				},
				CognitiveAccount: features.CognitiveAccountFeatures{
					PurgeSoftDeleteOnDestroy: false,
				},
				KeyVault: features.KeyVaultFeatures{
					PurgeSoftDeleteOnDestroy:    false,
					RecoverSoftDeletedKeyVaults: false,
				},
				LogAnalyticsWorkspace: features.LogAnalyticsWorkspaceFeatures{
					PermanentlyDeleteOnDestroy: false,
				},
				Network: features.NetworkFeatures{
					RelaxedLocking: false,
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
					PurgeSoftDeleteOnDestroy: false,
				},
			},
		},
		{
			Name: "Purge Soft Delete On Destroy Api Management Enabled",
			Input: []interface{}{
				map[string]interface{}{
					"api_management": []interface{}{
						map[string]interface{}{
							"purge_soft_delete_on_destroy": true,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				ApiManagement: features.ApiManagementFeatures{
					PurgeSoftDeleteOnDestroy: true,
				},
			},
		},
		{
			Name: "Purge Soft Delete On Destroy Api Management Disabled",
			Input: []interface{}{
				map[string]interface{}{
					"api_management": []interface{}{
						map[string]interface{}{
							"purge_soft_delete_on_destroy": false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				ApiManagement: features.ApiManagementFeatures{
					PurgeSoftDeleteOnDestroy: false,
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
					PurgeSoftDeleteOnDestroy:    true,
					RecoverSoftDeletedKeyVaults: true,
				},
			},
		},
		{
			Name: "Purge Soft Delete On Destroy and Recover Soft Deleted Key Vaults Enabled",
			Input: []interface{}{
				map[string]interface{}{
					"key_vault": []interface{}{
						map[string]interface{}{
							"purge_soft_delete_on_destroy":    true,
							"recover_soft_deleted_key_vaults": true,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				KeyVault: features.KeyVaultFeatures{
					PurgeSoftDeleteOnDestroy:    true,
					RecoverSoftDeletedKeyVaults: true,
				},
			},
		},
		{
			Name: "Purge Soft Delete On Destroy and Recover Soft Deleted Key Vaults Disabled",
			Input: []interface{}{
				map[string]interface{}{
					"key_vault": []interface{}{
						map[string]interface{}{
							"purge_soft_delete_on_destroy":    false,
							"recover_soft_deleted_key_vaults": false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				KeyVault: features.KeyVaultFeatures{
					PurgeSoftDeleteOnDestroy:    false,
					RecoverSoftDeletedKeyVaults: false,
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

func TestExpandFeaturesNetwork(t *testing.T) {
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
					"network": []interface{}{},
				},
			},
			Expected: features.UserFeatures{
				Network: features.NetworkFeatures{
					RelaxedLocking: false,
				},
			},
		},
		{
			Name: "Relaxed Locking Enabled",
			Input: []interface{}{
				map[string]interface{}{
					"network": []interface{}{
						map[string]interface{}{
							"relaxed_locking": true,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				Network: features.NetworkFeatures{
					RelaxedLocking: true,
				},
			},
		},
		{
			Name: "Relaxed Locking Disabled",
			Input: []interface{}{
				map[string]interface{}{
					"network": []interface{}{
						map[string]interface{}{
							"relaxed_locking": false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				Network: features.NetworkFeatures{
					RelaxedLocking: false,
				},
			},
		},
	}

	for _, testCase := range testData {
		t.Logf("[DEBUG] Test Case: %q", testCase.Name)
		result := expandFeatures(testCase.Input)
		if !reflect.DeepEqual(result.Network, testCase.Expected.Network) {
			t.Fatalf("Expected %+v but got %+v", result.Network, testCase.Expected.Network)
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
					PreventDeletionIfContainsResources: false,
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
