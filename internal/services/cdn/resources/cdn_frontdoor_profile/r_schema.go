// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_profile

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func cdnFrontdoorProfileSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FrontDoorName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"response_timeout_seconds": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      120,
			ValidateFunc: validation.IntBetween(16, 240),
		},

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(profiles.SkuNamePremiumAzureFrontDoor),
				string(profiles.SkuNameStandardAzureFrontDoor),
			}, false),
		},

		"log_scrubbing_rule": {
			Type:     pluginsdk.TypeSet,
			MaxItems: 3,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"match_variable": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice(
							profiles.PossibleValuesForScrubbingRuleEntryMatchVariable(),
							false,
						),
					},
				},
			},
		},

		"tags": commonschema.Tags(),

		"resource_guid": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
