// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package application_insights_analytics_item

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2015-05-01/analyticsitemsapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func applicationInsightsAnalyticsItemSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Application Insights Analytics Item.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"application_insights_id": {
			Description:  "The ID of the Application Insights component on which the Analytics Item exists.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: componentsapis.ValidateComponentID,
		},

		"version": {
			Description: "A string indicating the version of the query format.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"content": {
			Description: "The content for the Analytics Item, for example the query text if `type` is `query`.",
			Type:        pluginsdk.TypeString,
			Required:    true,
		},

		"scope": {
			Description:  "The scope for the Analytics Item.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(analyticsitemsapis.PossibleValuesForItemScope(), false),
		},

		"type": {
			Description: "The type of Analytics Item to create.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(analyticsitemsapis.ItemTypeQuery),
				string(analyticsitemsapis.ItemTypeFunction),
				string(analyticsitemsapis.ItemTypeParameterFolder),
				string(analyticsitemsapis.ItemTypeRecent),
			}, false),
		},

		"function_alias": {
			Description: "The alias to use for the function.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"time_created": {
			Description: "A string containing the time the Analytics Item was created.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"time_modified": {
			Description: "A string containing the time the Analytics Item was last modified.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
