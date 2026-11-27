// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dns_caa_record

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dnsCaaRecordSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the DNS CAA Record. If you are creating the record in the apex of the zone use `\"@\"` as the name.",
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
			Description: "A list of values that make up the CAA record.",
			Type:        pluginsdk.TypeSet,
			Required:    true,
			Set:         resourceDnsCaaRecordHash,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"flags": {
						Description: "Extensible CAA flags, currently only 1 is implemented to set the issuer critical flag.",
						Type:        pluginsdk.TypeInt,
						Required:    true,
					},

					"tag": {
						Description: "A property tag, options are `issue`, `issuewild`, `iodef`, and `contactemail`.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							"issue",
							"issuewild",
							"iodef",
							"contactemail",
						}, false),
					},

					"value": {
						Description: "A property value such as a registrar domain.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
				},
			},
		},

		"ttl": {
			Description: "The Time To Live (TTL) of the DNS record in seconds.",
			Type:        pluginsdk.TypeInt,
			Required:    true,
		},

		"fqdn": {
			Description: "The FQDN of the DNS CAA Record.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"tags": commonschema.Tags(),
	}
}
