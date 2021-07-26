package provider

import (
	"reflect"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
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
				KeyVault: features.KeyVaultFeatures{
					PurgeSoftDeleteOnDestroy:    true,
					RecoverSoftDeletedKeyVaults: true,
				},
				Network: features.NetworkFeatures{
					RelaxedLocking: false,
				},
				TemplateDeployment: features.TemplateDeploymentFeatures{
					DeleteNestedItemsDuringDeletion: true,
				},
				VirtualMachine: features.VirtualMachineFeatures{
					DeleteOSDiskOnDeletion: true,
				},
				VirtualMachineScaleSet: features.VirtualMachineScaleSetFeatures{
					RollInstancesWhenRequired: true,
				},
				LogAnalyticsWorkspace: features.LogAnalyticsWorkspaceFeatures{
					PermanentlyDeleteOnDestroy: false,
				},
			},
		},
		{
			Name: "Complete Enabled",
			Input: []interface{}{
				map[string]interface{}{
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
					"template_deployment": []interface{}{
						map[string]interface{}{
							"delete_nested_items_during_deletion": true,
						},
					},
					"virtual_machine": []interface{}{
						map[string]interface{}{
							"delete_os_disk_on_deletion": true,
							"graceful_shutdown":          true,
						},
					},
					"virtual_machine_scale_set": []interface{}{
						map[string]interface{}{
							"roll_instances_when_required": true,
						},
					},
				},
			},
			Expected: features.UserFeatures{
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
				TemplateDeployment: features.TemplateDeploymentFeatures{
					DeleteNestedItemsDuringDeletion: true,
				},
				VirtualMachine: features.VirtualMachineFeatures{
					DeleteOSDiskOnDeletion: true,
					GracefulShutdown:       true,
				},
				VirtualMachineScaleSet: features.VirtualMachineScaleSetFeatures{
					RollInstancesWhenRequired: true,
				},
			},
		},
		{
			Name: "Complete Disabled",
			Input: []interface{}{
				map[string]interface{}{
					"virtual_machine": []interface{}{
						map[string]interface{}{
							"delete_os_disk_on_deletion": false,
							"graceful_shutdown":          false,
						},
					},
					"network_locking": []interface{}{
						map[string]interface{}{
							"relaxed_locking": false,
						},
					},
					"template_deployment": []interface{}{
						map[string]interface{}{
							"delete_nested_items_during_deletion": false,
						},
					},
					"virtual_machine_scale_set": []interface{}{
						map[string]interface{}{
							"roll_instances_when_required": false,
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
				},
			},
			Expected: features.UserFeatures{
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
				TemplateDeployment: features.TemplateDeploymentFeatures{
					DeleteNestedItemsDuringDeletion: false,
				},
				VirtualMachine: features.VirtualMachineFeatures{
					DeleteOSDiskOnDeletion: false,
					GracefulShutdown:       false,
				},
				VirtualMachineScaleSet: features.VirtualMachineScaleSetFeatures{
					RollInstancesWhenRequired: false,
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
					DeleteOSDiskOnDeletion: true,
					GracefulShutdown:       false,
				},
			},
		},
		{
			Name: "Delete OS Disk and Graceful Shutdown Enabled",
			Input: []interface{}{
				map[string]interface{}{
					"virtual_machine": []interface{}{
						map[string]interface{}{
							"delete_os_disk_on_deletion": true,
							"graceful_shutdown":          true,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				VirtualMachine: features.VirtualMachineFeatures{
					DeleteOSDiskOnDeletion: true,
					GracefulShutdown:       true,
				},
			},
		},
		{
			Name: "Delete OS Disk and Graceful Shutdown Disabled",
			Input: []interface{}{
				map[string]interface{}{
					"virtual_machine": []interface{}{
						map[string]interface{}{
							"delete_os_disk_on_deletion": false,
							"graceful_shutdown":          false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				VirtualMachine: features.VirtualMachineFeatures{
					DeleteOSDiskOnDeletion: false,
					GracefulShutdown:       false,
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
				},
			},
		},
		{
			Name: "Roll Instances Enabled",
			Input: []interface{}{
				map[string]interface{}{
					"virtual_machine_scale_set": []interface{}{
						map[string]interface{}{
							"roll_instances_when_required": true,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				VirtualMachineScaleSet: features.VirtualMachineScaleSetFeatures{
					RollInstancesWhenRequired: true,
				},
			},
		},
		{
			Name: "Roll Instances Disabled",
			Input: []interface{}{
				map[string]interface{}{
					"virtual_machine_scale_set": []interface{}{
						map[string]interface{}{
							"roll_instances_when_required": false,
						},
					},
				},
			},
			Expected: features.UserFeatures{
				VirtualMachineScaleSet: features.VirtualMachineScaleSetFeatures{
					RollInstancesWhenRequired: false,
				},
			},
		},
	}

	for _, testCase := range testData {
		t.Logf("[DEBUG] Test Case: %q", testCase.Name)
		result := expandFeatures(testCase.Input)
		if !reflect.DeepEqual(result.VirtualMachineScaleSet, testCase.Expected.VirtualMachineScaleSet) {
			t.Fatalf("Expected %+v but got %+v", result.VirtualMachineScaleSet, testCase.Expected.VirtualMachineScaleSet)
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
