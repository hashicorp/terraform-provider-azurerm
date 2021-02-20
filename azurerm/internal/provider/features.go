package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func schemaFeatures(supportLegacyTestSuite bool) *schema.Schema {
	// NOTE: if there's only one nested field these want to be Required (since there's no point
	//       specifying the block otherwise) - however for 2+ they should be optional
	features := map[string]*schema.Schema{
		"key_vault": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"recover_soft_deleted_key_vaults": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"purge_soft_delete_on_destroy": {
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},

		"network": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"relaxed_locking": {
						Type:     schema.TypeBool,
						Required: true,
					},
				},
			},
		},

		"template_deployment": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"delete_nested_items_during_deletion": {
						Type:     schema.TypeBool,
						Required: true,
					},
				},
			},
		},

		"virtual_machine": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"delete_os_disk_on_deletion": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"graceful_shutdown": {
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},

		"virtual_machine_scale_set": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"roll_instances_when_required": {
						Type:     schema.TypeBool,
						Required: true,
					},
				},
			},
		},
	}

	// this is a temporary hack to enable us to gradually add provider blocks to test configurations
	// rather than doing it as a big-bang and breaking all open PR's
	if supportLegacyTestSuite {
		return &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: features,
			},
		}
	}

	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		MinItems: 1,
		Elem: &schema.Resource{
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

	if raw, ok := val["key_vault"]; ok {
		items := raw.([]interface{})
		if len(items) > 0 {
			keyVaultRaw := items[0].(map[string]interface{})
			if v, ok := keyVaultRaw["purge_soft_delete_on_destroy"]; ok {
				features.KeyVault.PurgeSoftDeleteOnDestroy = v.(bool)
			}
			if v, ok := keyVaultRaw["recover_soft_deleted_key_vaults"]; ok {
				features.KeyVault.RecoverSoftDeletedKeyVaults = v.(bool)
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
		if len(items) > 0 {
			virtualMachinesRaw := items[0].(map[string]interface{})
			if v, ok := virtualMachinesRaw["delete_os_disk_on_deletion"]; ok {
				features.VirtualMachine.DeleteOSDiskOnDeletion = v.(bool)
			}
			if v, ok := virtualMachinesRaw["graceful_shutdown"]; ok {
				features.VirtualMachine.GracefulShutdown = v.(bool)
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
		}
	}

	return features
}
