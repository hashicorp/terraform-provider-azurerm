package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ManagementGroupRemediationV0ToV1{}

type ManagementGroupRemediationV0ToV1 struct{}

func (ManagementGroupRemediationV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"management_group_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"policy_assignment_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"failure_percentage": {
			Type:     pluginsdk.TypeFloat,
			Optional: true,
		},

		"parallel_deployments": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"resource_count": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"location_filters": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"policy_definition_reference_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"policy_definition_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"resource_discovery_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}

func (ManagementGroupRemediationV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldIdRaw := rawState["id"].(string)

		// historically this allowed:
		// regex := regexp.MustCompile(`/providers/[Mm]icrosoft\.[Pp]olicy[Ii]nsights/remediations/`)
		// as such must be corrected
		parsed, err := parse.ParseManagementGroupRemediationIDInsensitively(oldIdRaw)
		if err != nil {
			return rawState, err
		}
		newId := parsed.ToRemediationID().ID()
		log.Printf("[DEBUG] Updating ID from %q to %q", oldIdRaw, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
