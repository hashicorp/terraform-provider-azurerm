// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dns_srv_record

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func dnsSrvRecordSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"zone_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"record": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"priority": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"weight": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"port": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"target": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
			Set: resourceDnsSrvRecordHash,
		},

		"ttl": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},

		"fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.Tags(),
	}
}
