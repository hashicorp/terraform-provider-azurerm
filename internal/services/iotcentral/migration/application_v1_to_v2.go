// Copyright (c) HashiCorp, Inc.
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
			Computed: true,
			Optional: true,
			Type:     pluginsdk.TypeString,
		},
		"identity": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"principal_id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"tenant_id": {
					Computed: true,
					Type:     pluginsdk.TypeString,
				},
				"type": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"location": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"public_network_access_enabled": {
			Optional: true,
			Type:     pluginsdk.TypeBool,
		},
		"resource_group_name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"sku": {
			Optional: true,
			Type:     pluginsdk.TypeString,
		},
		"sub_domain": {
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"tags": {
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			Optional: true,
			Type:     pluginsdk.TypeMap,
		},
		"template": {
			Computed: true,
			ForceNew: true,
			Optional: true,
			Type:     pluginsdk.TypeString,
		},
	}
}

func (ApplicationV1ToV2) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
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
