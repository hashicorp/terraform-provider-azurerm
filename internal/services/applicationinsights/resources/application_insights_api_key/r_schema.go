// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights_api_key

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func applicationInsightsApiKeySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.NoZeroValues,
		},

		"application_insights_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: componentsapis.ValidateComponentID,
		},

		"read_permissions": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			ForceNew: true,
			Set:      pluginsdk.HashString,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"agentconfig", "aggregate", "api", "draft", "extendqueries", "search"}, false),
			},
		},

		"write_permissions": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			ForceNew: true,
			Set:      pluginsdk.HashString,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"annotations"}, false),
			},
		},

		"api_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},
	}
}
