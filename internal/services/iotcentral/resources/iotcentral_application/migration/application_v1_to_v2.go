// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/iotcentral/2021-11-01-preview/apps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ApplicationV1ToV2{}

type ApplicationV1ToV2 struct{}

func (ApplicationV1ToV2) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"display_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},
		"identity": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"principal_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"tenant_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"type": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
			}},
		},
		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"sku": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"sub_domain": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},
		"template": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},
	}
}

func (ApplicationV1ToV2) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
		oldId := rawState["id"].(string)
		id, err := apps.ParseIotAppIDInsensitively(oldId)
		if err != nil {
			return rawState, err
		}

		newId := id.ID()
		log.Printf("Updating `id` from %q to %q", oldId, newId)
		rawState["id"] = newId
		return rawState, nil
	}
}
