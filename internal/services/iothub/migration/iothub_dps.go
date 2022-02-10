package migration

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/provisioningservices/mgmt/2021-10-15/iothub"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = IotHubDPSResourceV0ToV1{}

type IotHubDPSResourceV0ToV1 struct{}

func (IotHubDPSResourceV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return iotHubDPSSchemaForV0AndV1()
}

func (IotHubDPSResourceV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {

		log.Printf("[DEBUG] Updating `public_network_access_enabled` to `true`")

		rawState["public_network_access_enabled"] = true

		return rawState, nil
	}
}

func iotHubDPSSchemaForV0AndV1() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": azure.SchemaResourceGroupName(),

		"location": azure.SchemaLocation(),

		"sku": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"capacity": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
				},
			},
		},

		"linked_hub": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"connection_string": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						Sensitive: true,
					},
					"location": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"apply_allocation_policy": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  features.ThreePointOhBeta(),
					},
					"allocation_weight": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default: func() interface{} {
							if features.ThreePointOhBeta() {
								return 1
							}
							return 0
						}(),
					},
					"hostname": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"ip_filter_rule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"ip_mask": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"action": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"target": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},

		"allocation_policy": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(iothub.AllocationPolicyHashed),
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"device_provisioning_host_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"id_scope": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"service_operations_host_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": tags.Schema(),
	}
}
