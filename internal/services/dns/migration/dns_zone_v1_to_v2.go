// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = DnsZoneV1ToV2{}

type DnsZoneV1ToV2 struct{}

func (DnsZoneV1ToV2) Schema() map[string]*pluginsdk.Schema {
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

		"number_of_record_sets": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"max_number_of_record_sets": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"name_servers": {
			Type:     pluginsdk.TypeSet,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			Set: pluginsdk.HashString,
		},

		"soa_record": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"email": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"host_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"expire_time": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  2419200,
					},

					"minimum_ttl": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  300,
					},

					"refresh_time": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  3600,
					},

					"retry_time": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  300,
					},

					"serial_number": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  1,
					},

					"ttl": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  3600,
					},

					"tags": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},

					"fqdn": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
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

func (DnsZoneV1ToV2) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		id, err := zones.ParseDnsZoneIDInsensitively(oldId)
		if err != nil {
			return rawState, err
		}
		newId := id.ID()
		log.Printf("Updating `id` from %q to %q", oldId, newId)
		rawState["id"] = newId
		return rawState, nil
	}
}
