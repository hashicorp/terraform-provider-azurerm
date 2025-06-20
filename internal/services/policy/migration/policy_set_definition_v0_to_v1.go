package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2025-01-01/policysetdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type PolicySetDefinitionV0ToV1 struct{}

var _ pluginsdk.StateUpgrade = PolicySetDefinitionV0ToV1{}

func (p PolicySetDefinitionV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"policy_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"management_group_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"metadata": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"parameters": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"policy_definition_reference": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"policy_definition_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"parameter_values": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"reference_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},

					"policy_group_names": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"policy_definition_group": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"display_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"category": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"description": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"additional_metadata_resource_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},
	}
}

func (p PolicySetDefinitionV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		parsedId, err := policysetdefinitions.ParseProviderPolicySetDefinitionIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		newId := parsedId.ID()
		log.Printf("[DEBUG] Updating ID from `%s` to `%s`", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
