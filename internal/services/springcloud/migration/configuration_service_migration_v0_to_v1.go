// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SpringCloudConfigurationServiceV0ToV1 struct{}

func (s SpringCloudConfigurationServiceV0ToV1) Schema() map[string]*pluginsdk.Schema {
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

		"repository": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"label": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"patterns": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"uri": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"host_key": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"host_key_algorithm": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"password": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"private_key": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"search_paths": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"strict_host_key_checking": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"username": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},
	}
}

func (s SpringCloudConfigurationServiceV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId, err := parse.SpringCloudConfigurationServiceIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
