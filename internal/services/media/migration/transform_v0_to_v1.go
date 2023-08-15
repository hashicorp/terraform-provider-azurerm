// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2021-11-01/encodings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = TransformV0ToV1{}

type TransformV0ToV1 struct{}

func (TransformV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"media_services_account_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		//lintignore:XS003
		"output": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"on_error_action": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
					//lintignore:XS003
					"builtin_preset": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"preset_name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},
					//lintignore:XS003
					"audio_analyzer_preset": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"audio_language": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"audio_analysis_mode": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},
					//lintignore:XS003
					"video_analyzer_preset": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"audio_language": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"audio_analysis_mode": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"insights_type": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},
					//lintignore:XS003
					"face_detector_preset": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"analysis_resolution": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},
					"relative_priority": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},
	}
}

func (TransformV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldIdRaw := rawState["id"].(string)
		oldId, err := encodings.ParseTransformIDInsensitively(oldIdRaw)
		if err != nil {
			return nil, err
		}

		newId := oldId.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q..", oldIdRaw, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
