// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/streamingendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = StreamingEndpointV0ToV1{}

type StreamingEndpointV0ToV1 struct {
}

func (StreamingEndpointV0ToV1) Schema() map[string]*pluginsdk.Schema {
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
			Optional: true,
			Computed: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"scale_units": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},

		"access_control": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					//lintignore:XS003
					"akamai_signature_header_authentication_key": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"base64_key": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"expiration": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
								"identifier": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
						AtLeastOneOf: []string{"access_control.0.akamai_signature_header_authentication_key", "access_control.0.ip_allow"},
					},
					//lintignore:XS003
					"ip_allow": {
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
						AtLeastOneOf: []string{"access_control.0.akamai_signature_header_authentication_key", "access_control.0.ip_allow"},
					},
				},
			},
		},

		"cdn_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"cdn_profile": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"cdn_provider": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"cross_site_access_policy": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"client_access_policy": {
						Type:         pluginsdk.TypeString,
						Computed:     true,
						Optional:     true,
						AtLeastOneOf: []string{"cross_site_access_policy.0.client_access_policy", "cross_site_access_policy.0.cross_domain_policy"},
					},

					"cross_domain_policy": {
						Type:         pluginsdk.TypeString,
						Computed:     true,
						Optional:     true,
						AtLeastOneOf: []string{"cross_site_access_policy.0.client_access_policy", "cross_site_access_policy.0.cross_domain_policy"},
					},
				},
			},
		},

		"custom_host_names": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"max_cache_age_seconds": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"host_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
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

func (StreamingEndpointV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldIdRaw := rawState["id"].(string)
		oldId, err := streamingendpoints.ParseStreamingEndpointIDInsensitively(oldIdRaw)
		if err != nil {
			return nil, err
		}

		newId := oldId.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q..", oldIdRaw, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
