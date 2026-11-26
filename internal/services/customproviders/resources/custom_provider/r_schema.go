// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package custom_provider

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/customproviders/2018-09-01-preview/customresourceprovider"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/customproviders/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func customProviderSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Custom Provider.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CustomProviderName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"resource_type": {
			Type:         pluginsdk.TypeSet,
			Optional:     true,
			AtLeastOneOf: []string{"resource_type", "action"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the route definition.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.NoZeroValues,
					},
					"endpoint": {
						Description:  "The endpoint of the route definition.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsURLWithHTTPS,
					},
					"routing_type": {
						Description:  "The routing type that is supported for the resource request.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(customresourceprovider.ResourceTypeRoutingProxy),
						ValidateFunc: validation.StringInSlice(customresourceprovider.PossibleValuesForResourceTypeRouting(), false),
					},
				},
			},
		},

		"action": {
			Type:         pluginsdk.TypeSet,
			Optional:     true,
			AtLeastOneOf: []string{"resource_type", "action"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the action.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.NoZeroValues,
					},
					"endpoint": {
						Description:  "The endpoint of the action.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsURLWithHTTPS,
					},
				},
			},
		},

		"validation": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"specification": {
						Description:  "The endpoint where the validation specification is located.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.IsURLWithHTTPS,
					},
				},
			},
		},

		"tags": commonschema.TagsForceNew(),
	}
}
