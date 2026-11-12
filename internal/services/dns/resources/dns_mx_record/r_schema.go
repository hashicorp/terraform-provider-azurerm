// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dns_mx_record

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func dnsMxRecordSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  "@",
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
					"preference": {
						// TODO: this should become an Int
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"exchange": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
			Set: resourceDnsMxRecordHash,
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
