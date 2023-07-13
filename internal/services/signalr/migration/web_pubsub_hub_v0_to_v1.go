// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"

	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2023-02-01/webpubsub"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = WebPubsubHubV0ToV1{}

type WebPubsubHubV0ToV1 struct{}

func (WebPubsubHubV0ToV1) Schema() map[string]*pluginsdk.Schema {
	s := map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"web_pubsub_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"event_handler": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"url_template": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"user_event_pattern": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"system_events": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"auth": {
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						MinItems: 1,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"managed_identity_id": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},
				},
			},
		},

		"anonymous_connections_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
	return s
}

func (WebPubsubHubV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// the old segment is `WebPubsub` but should be `webPubsub`
		oldID := rawState["id"].(string)
		newID, err := webpubsub.ParseHubIDInsensitively(oldID)
		if err != nil {
			return nil, err
		}
		rawState["id"] = newID.ID()

		// parent webPubsub ID should also has the segment `webPubsub`
		oldWebPubSubID := rawState["web_pubsub_id"].(string)
		newWebPubsubID, err := webpubsub.ParseWebPubSubIDInsensitively(oldWebPubSubID)
		if err != nil {
			return nil, err
		}
		rawState["web_pubsub_id"] = newWebPubsubID.ID()

		return rawState, nil
	}
}
