// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SpringCloudGatewayRouteConfigV0ToV1 struct{}

func (s SpringCloudGatewayRouteConfigV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"spring_cloud_gateway_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"open_api": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"uri": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},

		"protocol": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"spring_cloud_app_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"route": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"description": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"filters": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"order": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},

					"predicates": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"sso_validation_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"title": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"token_relay": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"uri": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"classification_tags": {
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

func (s SpringCloudGatewayRouteConfigV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId, err := parse.SpringCloudGatewayRouteConfigIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
