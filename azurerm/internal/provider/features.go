package provider

import (
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func schemaFeatures() *schema.Schema {
	features := map[string]*schema.Schema{
		"virtual_machine": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"delete_os_disk_on_deletion": {
						Type:     schema.TypeBool,
						Required: true,
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

	runningAcceptanceTests := os.Getenv("TF_ACC") != ""
	if runningAcceptanceTests {
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
	features := features.UserFeatures{
		// NOTE: ensure all nested objects are fully populated
		VirtualMachine: features.VirtualMachineFeatures{
			DeleteOSDiskOnDeletion: true,
		},
		VirtualMachineScaleSet: features.VirtualMachineScaleSetFeatures{
			RollInstancesWhenRequired: true,
		},
	}

	if len(input) == 0 || input[0] == nil {
		return features
	}

	val := input[0].(map[string]interface{})

	if raw, ok := val["virtual_machine"]; ok {
		items := raw.([]interface{})
		if len(items) > 0 {
			virtualMachinesRaw := items[0].(map[string]interface{})
			if v, ok := virtualMachinesRaw["delete_os_disk_on_deletion"]; ok {
				features.VirtualMachine.DeleteOSDiskOnDeletion = v.(bool)
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
