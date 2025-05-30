// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package framework

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var testConfig = ProviderConfig{}

func TestProviderConfig_LoadDefault(t *testing.T) {
	if os.Getenv("ARM_CLIENT_ID") == "" {
		t.Skip("ARM_CLIENT_ID env var not set")
	}

	if os.Getenv("ARM_CLIENT_SECRET") == "" {
		t.Skip("ARM_CLIENT_SECRET env var not set")
	}

	// Skip enhanced validation
	t.Setenv("ARM_PROVIDER_ENHANCED_VALIDATION", "false")

	testData := &ProviderModel{
		ResourceProviderRegistrations: types.StringValue("none"),
		Features:                      defaultFeaturesList(),
	}

	testConfig.Load(context.Background(), testData, "unittest", &diag.Diagnostics{})

	if testConfig.Client == nil {
		t.Fatal("client nil after Load")
	}

	client := *testConfig.Client

	if account := client.Account; account != nil {
		if account.SubscriptionId == "" {
			t.Errorf("expected a value for subscription ID, but got an empty string")
		}
		if account.ClientId == "" {
			t.Errorf("expected a value for Client ID, but got an empty string")
		}
		if account.TenantId == "" {
			t.Errorf("expected a value for Tenant ID, but got an empty string")
		}
		if account.ObjectId == "" {
			t.Errorf("expected a value for Object ID, but got an empty string")
		}
		if account.Environment.Name != "Public" {
			t.Errorf("expected Environment name to be `Public` got %s", account.Environment.Name)
		}
		if suffix, _ := account.Environment.Storage.DomainSuffix(); suffix == nil || *suffix != "core.windows.net" {
			t.Errorf("expected `core.windows.net` got %+v", suffix)
		}
	} else {
		t.Error("account nil after Load")
	}

	features := client.Features

	if !features.ApiManagement.PurgeSoftDeleteOnDestroy {
		t.Errorf("expected api_management.purge_soft_delete_on_destroy to be true")
	}

	if !features.ApiManagement.RecoverSoftDeleted {
		t.Errorf("expected api_management.recover_soft_deleted to be true")
	}

	if !features.AppConfiguration.PurgeSoftDeleteOnDestroy {
		t.Errorf("expected app_configuration.purge_soft_delete_on_destroy to be true")
	}

	if !features.AppConfiguration.RecoverSoftDeleted {
		t.Errorf("expected app_configuration.recover_soft_deleted to be true")
	}

	if features.ApplicationInsights.DisableGeneratedRule {
		t.Errorf("expected application_insights.disable_generated_rule to be false")
	}

	if !features.CognitiveAccount.PurgeSoftDeleteOnDestroy {
		t.Errorf("expected cognitive_account.purge_soft_delete_on_destroy to be true")
	}

	if !features.KeyVault.PurgeSoftDeleteOnDestroy {
		t.Errorf("expected key_vault.purge_soft_delete_on_destroy to be true")
	}

	if !features.KeyVault.PurgeSoftDeletedCertsOnDestroy {
		t.Errorf("expected key_vault.purge_soft_deleted_certificates_on_destroy to be true")
	}

	if !features.KeyVault.PurgeSoftDeletedKeysOnDestroy {
		t.Errorf("expected key_vault.purge_soft_deleted_keys_on_destroy to be true")
	}

	if !features.KeyVault.PurgeSoftDeletedSecretsOnDestroy {
		t.Errorf("expected key_vault.purge_soft_deleted_secrets_on_destroy to be true")
	}

	if !features.KeyVault.PurgeSoftDeletedHSMsOnDestroy {
		t.Errorf("expected key_vault.purge_soft_deleted_hardware_security_modules_on_destroy to be true")
	}

	if !features.KeyVault.PurgeSoftDeletedHSMKeysOnDestroy {
		t.Errorf("expected key_vault.purge_soft_deleted_hardware_security_module_keys_on_destroy to be true")
	}

	if !features.KeyVault.RecoverSoftDeletedCerts {
		t.Errorf("expected key_vault.recover_soft_deleted_certificates to be true")
	}

	if !features.KeyVault.RecoverSoftDeletedKeyVaults {
		t.Errorf("expected key_vault.recover_soft_deleted_key_vaults to be true")
	}

	if !features.KeyVault.RecoverSoftDeletedKeys {
		t.Errorf("expected key_vault.recover_soft_deleted_keys to be true")
	}

	if !features.KeyVault.RecoverSoftDeletedSecrets {
		t.Errorf("expected key_vault.recover_soft_deleted_secrets to be true")
	}

	if !features.KeyVault.RecoverSoftDeletedHSMKeys {
		t.Errorf("expected key_vault.recover_soft_deleted_hsm_keys to be true")
	}

	if !features.LogAnalyticsWorkspace.PermanentlyDeleteOnDestroy {
		t.Errorf("expected log_analytics_workspace.permanently_delete_on_destroy to be true")
	}

	if features.TemplateDeployment.DeleteNestedItemsDuringDeletion {
		t.Errorf("expected template_deployment.delete_nested_items_during_deletion to be false")
	}

	if features.VirtualMachine.DeleteOSDiskOnDeletion {
		t.Errorf("expected virtual_machine.delete_os_disk_on_deletion to be false")
	}

	if features.VirtualMachine.DetachImplicitDataDiskOnDeletion {
		t.Errorf("expected virtual_machine.detach_implicit_data_disk_on_deletion to be false")
	}

	if features.VirtualMachine.SkipShutdownAndForceDelete {
		t.Errorf("expected virtual_machine.skip_shutdown_and_force_delete to be false")
	}

	if features.VirtualMachineScaleSet.ForceDelete {
		t.Errorf("expected virtual_machine.force_delete to be false")
	}

	if !features.VirtualMachineScaleSet.ReimageOnManualUpgrade {
		t.Errorf("expected virtual_machine.reimage_on_manual_upgrade to be true")
	}

	if !features.VirtualMachineScaleSet.RollInstancesWhenRequired {
		t.Errorf("expected virtual_machine.roll_instances_when_required to be true")
	}

	if features.VirtualMachineScaleSet.ScaleToZeroOnDelete {
		t.Errorf("expected virtual_machine.scale_to_zero_on_delete to be false")
	}

	if !features.ManagedDisk.ExpandWithoutDowntime {
		t.Errorf("expected managed_disk.expand_without_downtime to be true")
	}

	if features.Subscription.PreventCancellationOnDestroy {
		t.Errorf("expected subscription.prevent_cancellation_on_destroy to be false")
	}

	if !features.PostgresqlFlexibleServer.RestartServerOnConfigurationValueChange {
		t.Errorf("expected postgresql.restart_server_on_configuration_value_change to be true")
	}

	if features.MachineLearning.PurgeSoftDeletedWorkspaceOnDestroy {
		t.Errorf("expected machine_learning.PurgeSoftDeletedWorkspaceOnDestroy to be false")
	}

	if features.RecoveryService.VMBackupStopProtectionAndRetainDataOnDestroy {
		t.Errorf("expected recovery_service.vm_backup_stop_protection_and_retain_data_on_destroy to be false")
	}

	if features.RecoveryService.VMBackupSuspendProtectionAndRetainDataOnDestroy {
		t.Errorf("expected recovery_service.vm_backup_suspend_protection_and_retain_data_on_destroy to be false")
	}

	if features.RecoveryService.PurgeProtectedItemsFromVaultOnDestroy {
		t.Errorf("expected recovery_service.PurgeProtectedItemsFromVaultOnDestroy to be false")
	}

	if features.RecoveryService.PurgeProtectedItemsFromVaultOnDestroy {
		t.Errorf("expected recovery_service.PurgeProtectedItemsFromVaultOnDestroy to be false")
	}

	if features.NetApp.DeleteBackupsOnBackupVaultDestroy {
		t.Errorf("expected netapp.DeleteBackupsOnBackupVaultDestroy to be false")
	}

	if !features.NetApp.PreventVolumeDestruction {
		t.Errorf("expected netapp.PreventVolumeDestruction to be true")
	}

	if features.DatabricksWorkspace.ForceDelete {
		t.Errorf("expected databricks_workspace.ForceDelete to be false")
	}
}

// TODO - helper functions to make setting up test date more easily so we can add more configuration coverage

func defaultFeaturesList() types.List {
	apiManagement, _ := basetypes.NewObjectValueFrom(context.Background(), APIManagementAttributes, map[string]attr.Value{
		"purge_soft_delete_on_destroy": basetypes.NewBoolNull(),
		"recover_soft_deleted":         basetypes.NewBoolNull(),
	})
	apiManagementList, _ := basetypes.NewListValue(types.ObjectType{}.WithAttributeTypes(APIManagementAttributes), []attr.Value{apiManagement})

	appConfiguration, _ := basetypes.NewObjectValueFrom(context.Background(), AppConfigurationAttributes, map[string]attr.Value{
		"purge_soft_delete_on_destroy": basetypes.NewBoolNull(),
		"recover_soft_deleted":         basetypes.NewBoolNull(),
	})
	appConfigurationList, _ := basetypes.NewListValue(types.ObjectType{}.WithAttributeTypes(AppConfigurationAttributes), []attr.Value{appConfiguration})

	applicationInsights, _ := basetypes.NewObjectValueFrom(context.Background(), ApplicationInsightsAttributes, map[string]attr.Value{
		"disable_generated_rule": basetypes.NewBoolNull(),
	})
	applicationInsightsList, _ := basetypes.NewListValue(types.ObjectType{}.WithAttributeTypes(ApplicationInsightsAttributes), []attr.Value{applicationInsights})

	cognitiveAccount, _ := basetypes.NewObjectValueFrom(context.Background(), CognitiveAccountAttributes, map[string]attr.Value{
		"purge_soft_delete_on_destroy": basetypes.NewBoolNull(),
	})
	cognitiveAccountList, _ := basetypes.NewListValue(types.ObjectType{}.WithAttributeTypes(CognitiveAccountAttributes), []attr.Value{cognitiveAccount})

	keyVault, _ := basetypes.NewObjectValueFrom(context.Background(), KeyVaultAttributes, map[string]attr.Value{
		"purge_soft_delete_on_destroy":                            basetypes.NewBoolNull(),
		"purge_soft_deleted_certificates_on_destroy":              basetypes.NewBoolNull(),
		"purge_soft_deleted_keys_on_destroy":                      basetypes.NewBoolNull(),
		"purge_soft_deleted_secrets_on_destroy":                   basetypes.NewBoolNull(),
		"purge_soft_deleted_hardware_security_modules_on_destroy": basetypes.NewBoolNull(),
		"recover_soft_deleted_certificates":                       basetypes.NewBoolNull(),
		"recover_soft_deleted_key_vaults":                         basetypes.NewBoolNull(),
		"recover_soft_deleted_keys":                               basetypes.NewBoolNull(),
		"recover_soft_deleted_secrets":                            basetypes.NewBoolNull(),
	})
	keyVaultList, _ := basetypes.NewListValue(types.ObjectType{}.WithAttributeTypes(KeyVaultAttributes), []attr.Value{keyVault})

	logAnalyticsWorkspace, _ := basetypes.NewObjectValueFrom(context.Background(), LogAnalyticsWorkspaceAttributes, map[string]attr.Value{
		"permanently_delete_on_destroy": basetypes.NewBoolNull(),
	})
	logAnalyticsWorkspaceList, _ := basetypes.NewListValue(types.ObjectType{}.WithAttributeTypes(LogAnalyticsWorkspaceAttributes), []attr.Value{logAnalyticsWorkspace})

	templateDeployment, _ := basetypes.NewObjectValueFrom(context.Background(), TemplateDeploymentAttributes, map[string]attr.Value{
		"delete_nested_items_during_deletion": basetypes.NewBoolNull(),
	})
	templateDeploymentList, _ := basetypes.NewListValue(types.ObjectType{}.WithAttributeTypes(TemplateDeploymentAttributes), []attr.Value{templateDeployment})

	virtualMachine, _ := basetypes.NewObjectValueFrom(context.Background(), VirtualMachineAttributes, map[string]attr.Value{
		"delete_os_disk_on_deletion":     basetypes.NewBoolNull(),
		"skip_shutdown_and_force_delete": basetypes.NewBoolNull(),
	})
	virtualMachineList, _ := basetypes.NewListValue(types.ObjectType{}.WithAttributeTypes(VirtualMachineAttributes), []attr.Value{virtualMachine})

	virtualMachineScaleSet, _ := basetypes.NewObjectValueFrom(context.Background(), VirtualMachineScaleSetAttributes, map[string]attr.Value{
		"force_delete":                  basetypes.NewBoolNull(),
		"reimage_on_manual_upgrade":     basetypes.NewBoolNull(),
		"roll_instances_when_required":  basetypes.NewBoolNull(),
		"scale_to_zero_before_deletion": basetypes.NewBoolNull(),
	})
	virtualMachineScaleSetList, _ := basetypes.NewListValue(types.ObjectType{}.WithAttributeTypes(VirtualMachineScaleSetAttributes), []attr.Value{virtualMachineScaleSet})

	resourceGroup, _ := basetypes.NewObjectValueFrom(context.Background(), ResourceGroupAttributes, map[string]attr.Value{
		"prevent_deletion_if_contains_resources": basetypes.NewBoolNull(),
	})
	resourceGroupList, _ := basetypes.NewListValue(types.ObjectType{}.WithAttributeTypes(ResourceGroupAttributes), []attr.Value{resourceGroup})

	managedDisk, _ := basetypes.NewObjectValueFrom(context.Background(), ManagedDiskAttributes, map[string]attr.Value{
		"expand_without_downtime": basetypes.NewBoolNull(),
	})
	managedDiskList, _ := basetypes.NewListValue(types.ObjectType{}.WithAttributeTypes(ManagedDiskAttributes), []attr.Value{managedDisk})

	storage, _ := basetypes.NewObjectValueFrom(context.Background(), StorageAttributes, map[string]attr.Value{
		"data_plane_available": basetypes.NewBoolNull(),
	})
	storageList, _ := basetypes.NewListValue(types.ObjectType{}.WithAttributeTypes(StorageAttributes), []attr.Value{storage})

	subscription, _ := basetypes.NewObjectValueFrom(context.Background(), SubscriptionAttributes, map[string]attr.Value{
		"prevent_cancellation_on_destroy": basetypes.NewBoolNull(),
	})
	subscriptionList, _ := basetypes.NewListValue(types.ObjectType{}.WithAttributeTypes(SubscriptionAttributes), []attr.Value{subscription})

	postgresqlFlexibleServer, _ := basetypes.NewObjectValueFrom(context.Background(), PostgresqlFlexibleServerAttributes, map[string]attr.Value{
		"restart_server_on_configuration_value_change": basetypes.NewBoolNull(),
	})
	postgresqlFlexibleServerList, _ := basetypes.NewListValue(types.ObjectType{}.WithAttributeTypes(PostgresqlFlexibleServerAttributes), []attr.Value{postgresqlFlexibleServer})

	machineLearning, _ := basetypes.NewObjectValueFrom(context.Background(), MachineLearningAttributes, map[string]attr.Value{
		"purge_soft_deleted_workspace_on_destroy": basetypes.NewBoolNull(),
	})
	machineLearningList, _ := basetypes.NewListValue(types.ObjectType{}.WithAttributeTypes(MachineLearningAttributes), []attr.Value{machineLearning})

	recoveryServices, _ := basetypes.NewObjectValueFrom(context.Background(), RecoveryServiceAttributes, map[string]attr.Value{
		"vm_backup_stop_protection_and_retain_data_on_destroy":    basetypes.NewBoolNull(),
		"vm_backup_suspend_protection_and_retain_data_on_destroy": basetypes.NewBoolNull(),
		"purge_protected_items_from_vault_on_destroy":             basetypes.NewBoolNull(),
	})
	recoveryServicesList, _ := basetypes.NewListValue(types.ObjectType{}.WithAttributeTypes(RecoveryServiceAttributes), []attr.Value{recoveryServices})

	recoveryServicesVaults, _ := basetypes.NewObjectValueFrom(context.Background(), RecoveryServiceVaultsAttributes, map[string]attr.Value{
		"vm_backup_stop_protection_and_retain_data_on_destroy":    basetypes.NewBoolNull(),
		"vm_backup_suspend_protection_and_retain_data_on_destroy": basetypes.NewBoolNull(),
		"purge_protected_items_from_vault_on_destroy":             basetypes.NewBoolNull(),
	})
	recoveryServicesVaultsList, _ := basetypes.NewListValue(types.ObjectType{}.WithAttributeTypes(RecoveryServiceVaultsAttributes), []attr.Value{recoveryServicesVaults})

	netapp, _ := basetypes.NewObjectValueFrom(context.Background(), NetAppAttributes, map[string]attr.Value{
		"delete_backups_on_backup_vault_destroy": basetypes.NewBoolNull(),
		"prevent_volume_destruction":             basetypes.NewBoolNull(),
	})
	netappList, _ := basetypes.NewListValue(types.ObjectType{}.WithAttributeTypes(NetAppAttributes), []attr.Value{netapp})

	databricksWorkspace, _ := basetypes.NewObjectValueFrom(context.Background(), DatabricksWorkspaceAttributes, map[string]attr.Value{
		"force_delete": basetypes.NewBoolNull(),
	})
	databricksWorkspaceList, _ := basetypes.NewListValue(types.ObjectType{}.WithAttributeTypes(DatabricksWorkspaceAttributes), []attr.Value{databricksWorkspace})

	fData, d := basetypes.NewObjectValue(FeaturesAttributes, map[string]attr.Value{
		"api_management":             apiManagementList,
		"app_configuration":          appConfigurationList,
		"application_insights":       applicationInsightsList,
		"cognitive_account":          cognitiveAccountList,
		"key_vault":                  keyVaultList,
		"log_analytics_workspace":    logAnalyticsWorkspaceList,
		"template_deployment":        templateDeploymentList,
		"virtual_machine":            virtualMachineList,
		"virtual_machine_scale_set":  virtualMachineScaleSetList,
		"resource_group":             resourceGroupList,
		"managed_disk":               managedDiskList,
		"storage":                    storageList,
		"subscription":               subscriptionList,
		"postgresql_flexible_server": postgresqlFlexibleServerList,
		"machine_learning":           machineLearningList,
		"recovery_service":           recoveryServicesList,
		"recovery_services_vaults":   recoveryServicesVaultsList,
		"netapp":                     netappList,
		"databricks_workspace":       databricksWorkspaceList,
	})

	fmt.Printf("%+v", d)

	f, _ := basetypes.NewListValue(types.ObjectType{}.WithAttributeTypes(FeaturesAttributes), []attr.Value{fData})

	return f
}
