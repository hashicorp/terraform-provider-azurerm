package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type HealthCareIoTConnectorV0ToV1 struct{}

func (s HealthCareIoTConnectorV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"workspace_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": commonschema.Location(),

		"identity": commonschema.SystemAssignedIdentityOptional(),

		"eventhub_namespace_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"eventhub_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"eventhub_consumer_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"device_mapping_json": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"tags": commonschema.Tags(),
	}
}

func (s HealthCareIoTConnectorV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId, err := parse.MedTechServiceIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
