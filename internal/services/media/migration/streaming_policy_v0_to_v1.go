// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/streamingpoliciesandstreaminglocators"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = StreamingPolicyV0ToV1{}

type StreamingPolicyV0ToV1 struct {
}

func (StreamingPolicyV0ToV1) Schema() map[string]*pluginsdk.Schema {
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

		"no_encryption_enabled_protocols": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"dash": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
					},

					"download": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
					},

					"hls": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
					},

					"smooth_streaming": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						ForceNew: true,
					},
				},
			},
		},

		"common_encryption_cenc": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled_protocols": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"dash": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									ForceNew: true,
								},

								"download": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									ForceNew: true,
								},

								"hls": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									ForceNew: true,
								},

								"smooth_streaming": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									ForceNew: true,
								},
							},
						},
					},

					"drm_widevine_custom_license_acquisition_url_template": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},

					"drm_playready": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"custom_license_acquisition_url_template": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ForceNew:     true,
									AtLeastOneOf: []string{"common_encryption_cenc.0.drm_playready.0.custom_license_acquisition_url_template", "common_encryption_cenc.0.drm_playready.0.custom_attributes"},
								},

								"custom_attributes": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ForceNew:     true,
									AtLeastOneOf: []string{"common_encryption_cenc.0.drm_playready.0.custom_license_acquisition_url_template", "common_encryption_cenc.0.drm_playready.0.custom_attributes"},
								},
							},
						},
					},

					"default_content_key": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"label": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ForceNew: true,
								},

								"policy_name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ForceNew: true,
								},
							},
						},
					},
				},
			},
		},

		"common_encryption_cbcs": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled_protocols": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"dash": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									ForceNew: true,
								},

								"download": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									ForceNew: true,
								},

								"hls": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									ForceNew: true,
								},

								"smooth_streaming": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									ForceNew: true,
								},
							},
						},
					},

					"drm_fairplay": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"custom_license_acquisition_url_template": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ForceNew:     true,
									AtLeastOneOf: []string{"common_encryption_cbcs.0.drm_fairplay.0.custom_license_acquisition_url_template", "common_encryption_cbcs.0.drm_fairplay.0.allow_persistent_license"},
								},

								"allow_persistent_license": {
									Type:         pluginsdk.TypeBool,
									Optional:     true,
									ForceNew:     true,
									AtLeastOneOf: []string{"common_encryption_cbcs.0.drm_fairplay.0.custom_license_acquisition_url_template", "common_encryption_cbcs.0.drm_fairplay.0.allow_persistent_license"},
								},
							},
						},
					},

					"default_content_key": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"label": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ForceNew: true,
								},

								"policy_name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									ForceNew: true,
								},
							},
						},
					},
				},
			},
		},

		"default_content_key_policy_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},
	}
}

func (StreamingPolicyV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldIdRaw := rawState["id"].(string)
		oldId, err := streamingpoliciesandstreaminglocators.ParseStreamingPolicyIDInsensitively(oldIdRaw)
		if err != nil {
			return nil, err
		}

		newId := oldId.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q..", oldIdRaw, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
