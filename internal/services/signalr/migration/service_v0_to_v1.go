// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2023-02-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ServiceV0ToV1{}

type ServiceV0ToV1 struct{}

func (s ServiceV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

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
						Required: true,
					},
				},
			},
		},

		"features": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"flag": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"value": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"upstream_endpoint": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"category_pattern": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"event_pattern": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"hub_pattern": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"url_template": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"cors": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"allowed_origins": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"hostname": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"ip_address": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"public_port": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"server_port": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"primary_access_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"primary_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_access_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
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

func (s ServiceV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		log.Println("[DEBUG] Migrating SignalR from v0 to v1 format")

		// the old segment is `SignalR` but should be `signalR`
		oldId := rawState["id"].(string)
		parsed, err := signalr.ParseSignalRIDInsensitively(oldId)
		if err != nil {
			return rawState, fmt.Errorf("parsing Old Resource ID %q: %+v", oldId, err)
		}

		newId := parsed.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
