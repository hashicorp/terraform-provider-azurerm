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
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"soa_record": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"email": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.DnsZoneSOARecordEmail,
					},

					"expire_time": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      2419200,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"minimum_ttl": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      300,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"refresh_time": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      3600,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"retry_time": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      300,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"serial_number": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      1,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"ttl": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      3600,
						ValidateFunc: validation.IntBetween(0, math.MaxInt32),
					},

					"tags": commonschema.Tags(),

					"fqdn": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"host_name": {
						Type:     pluginsdk.TypeString,
						Computed: true, // (@jackofallops) - This should not be set or updatable to meet API design, see https://learn.microsoft.com/azure/dns/dns-zones-records#soa-records
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
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
			Set:      pluginsdk.HashString,
		},
	}
}
