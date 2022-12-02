package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type KustoAttachedClusterV0ToV1 struct{}

func (s KustoAttachedClusterV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"sku": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"capacity": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Computed: true,
					},
				},
			},
		},

		"allowed_fqdns": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"allowed_ip_ranges": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"trusted_external_tenants": {
			Type:       pluginsdk.TypeList,
			Optional:   true,
			Computed:   true,
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"optimized_auto_scale": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"minimum_instances": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
					"maximum_instances": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
				},
			},
		},

		"virtual_network_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"subnet_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"engine_public_ip_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"data_management_public_ip_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"language_extensions": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"engine": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"uri": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"data_ingestion_uri": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"public_ip_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"outbound_network_access_restricted": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"double_encryption_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},

		"auto_stop_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"disk_encryption_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"streaming_ingestion_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"purge_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"zones": commonschema.ZonesMultipleOptionalForceNew(),

		"tags": tags.Schema(),
	}
}

func (s KustoAttachedClusterV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId, err := parse.ClusterIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
