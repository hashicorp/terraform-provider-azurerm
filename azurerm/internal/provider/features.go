package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

func schemaFeatures(supportLegacyTestSuite bool) *pluginsdk.Schema {
	// NOTE: if there's only one nested field these want to be Required (since there's no point
	//       specifying the block otherwise) - however for 2+ they should be optional
	features := map[string]*pluginsdk.Schema{
		// lintignore:XS003
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
					"recover_soft_deleted_key_vaults": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"purge_soft_delete_on_destroy": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
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
						Required: true,
					},
				},
			},
		},

		"network": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"relaxed_locking": {
						Type:     pluginsdk.TypeBool,
						Required: true,
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
					"delete_os_disk_on_deletion": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"graceful_shutdown": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
					"skip_shutdown_and_force_delete": {
						Type:     schema.TypeBool,
						Optional: true,
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
					},
					"roll_instances_when_required": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
				},
			},
		},
	}

	// this is a temporary hack to enable us to gradually add provider blocks to test configurations
	// rather than doing it as a big-bang and breaking all open PR's
	if supportLegacyTestSuite {
		return &pluginsdk.Schema{
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: features,
			},
		}
	}

	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		MinItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: features,
		},
	}
}

func expandFeatures(input []interface{}) features.UserFeatures {
	// these are the defaults if omitted from the config
	features := features.Default()

	if len(input) == 0 || input[0] == nil {
		return features
	}

	val := input[0].(map[string]interface{})

	if raw, ok := val["cognitive_account"]; ok {
		items := raw.([]interface{})
		if len(items) > 0 && items[0] != nil {
			cognitiveRaw := items[0].(map[string]interface{})
			if v, ok := cognitiveRaw["purge_soft_delete_on_destroy"]; ok {
				features.CognitiveAccount.PurgeSoftDeleteOnDestroy = v.(bool)
			}
		}
	}

	if raw, ok := val["key_vault"]; ok {
		items := raw.([]interface{})
		if len(items) > 0 && items[0] != nil {
			keyVaultRaw := items[0].(map[string]interface{})
			if v, ok := keyVaultRaw["purge_soft_delete_on_destroy"]; ok {
				features.KeyVault.PurgeSoftDeleteOnDestroy = v.(bool)
			}
			if v, ok := keyVaultRaw["recover_soft_deleted_key_vaults"]; ok {
				features.KeyVault.RecoverSoftDeletedKeyVaults = v.(bool)
			}
		}
	}

	if raw, ok := val["log_analytics_workspace"]; ok {
		items := raw.([]interface{})
		if len(items) > 0 {
			logAnalyticsWorkspaceRaw := items[0].(map[string]interface{})
			if v, ok := logAnalyticsWorkspaceRaw["permanently_delete_on_destroy"]; ok {
				features.LogAnalyticsWorkspace.PermanentlyDeleteOnDestroy = v.(bool)
			}
		}
	}

	if raw, ok := val["network"]; ok {
		items := raw.([]interface{})
		if len(items) > 0 {
			networkRaw := items[0].(map[string]interface{})
			if v, ok := networkRaw["relaxed_locking"]; ok {
				features.Network.RelaxedLocking = v.(bool)
			}
		}
	}

	if raw, ok := val["template_deployment"]; ok {
		items := raw.([]interface{})
		if len(items) > 0 {
			networkRaw := items[0].(map[string]interface{})
			if v, ok := networkRaw["delete_nested_items_during_deletion"]; ok {
				features.TemplateDeployment.DeleteNestedItemsDuringDeletion = v.(bool)
			}
		}
	}

	if raw, ok := val["virtual_machine"]; ok {
		items := raw.([]interface{})
		if len(items) > 0 && items[0] != nil {
			virtualMachinesRaw := items[0].(map[string]interface{})
			if v, ok := virtualMachinesRaw["delete_os_disk_on_deletion"]; ok {
				features.VirtualMachine.DeleteOSDiskOnDeletion = v.(bool)
			}
			if v, ok := virtualMachinesRaw["graceful_shutdown"]; ok {
				features.VirtualMachine.GracefulShutdown = v.(bool)
			}
			if v, ok := virtualMachinesRaw["skip_shutdown_and_force_delete"]; ok {
				features.VirtualMachine.SkipShutdownAndForceDelete = v.(bool)
			}
		}
	}

	if raw, ok := val["virtual_machine_scale_set"]; ok {
		items := raw.([]interface{})
		if len(items) > 0 {
			scaleSetRaw := items[0].(map[string]interface{})
			if v, ok := scaleSetRaw["roll_instances_when_required"]; ok {
				features.VirtualMachineScaleSet.RollInstancesWhenRequired = v.(bool)
			}
			if v, ok := scaleSetRaw["force_delete"]; ok {
				features.VirtualMachineScaleSet.ForceDelete = v.(bool)
			}
		}
	}

	return features
}
