// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func schemaFeatures(supportLegacyTestSuite bool) *pluginsdk.Schema {
	// NOTE: if there's only one nested field these want to be Required (since there's no point
	//       specifying the block otherwise) - however for 2+ they should be optional
	featuresMap := map[string]*pluginsdk.Schema{
		// lintignore:XS003
		"api_management": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"purge_soft_delete_on_destroy": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"recover_soft_deleted": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
				},
			},
		},

		"app_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"purge_soft_delete_on_destroy": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},

					"recover_soft_deleted": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
				},
			},
		},

		"application_insights": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"disable_generated_rule": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"cognitive_account": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"purge_soft_delete_on_destroy": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
				},
			},
		},

		"key_vault": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"purge_soft_delete_on_destroy": {
						Description: "When enabled soft-deleted `azurerm_key_vault` resources will be permanently deleted (e.g purged), when destroyed",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"purge_soft_deleted_certificates_on_destroy": {
						Description: "When enabled soft-deleted `azurerm_key_vault_certificate` resources will be permanently deleted (e.g purged), when destroyed",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"purge_soft_deleted_keys_on_destroy": {
						Description: "When enabled soft-deleted `azurerm_key_vault_key` resources will be permanently deleted (e.g purged), when destroyed",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"purge_soft_deleted_secrets_on_destroy": {
						Description: "When enabled soft-deleted `azurerm_key_vault_secret` resources will be permanently deleted (e.g purged), when destroyed",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"purge_soft_deleted_hardware_security_modules_on_destroy": {
						Description: "When enabled soft-deleted `azurerm_key_vault_managed_hardware_security_module` resources will be permanently deleted (e.g purged), when destroyed",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"purge_soft_deleted_hardware_security_module_keys_on_destroy": {
						Description: "When enabled soft-deleted `azurerm_key_vault_managed_hardware_security_module_key` resources will be permanently deleted (e.g purged), when destroyed",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"recover_soft_deleted_certificates": {
						Description: "When enabled soft-deleted `azurerm_key_vault_certificate` resources will be restored, instead of creating new ones",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"recover_soft_deleted_key_vaults": {
						Description: "When enabled soft-deleted `azurerm_key_vault` resources will be restored, instead of creating new ones",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"recover_soft_deleted_keys": {
						Description: "When enabled soft-deleted `azurerm_key_vault_key` resources will be restored, instead of creating new ones",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"recover_soft_deleted_secrets": {
						Description: "When enabled soft-deleted `azurerm_key_vault_secret` resources will be restored, instead of creating new ones",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"recover_soft_deleted_hardware_security_module_keys": {
						Description: "When enabled soft-deleted `azurerm_key_vault_managed_hardware_security_module_key` resources will be restored, instead of creating new ones",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},
				},
			},
		},

		"log_analytics_workspace": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"permanently_delete_on_destroy": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"template_deployment": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"delete_nested_items_during_deletion": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
				},
			},
		},

		// lintignore:XS003
		"virtual_machine": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"detach_implicit_data_disk_on_deletion": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"delete_os_disk_on_deletion": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"skip_shutdown_and_force_delete": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"virtual_machine_scale_set": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"force_delete": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
					"reimage_on_manual_upgrade": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
					"roll_instances_when_required": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
					"scale_to_zero_before_deletion": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"resource_group": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"prevent_deletion_if_contains_resources": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  os.Getenv("TF_ACC") == "",
					},
				},
			},
		},

		"recovery_services_vaults": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"recover_soft_deleted_backup_protected_vm": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"managed_disk": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"expand_without_downtime": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
				},
			},
		},

		"storage": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"data_plane_available": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
				},
			},
		},

		"subscription": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"prevent_cancellation_on_destroy": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"postgresql_flexible_server": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"restart_server_on_configuration_value_change": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					},
				},
			},
		},
		"machine_learning": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"purge_soft_deleted_workspace_on_destroy": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"recovery_service": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"vm_backup_stop_protection_and_retain_data_on_destroy": {
						Type:         pluginsdk.TypeBool,
						Optional:     true,
						Default:      false,
						ExactlyOneOf: []string{"features.0.recovery_service.0.vm_backup_stop_protection_and_retain_data_on_destroy", "features.0.recovery_service.0.vm_backup_suspend_protection_and_retain_data_on_destroy"},
					},
					"vm_backup_suspend_protection_and_retain_data_on_destroy": {
						Type:         pluginsdk.TypeBool,
						Optional:     true,
						Default:      false,
						ExactlyOneOf: []string{"features.0.recovery_service.0.vm_backup_stop_protection_and_retain_data_on_destroy", "features.0.recovery_service.0.vm_backup_suspend_protection_and_retain_data_on_destroy"},
					},
					"purge_protected_items_from_vault_on_destroy": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"netapp": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"delete_backups_on_backup_vault_destroy": {
						Description: "When enabled, backups will be deleted when the `azurerm_netapp_backup_vault` resource is destroyed",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
					"prevent_volume_destruction": {
						Description: "When enabled, the volume will not be destroyed, safeguarding from severe data loss",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},
				},
			},
		},

		"databricks_workspace": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"force_delete": {
						Description: "When enabled, the managed resource group that contains the Unity Catalog data will be forcibly deleted when the workspace is destroyed, regardless of contents.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
				},
			},
		},
	}

	if !features.FivePointOh() {
		featuresMap["virtual_machine"].Elem.(*pluginsdk.Resource).Schema["graceful_shutdown"] = &pluginsdk.Schema{
			Type:       pluginsdk.TypeBool,
			Optional:   true,
			Default:    false,
			Deprecated: "'graceful_shutdown' has been deprecated and will be removed from v5.0 of the AzureRM provider.",
		}
	}

	// this is a temporary hack to enable us to gradually add provider blocks to test configurations
	// rather than doing it as a big-bang and breaking all open PR's
	if supportLegacyTestSuite {
		return &pluginsdk.Schema{
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: featuresMap,
			},
		}
	}

	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		MinItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: featuresMap,
		},
	}
}

func expandFeatures(input []interface{}) features.UserFeatures {
	// these are the defaults if omitted from the config
	featuresMap := features.Default()

	if len(input) == 0 || input[0] == nil {
		return featuresMap
	}

	val := input[0].(map[string]interface{})

	if raw, ok := val["api_management"]; ok {
		items := raw.([]interface{})
		if len(items) > 0 && items[0] != nil {
			apimRaw := items[0].(map[string]interface{})
			if v, ok := apimRaw["purge_soft_delete_on_destroy"]; ok {
				featuresMap.ApiManagement.PurgeSoftDeleteOnDestroy = v.(bool)
			}
			if v, ok := apimRaw["recover_soft_deleted"]; ok {
				featuresMap.ApiManagement.RecoverSoftDeleted = v.(bool)
			}
		}
	}

	if raw, ok := val["app_configuration"]; ok {
		items := raw.([]interface{})
		if len(items) > 0 && items[0] != nil {
			appConfRaw := items[0].(map[string]interface{})
			if v, ok := appConfRaw["purge_soft_delete_on_destroy"]; ok {
				featuresMap.AppConfiguration.PurgeSoftDeleteOnDestroy = v.(bool)
			}
			if v, ok := appConfRaw["recover_soft_deleted"]; ok {
				featuresMap.AppConfiguration.RecoverSoftDeleted = v.(bool)
			}
		}
	}

	if raw, ok := val["application_insights"]; ok {
		items := raw.([]interface{})
		if len(items) > 0 && items[0] != nil {
			applicationInsightsRaw := items[0].(map[string]interface{})
			if v, ok := applicationInsightsRaw["disable_generated_rule"]; ok {
				featuresMap.ApplicationInsights.DisableGeneratedRule = v.(bool)
			}
		}
	}

	if raw, ok := val["cognitive_account"]; ok {
		items := raw.([]interface{})
		if len(items) > 0 && items[0] != nil {
			cognitiveRaw := items[0].(map[string]interface{})
			if v, ok := cognitiveRaw["purge_soft_delete_on_destroy"]; ok {
				featuresMap.CognitiveAccount.PurgeSoftDeleteOnDestroy = v.(bool)
			}
		}
	}

	if raw, ok := val["key_vault"]; ok {
		items := raw.([]interface{})
		if len(items) > 0 && items[0] != nil {
			keyVaultRaw := items[0].(map[string]interface{})
			if v, ok := keyVaultRaw["purge_soft_delete_on_destroy"]; ok {
				featuresMap.KeyVault.PurgeSoftDeleteOnDestroy = v.(bool)
			}
			if v, ok := keyVaultRaw["purge_soft_deleted_certificates_on_destroy"]; ok {
				featuresMap.KeyVault.PurgeSoftDeletedCertsOnDestroy = v.(bool)
			}
			if v, ok := keyVaultRaw["purge_soft_deleted_keys_on_destroy"]; ok {
				featuresMap.KeyVault.PurgeSoftDeletedKeysOnDestroy = v.(bool)
			}
			if v, ok := keyVaultRaw["purge_soft_deleted_secrets_on_destroy"]; ok {
				featuresMap.KeyVault.PurgeSoftDeletedSecretsOnDestroy = v.(bool)
			}
			if v, ok := keyVaultRaw["purge_soft_deleted_hardware_security_modules_on_destroy"]; ok {
				featuresMap.KeyVault.PurgeSoftDeletedHSMsOnDestroy = v.(bool)
			}
			if v, ok := keyVaultRaw["purge_soft_deleted_hardware_security_module_keys_on_destroy"]; ok {
				featuresMap.KeyVault.PurgeSoftDeletedHSMKeysOnDestroy = v.(bool)
			}
			if v, ok := keyVaultRaw["recover_soft_deleted_certificates"]; ok {
				featuresMap.KeyVault.RecoverSoftDeletedCerts = v.(bool)
			}
			if v, ok := keyVaultRaw["recover_soft_deleted_key_vaults"]; ok {
				featuresMap.KeyVault.RecoverSoftDeletedKeyVaults = v.(bool)
			}
			if v, ok := keyVaultRaw["recover_soft_deleted_keys"]; ok {
				featuresMap.KeyVault.RecoverSoftDeletedKeys = v.(bool)
			}
			if v, ok := keyVaultRaw["recover_soft_deleted_secrets"]; ok {
				featuresMap.KeyVault.RecoverSoftDeletedSecrets = v.(bool)
			}
			if v, ok := keyVaultRaw["recover_soft_deleted_hardware_security_module_keys"]; ok {
				featuresMap.KeyVault.RecoverSoftDeletedHSMKeys = v.(bool)
			}
		}
	}

	if raw, ok := val["log_analytics_workspace"]; ok {
		items := raw.([]interface{})
		if len(items) > 0 {
			logAnalyticsWorkspaceRaw := items[0].(map[string]interface{})
			if v, ok := logAnalyticsWorkspaceRaw["permanently_delete_on_destroy"]; ok {
				featuresMap.LogAnalyticsWorkspace.PermanentlyDeleteOnDestroy = v.(bool)
			}
		}
	}

	if raw, ok := val["template_deployment"]; ok {
		items := raw.([]interface{})
		if len(items) > 0 {
			templateRaw := items[0].(map[string]interface{})
			if v, ok := templateRaw["delete_nested_items_during_deletion"]; ok {
				featuresMap.TemplateDeployment.DeleteNestedItemsDuringDeletion = v.(bool)
			}
		}
	}

	if raw, ok := val["virtual_machine"]; ok {
		items := raw.([]interface{})
		if len(items) > 0 && items[0] != nil {
			virtualMachinesRaw := items[0].(map[string]interface{})
			if v, ok := virtualMachinesRaw["detach_implicit_data_disk_on_deletion"]; ok {
				featuresMap.VirtualMachine.DetachImplicitDataDiskOnDeletion = v.(bool)
			}
			if v, ok := virtualMachinesRaw["delete_os_disk_on_deletion"]; ok {
				featuresMap.VirtualMachine.DeleteOSDiskOnDeletion = v.(bool)
			}
			if v, ok := virtualMachinesRaw["skip_shutdown_and_force_delete"]; ok {
				featuresMap.VirtualMachine.SkipShutdownAndForceDelete = v.(bool)
			}
		}
	}

	if raw, ok := val["virtual_machine_scale_set"]; ok {
		items := raw.([]interface{})
		if len(items) > 0 {
			scaleSetRaw := items[0].(map[string]interface{})
			if v, ok := scaleSetRaw["reimage_on_manual_upgrade"]; ok {
				featuresMap.VirtualMachineScaleSet.ReimageOnManualUpgrade = v.(bool)
			}
			if v, ok := scaleSetRaw["roll_instances_when_required"]; ok {
				featuresMap.VirtualMachineScaleSet.RollInstancesWhenRequired = v.(bool)
			}
			if v, ok := scaleSetRaw["force_delete"]; ok {
				featuresMap.VirtualMachineScaleSet.ForceDelete = v.(bool)
			}
			if v, ok := scaleSetRaw["scale_to_zero_before_deletion"]; ok {
				featuresMap.VirtualMachineScaleSet.ScaleToZeroOnDelete = v.(bool)
			}
		}
	}

	if raw, ok := val["resource_group"]; ok {
		items := raw.([]interface{})
		if len(items) > 0 {
			resourceGroupRaw := items[0].(map[string]interface{})
			if v, ok := resourceGroupRaw["prevent_deletion_if_contains_resources"]; ok {
				featuresMap.ResourceGroup.PreventDeletionIfContainsResources = v.(bool)
			}
		}
	}

	if raw, ok := val["recovery_services_vaults"]; ok {
		items := raw.([]interface{})
		if len(items) > 0 && items[0] != nil {
			appConfRaw := items[0].(map[string]interface{})
			if v, ok := appConfRaw["recover_soft_deleted_backup_protected_vm"]; ok {
				featuresMap.RecoveryServicesVault.RecoverSoftDeletedBackupProtectedVM = v.(bool)
			}
		}
	}

	if raw, ok := val["managed_disk"]; ok {
		items := raw.([]interface{})
		if len(items) > 0 {
			managedDiskRaw := items[0].(map[string]interface{})
			if v, ok := managedDiskRaw["expand_without_downtime"]; ok {
				featuresMap.ManagedDisk.ExpandWithoutDowntime = v.(bool)
			}
		}
	}
	if raw, ok := val["storage"]; ok {
		items := raw.([]interface{})
		if len(items) > 0 {
			storageRaw := items[0].(map[string]interface{})
			if v, ok := storageRaw["data_plane_available"]; ok {
				featuresMap.Storage.DataPlaneAvailable = v.(bool)
			}
		}
	}

	if raw, ok := val["subscription"]; ok {
		items := raw.([]interface{})
		if len(items) > 0 {
			subscriptionRaw := items[0].(map[string]interface{})
			if v, ok := subscriptionRaw["prevent_cancellation_on_destroy"]; ok {
				featuresMap.Subscription.PreventCancellationOnDestroy = v.(bool)
			}
		}
	}

	if raw, ok := val["postgresql_flexible_server"]; ok {
		items := raw.([]interface{})
		if len(items) > 0 {
			subscriptionRaw := items[0].(map[string]interface{})
			if v, ok := subscriptionRaw["restart_server_on_configuration_value_change"]; ok {
				featuresMap.PostgresqlFlexibleServer.RestartServerOnConfigurationValueChange = v.(bool)
			}
		}
	}

	if raw, ok := val["machine_learning"]; ok {
		items := raw.([]interface{})
		if len(items) > 0 {
			subscriptionRaw := items[0].(map[string]interface{})
			if v, ok := subscriptionRaw["purge_soft_deleted_workspace_on_destroy"]; ok {
				featuresMap.MachineLearning.PurgeSoftDeletedWorkspaceOnDestroy = v.(bool)
			}
		}
	}

	if raw, ok := val["recovery_service"]; ok {
		items := raw.([]interface{})
		if len(items) > 0 {
			recoveryServicesRaw := items[0].(map[string]interface{})
			if v, ok := recoveryServicesRaw["vm_backup_stop_protection_and_retain_data_on_destroy"]; ok {
				featuresMap.RecoveryService.VMBackupStopProtectionAndRetainDataOnDestroy = v.(bool)
			}
			if v, ok := recoveryServicesRaw["vm_backup_suspend_protection_and_retain_data_on_destroy"]; ok {
				featuresMap.RecoveryService.VMBackupSuspendProtectionAndRetainDataOnDestroy = v.(bool)
			}
			if v, ok := recoveryServicesRaw["purge_protected_items_from_vault_on_destroy"]; ok {
				featuresMap.RecoveryService.PurgeProtectedItemsFromVaultOnDestroy = v.(bool)
			}
		}
	}

	if raw, ok := val["netapp"]; ok {
		items := raw.([]interface{})
		if len(items) > 0 {
			netappRaw := items[0].(map[string]interface{})
			if v, ok := netappRaw["delete_backups_on_backup_vault_destroy"]; ok {
				featuresMap.NetApp.DeleteBackupsOnBackupVaultDestroy = v.(bool)
			}
			if v, ok := netappRaw["prevent_volume_destruction"]; ok {
				featuresMap.NetApp.PreventVolumeDestruction = v.(bool)
			}
		}
	}

	if raw, ok := val["databricks_workspace"]; ok {
		items := raw.([]interface{})
		if len(items) > 0 {
			databricksRaw := items[0].(map[string]interface{})
			if v, ok := databricksRaw["force_delete"]; ok {
				featuresMap.DatabricksWorkspace.ForceDelete = v.(bool)
			}
		}
	}

	return featuresMap
}
