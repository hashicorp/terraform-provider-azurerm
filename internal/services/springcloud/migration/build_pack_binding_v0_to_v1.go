// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = BuildPackBindingV0ToV1{}

type BuildPackBindingV0ToV1 struct {
}

func (BuildPackBindingV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"spring_cloud_builder_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"binding_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"launch": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"properties": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"secrets": {
						Type:     pluginsdk.TypeMap,
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

func (BuildPackBindingV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldIdRaw := rawState["id"].(string)
		oldId, err := parse.SpringCloudBuildPackBindingIDInsensitively(oldIdRaw)
		if err != nil {
			return nil, err
		}

		newId := oldId.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q..", oldIdRaw, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
