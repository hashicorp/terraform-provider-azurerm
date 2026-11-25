// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dns_a_record

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func dnsARecordSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the DNS A Record.",
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

		"records": {
			Description:   "List of IPv4 Addresses.",
			Type:          pluginsdk.TypeSet,
			Optional:      true,
			Elem:          &pluginsdk.Schema{Type: pluginsdk.TypeString},
			Set:           pluginsdk.HashString,
			ConflictsWith: []string{"target_resource_id"},
		},

		"ttl": {
			Description: "The Time To Live (TTL) of the DNS record in seconds.",
			Type:        pluginsdk.TypeInt,
			Required:    true,
		},

		"fqdn": {
			Description: "The FQDN of the DNS A Record.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"target_resource_id": {
			Description:   "The Azure resource ID of the target object.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  azure.ValidateResourceID,
			ConflictsWith: []string{"records"},

			// TODO: switch ConflictsWith for ExactlyOneOf when the Provider SDK's updated
		},

		"tags": commonschema.Tags(),
	}
}
