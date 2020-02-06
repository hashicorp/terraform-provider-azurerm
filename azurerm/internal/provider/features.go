package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func schemaFeatures() *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeList,
		// TODO: make this Required in 2.0
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
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
			},
		},
	}
}

func expandFeatures(input []interface{}) features.UserFeatures {
	// TODO: in 2.0 when Required this can become:
	//val := input[0].(map[string]interface{})

	var val map[string]interface{}
	if len(input) > 0 {
		val = input[0].(map[string]interface{})
	}

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
