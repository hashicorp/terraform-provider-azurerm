// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SpringCloudGatewayV0ToV1 struct{}

func (s SpringCloudGatewayV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"spring_cloud_service_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"api_metadata": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"description": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"documentation_url": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"server_url": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"title": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"version": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},

		"cors": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"credentials_allowed": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"allowed_headers": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"allowed_methods": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"allowed_origins": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"exposed_headers": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"max_age_seconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
				},
			},
		},

		"https_only": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"instance_count": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Default:  1,
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"quota": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"cpu": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "1",
					},

					"memory": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "2Gi",
					},
				},
			},
		},

		"sso": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"client_secret": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"issuer_uri": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"scope": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (s SpringCloudGatewayV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId, err := parse.SpringCloudGatewayIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
