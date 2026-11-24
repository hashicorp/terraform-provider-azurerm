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
			Description: "The name of the DNS MX Record.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     "@",
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"zone_name": {
			Description: "The DNS Zone where the resource exists.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"record": {
			Description: "A list of values that make up the MX record.",
			Type:        pluginsdk.TypeSet,
			Required:    true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"preference": {
						Description: "String representing the `preference` value of the MX records. Records with lower preference value take priority.",
						// TODO: this should become an Int
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"exchange": {
						Description: "The mail server responsible for the domain covered by the MX record.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
				},
			},
			Set: resourceDnsMxRecordHash,
		},

		"ttl": {
			Description: "The Time To Live (TTL) of the DNS record in seconds.",
			Type:        pluginsdk.TypeInt,
			Required:    true,
		},

		"fqdn": {
			Description: "The FQDN of the DNS MX Record.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"tags": commonschema.Tags(),
	}
}
