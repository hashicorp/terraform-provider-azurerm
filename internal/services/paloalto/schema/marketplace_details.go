// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2023-09-01/firewalls"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type MarketplaceDetails struct {
	OfferId     string `tfschema:"offer_id"`
	PublisherId string `tfschema:"publisher_id"`
}

func MarketplaceDetailsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"offer_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					Default:      "pan_swfw_cloud_ngfw",
					ValidateFunc: validation.StringIsNotEmpty,
				},

				"publisher_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ForceNew:     true,
					Default:      "paloaltonetworks",
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func ExpandMarketplaceDetails(input []MarketplaceDetails) firewalls.MarketplaceDetails {
	result := firewalls.MarketplaceDetails{
		OfferId:     "pan_swfw_cloud_ngfw",
		PublisherId: "paloaltonetworks",
	}

	if len(input) == 1 {
		p := input[0]

		if p.OfferId != "" {
			result.OfferId = p.OfferId
		}

		if p.PublisherId != "" {
			result.PublisherId = p.PublisherId
		}
	}

	return result
}

func FlattenMarketplaceDetails(input firewalls.MarketplaceDetails) []MarketplaceDetails {
	result := MarketplaceDetails{}

	result.OfferId = input.OfferId
	result.PublisherId = input.PublisherId

	return []MarketplaceDetails{result}
}
