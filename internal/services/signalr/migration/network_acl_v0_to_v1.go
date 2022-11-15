package migration

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2022-02-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = NetworkAclV0ToV1{}

type NetworkAclV0ToV1 struct{}

func (n NetworkAclV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"signalr_service_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"default_action": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"public_network": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"allowed_request_types": {
						Type:          pluginsdk.TypeSet,
						Optional:      true,
						ConflictsWith: []string{"public_network.0.denied_request_types"},
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"denied_request_types": {
						Type:          pluginsdk.TypeSet,
						Optional:      true,
						ConflictsWith: []string{"public_network.0.allowed_request_types"},
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"private_endpoint": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"allowed_request_types": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"denied_request_types": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},
	}
}

func (n NetworkAclV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		log.Println("[DEBUG] Migrating SignalR Network ACL from v0 to v1 format")

		// the old segment is `SignalR` but should be `signalR`
		oldId := rawState["id"].(string)
		parsed, err := signalr.ParseSignalRIDInsensitively(oldId)
		if err != nil {
			return rawState, fmt.Errorf("parsing Old Resource ID %q: %+v", oldId, err)
		}

		newId := parsed.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
		rawState["id"] = newId

		// SignalR Service ID is the same as ID, but let's call the values out specifically
		oldServiceId := rawState["signalr_service_id"].(string)
		newServiceId := parsed.ID()
		log.Printf("[DEBUG] Updating SignalR Service ID from %q to %q", oldServiceId, newServiceId)
		rawState["signalr_service_id"] = newId

		return rawState, nil
	}
}
