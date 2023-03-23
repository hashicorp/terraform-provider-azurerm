package migration

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/notificationhubs/2017-04-01/notificationhubs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = NotificationHubResourceV0ToV1{}

type NotificationHubResourceV0ToV1 struct{}

func (NotificationHubResourceV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"namespace_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"apns_credential": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					// NOTE: APNS supports two modes, certificate auth (v1) and token auth (v2)
					// certificate authentication/v1 is marked for deprecation; as such we're not
					// supporting it at this time.
					"application_mode": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"bundle_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"key_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					// Team ID (within Apple & the Portal) == "AppID" (within the API)
					"team_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"token": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						Sensitive: true,
					},
				},
			},
		},

		"gcm_credential": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"api_key": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						Sensitive: true,
					},
				},
			},
		},

		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (NotificationHubResourceV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldIdRaw := rawState["id"].(string)
		oldId, err := notificationhubs.ParseNotificationHubIDInsensitively(oldIdRaw)
		if err != nil {
			return rawState, fmt.Errorf("parsing ID %q to upgrade: %+v", oldIdRaw, err)
		}

		rawState["id"] = oldId.ID()
		return rawState, nil
	}
}
