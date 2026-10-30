// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_zone

import (
	"math"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/privatedns/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func privateDnsZoneSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		// TODO: this can become case-sensitive with a state migration
		"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

		"number_of_record_sets": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"max_number_of_record_sets": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"max_number_of_virtual_network_links": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"max_number_of_virtual_network_links_with_registration": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"soa_record": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"email": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.PrivateDnsZoneSOARecordEmail,
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
						Default:      10,
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

					// This field should be able to be updated since DNS Record Sets API allows to update it.
					// So the issue is submitted on https://github.com/Azure/azure-rest-api-specs/issues/11674
					// Once the issue is fixed, the field will be updated to `Required` property.
					"host_name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					// This field should be able to be updated since DNS Record Sets API allows to update it.
					// So the issue is submitted on https://github.com/Azure/azure-rest-api-specs/issues/11674
					// Once the issue is fixed, the field will be updated to `Optional` property with `Default` attribute.
					"serial_number": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}
