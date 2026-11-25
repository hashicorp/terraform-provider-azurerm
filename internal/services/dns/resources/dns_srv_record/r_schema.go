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
			Description: "The name of the DNS SRV Record.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"zone_name": {
			Description: "The DNS Zone where the resource exists.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"record": {
			Description: "A list of values that make up the SRV record.",
			Type:        pluginsdk.TypeSet,
			Required:    true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"priority": {
						Description: "Priority of the SRV record.",
						Type:        pluginsdk.TypeInt,
						Required:    true,
					},

					"weight": {
						Description: "Weight of the SRV record.",
						Type:        pluginsdk.TypeInt,
						Required:    true,
					},

					"port": {
						Description: "Port the service is listening on.",
						Type:        pluginsdk.TypeInt,
						Required:    true,
					},

					"target": {
						Description: "FQDN of the service.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
				},
			},
			Set: resourceDnsSrvRecordHash,
		},

		"ttl": {
			Description: "The Time To Live (TTL) of the DNS record in seconds.",
			Type:        pluginsdk.TypeInt,
			Required:    true,
		},

		"fqdn": {
			Description: "The FQDN of the DNS SRV Record.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"tags": commonschema.Tags(),
	}
}
