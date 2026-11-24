// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package route_filter

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routefilters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func routeFilterSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The Name which should be used for this Route Filter.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"rule": {
			Type:       pluginsdk.TypeList,
			Optional:   true,
			Computed:   true, // azignore:AZS007 - pre-existing violation
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			MaxItems:   1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the Route Filter rule.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"access": {
						Description: "The access type of the rule.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							string(routefilters.AccessAllow),
						}, false),
					},

					"rule_type": {
						Description: "The rule type of the rule.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							"Community",
						}, false),
					},

					"communities": {
						Description: "The collection for BGP community values to filter on. e.g. `['12076:5010','12076:5020']`.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						MinItems:    1,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}
