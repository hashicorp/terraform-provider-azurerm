// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SpringCloudAppV0ToV1 struct{}

func (s SpringCloudAppV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"service_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"addon_json": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"identity": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"type": {
						Type:     schema.TypeString,
						Required: true,
					},
					"identity_ids": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"principal_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"tenant_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},

		"custom_persistent_disk": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"storage_name": {
						Type:     schema.TypeString,
						Required: true,
					},

					"mount_path": {
						Type:     schema.TypeString,
						Required: true,
					},

					"share_name": {
						Type:     schema.TypeString,
						Required: true,
					},

					"mount_options": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},

					"read_only_enabled": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"is_public": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"https_only": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"ingress_settings": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"backend_protocol": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"read_timeout_in_seconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  300,
					},

					"send_timeout_in_seconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  60,
					},

					"session_affinity": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"session_cookie_max_age": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
					},
				},
			},
		},

		"persistent_disk": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"size_in_gb": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"mount_path": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  "/persistent",
					},
				},
			},
		},

		"public_endpoint_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"tls_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (s SpringCloudAppV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId, err := parse.SpringCloudAppIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
