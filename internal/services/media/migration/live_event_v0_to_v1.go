// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/liveevents"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = LiveEventV0ToV1{}

type LiveEventV0ToV1 struct {
}

func (LiveEventV0ToV1) Schema() map[string]*pluginsdk.Schema {
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

		"auto_start_enabled": {
			Type:     pluginsdk.TypeBool,
			ForceNew: true,
			Optional: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"input": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					//lintignore:XS003
					"ip_access_control_allow": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"address": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"subnet_prefix_length": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},
							},
						},
						AtLeastOneOf: []string{
							"input.0.ip_access_control_allow", "input.0.access_token",
							"input.0.key_frame_interval_duration", "input.0.streaming_protocol",
						},
					},

					"access_token": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
						ForceNew: true,
						AtLeastOneOf: []string{
							"input.0.ip_access_control_allow", "input.0.access_token",
							"input.0.key_frame_interval_duration", "input.0.streaming_protocol",
						},
					},

					"endpoint": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"protocol": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"url": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},

					"key_frame_interval_duration": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						AtLeastOneOf: []string{
							"input.0.ip_access_control_allow", "input.0.access_token",
							"input.0.key_frame_interval_duration", "input.0.streaming_protocol",
						},
					},

					"streaming_protocol": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						AtLeastOneOf: []string{
							"input.0.ip_access_control_allow", "input.0.access_token",
							"input.0.key_frame_interval_duration", "input.0.streaming_protocol",
						},
					},
				},
			},
		},

		"cross_site_access_policy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_access_policy": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						AtLeastOneOf: []string{"cross_site_access_policy.0.client_access_policy", "cross_site_access_policy.0.cross_domain_policy"},
					},

					"cross_domain_policy": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						AtLeastOneOf: []string{"cross_site_access_policy.0.client_access_policy", "cross_site_access_policy.0.cross_domain_policy"},
					},
				},
			},
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"encoding": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						Default:  "None",
					},

					"key_frame_interval": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "PT2S",
					},

					"preset_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
					},

					"stretch_mode": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "None",
					},
				},
			},
		},

		"hostname_prefix": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"preview": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					//lintignore:XS003
					"ip_access_control_allow": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"address": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"subnet_prefix_length": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
								},
							},
						},
						AtLeastOneOf: []string{
							"preview.0.ip_access_control_allow", "preview.0.alternative_media_id",
							"preview.0.preview_locator", "preview.0.streaming_policy_name",
						},
					},

					"alternative_media_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						AtLeastOneOf: []string{
							"preview.0.ip_access_control_allow", "preview.0.alternative_media_id",
							"preview.0.preview_locator", "preview.0.streaming_policy_name",
						},
					},

					"endpoint": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"protocol": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
								"url": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},
							},
						},
					},

					"preview_locator": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ForceNew: true,
						Computed: true,
						AtLeastOneOf: []string{
							"preview.0.ip_access_control_allow", "preview.0.alternative_media_id",
							"preview.0.preview_locator", "preview.0.streaming_policy_name",
						},
					},

					"streaming_policy_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
						Optional: true,
						ForceNew: true,
						AtLeastOneOf: []string{
							"preview.0.ip_access_control_allow", "preview.0.alternative_media_id",
							"preview.0.preview_locator", "preview.0.streaming_policy_name",
						},
					},
				},
			},
		},

		"transcription_languages": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"use_static_hostname": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},

		"tags": tags.Schema(),
	}
}

func (LiveEventV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldIdRaw := rawState["id"].(string)
		oldId, err := liveevents.ParseLiveEventIDInsensitively(oldIdRaw)
		if err != nil {
			return nil, err
		}

		newId := oldId.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q..", oldIdRaw, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
