// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"

	"github.com/hashicorp/go-azure-sdk/resource-manager/webpubsub/2023-02-01/webpubsub"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = WebPubsubV0ToV1{}

type WebPubsubV0ToV1 struct{}

func (WebPubsubV0ToV1) Schema() map[string]*pluginsdk.Schema {
	s := map[string]*pluginsdk.Schema{
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

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"sku": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"capacity": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"live_trace": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"connectivity_logs_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"messaging_logs_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"http_request_logs_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
				},
			},
		},

		"public_network_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"local_auth_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"aad_auth_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"tls_client_cert_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
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

		"public_port": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"server_port": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"external_ip": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"hostname": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"version": {
			Type:     pluginsdk.TypeString,
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
	return s
}

func (WebPubsubV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// the old segment is `WebPubsub` but should be `webPubsub`
		oldID := rawState["id"].(string)

		newID, err := webpubsub.ParseWebPubSubIDInsensitively(oldID)
		if err != nil {
			return nil, err
		}

		rawState["id"] = newID.ID()

		return rawState, nil
	}
}
