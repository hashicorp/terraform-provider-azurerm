// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package availability_set

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func availabilitySetSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z0-9]([-._a-zA-Z0-9]{0,78}[a-zA-Z0-9_])?$"),
				"The Availability set name can contain only letters, numbers, periods (.), hyphens (-),and underscores (_), up to 80 characters, and it must begin a letter or number and end with a letter, number or underscore.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"platform_update_domain_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      5,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 20),
		},

		"platform_fault_domain_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      3,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 3),
		},

		"managed": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
			ForceNew: true,
		},

		"proximity_placement_group_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,

			// We have to ignore case due to incorrect capitalisation of resource group name in
			// proximity placement group ID in the response we get from the API request
			//
			// todo can be removed when https://github.com/Azure/azure-sdk-for-go/issues/5699 is fixed
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"tags": commonschema.Tags(),
	}
}
