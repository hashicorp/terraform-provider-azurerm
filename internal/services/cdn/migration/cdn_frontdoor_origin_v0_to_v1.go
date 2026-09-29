// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigins"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = CdnFrontDoorOriginV0ToV1{}

type CdnFrontDoorOriginV0ToV1 struct{}

func (CdnFrontDoorOriginV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"cdn_frontdoor_origin_group_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"host_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"certificate_name_check_enabled": {
			Type:     pluginsdk.TypeBool,
			Required: true,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"http_port": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"https_port": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"origin_host_header": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"priority": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"private_link": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"location": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"private_link_target_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"request_message": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"target_type": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},

		"weight": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},
	}
}

func (CdnFrontDoorOriginV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		id, err := afdorigins.ParseOriginGroupOriginIDInsensitively(oldId)
		if err != nil {
			return rawState, err
		}

		newId := id.ID()
		log.Printf("[DEBUG] Updating ID from `%s` to `%s`", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
