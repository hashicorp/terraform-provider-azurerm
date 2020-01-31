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

	features := features.UserFeatures{
		// NOTE: ensure all nested objects are fully populated
		VirtualMachine: features.VirtualMachineFeatures{
			DeleteOSDiskOnDeletion: true,
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

	return features
}
