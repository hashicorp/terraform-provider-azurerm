package migration

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-11-01/dedicatedhosts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2022-07-01-preview/configurationassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = AssignmentDedicatedHostV0ToV1{}

type AssignmentDedicatedHostV0ToV1 struct {
}

func (AssignmentDedicatedHostV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"maintenance_configuration_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"dedicated_host_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
	}
}

func (AssignmentDedicatedHostV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldIdRaw := rawState["id"].(string)
		oldId, err := configurationassignments.ParseScopedConfigurationAssignmentIDInsensitively(oldIdRaw)
		if err != nil {
			return nil, fmt.Errorf("parsing old id %q: %+v", oldIdRaw, err)
		}

		dedicatedHostId, err := dedicatedhosts.ParseHostIDInsensitively(oldId.Scope)
		if err != nil {
			return nil, fmt.Errorf("parsing %q as a dedicated hosts id: %+v", oldId.Scope, err)
		}

		newId := configurationassignments.NewScopedConfigurationAssignmentID(dedicatedHostId.ID(), oldId.ConfigurationAssignmentName)
		newIdRaw := newId.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q", oldIdRaw, newIdRaw)
		rawState["id"] = newIdRaw

		return rawState, nil
	}
}
