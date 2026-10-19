// Copyright IBM Corp. 2014, 2025
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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"application_insights_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: componentsapis.ValidateComponentID,
		},

		"version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"content": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"scope": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(analyticsitemsapis.PossibleValuesForItemScope(), false),
		},

		"type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(analyticsitemsapis.ItemTypeQuery),
				string(analyticsitemsapis.ItemTypeFunction),
				string(analyticsitemsapis.ItemTypeParameterFolder),
				string(analyticsitemsapis.ItemTypeRecent),
			}, false),
		},

		"function_alias": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"time_created": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_modified": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
