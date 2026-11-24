// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_txt_record

import (
	"math"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatezones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func privateDnsTxtRecordSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the DNS TXT Record.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},

		"private_dns_zone_id": {
			Description:  "The ID of the Private DNS Zone where the resource exists.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: privatezones.ValidatePrivateDnsZoneID,
		},

		"record": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"value": {
						Description:  "The value of the TXT record.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringLenBetween(1, 1024),
					},
				},
			},
		},

		"ttl": {
			Description:  "The Time To Live (TTL) of the DNS record in seconds.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, math.MaxInt32),
		},

		"fqdn": {
			Description: "The FQDN of the DNS TXT Record.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"tags": commonschema.Tags(),
	}
}
