// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dns_zone

import (
	"math"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the DNS Zone.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"soa_record": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"email": {
						Description:  "The email contact for the SOA record.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.DnsZoneSOARecordEmail,
					},

					"expire_time": {
						Description:  "The expire time for the SOA record.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      2419200,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"minimum_ttl": {
						Description:  "The minimum Time To Live for the SOA record. By convention, it is used to determine the negative caching duration.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      300,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"refresh_time": {
						Description:  "The refresh time for the SOA record.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      3600,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"retry_time": {
						Description:  "The retry time for the SOA record.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      300,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"serial_number": {
						Description:  "The serial number for the SOA record.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      1,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"ttl": {
						Description:  "The Time To Live of the SOA Record in seconds.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      3600,
						ValidateFunc: validation.IntBetween(0, math.MaxInt32),
					},

					"tags": commonschema.Tags(),

					"fqdn": {
						Description: "The fully qualified domain name.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"host_name": {
						Description: "The domain name of the authoritative name server for the SOA record.",
						Type:        pluginsdk.TypeString,
						Computed:    true, // (@jackofallops) - This should not be set or updatable to meet API design, see https://learn.microsoft.com/azure/dns/dns-zones-records#soa-records
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"number_of_record_sets": {
			Description: "The number of records already in the zone.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"max_number_of_record_sets": {
			Description: "Maximum number of Records in the zone.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"name_servers": {
			Description: "A list of values that make up the NS record for the zone.",
			Type:        pluginsdk.TypeSet,
			Computed:    true,
			Set:         pluginsdk.HashString,
			Elem:        &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},
	}
}
