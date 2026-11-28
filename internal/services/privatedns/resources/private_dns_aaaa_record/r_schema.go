// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_aaaa_record

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatezones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func privateDnsAaaaRecordSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the DNS A Record.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"private_dns_zone_id": {
			Description:  "The ID of the Private DNS Zone where the resource exists.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: privatezones.ValidatePrivateDnsZoneID,
		},

		"records": {
			Description: "A list of IPv6 Addresses.",
			Type:        pluginsdk.TypeSet,
			Required:    true,
			Set:         pluginsdk.HashString, Elem: &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},

		"ttl": {
			Description: "The Time To Live (TTL) of the DNS record in seconds.",
			Type:        pluginsdk.TypeInt,
			Required:    true,
		},

		"fqdn": {
			Description: "The FQDN of the DNS AAAA Record.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"tags": commonschema.Tags(),
	}
}
