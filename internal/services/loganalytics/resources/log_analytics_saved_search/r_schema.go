// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_saved_search

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/savedsearches"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func logAnalyticsSavedSearchSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"log_analytics_workspace_id": {
			Description:  "The ID of the Log Analytics Workspace that the Saved Search will be associated with.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: savedsearches.ValidateWorkspaceID,
		},

		"name": {
			Description: "The name of the Log Analytics Saved Search.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"category": {
			Description:  "The category that the Saved Search will be listed under.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"display_name": {
			Description:  "The name that Saved Search will be displayed as.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"query": {
			Description:  "The query expression for the saved search.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"function_alias": {
			Description:  "The function alias if the query serves as a function.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"function_parameters": {
			Description: "The function parameters if the query serves as a function. For more examples and proper syntax refer to [this document](https://learn.microsoft.com/azure/data-explorer/kusto/query/functions/user-defined-functions).",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			ForceNew:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				// https://learn.microsoft.com/azure/data-explorer/kusto/query/functions/user-defined-functions
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9-_]*:([a-z]+(=[^,\n]+)?|\(\*\)|(\([a-zA-Z_][a-zA-Z0-9-_]*:[a-z]+(,[a-zA-Z_][a-zA-Z0-9-_]*:([a-z]+))*\)))(,\s*[a-zA-Z_][a-zA-Z0-9-_]*:([a-z]+(=[^,\n]+)?|\(\*\)|(\([a-zA-Z_][a-zA-Z0-9-_]*:[a-z]+(,\s*[a-zA-Z_][a-zA-Z0-9-_]*:([a-z]+))*\))))*$`),
					"Log Analytics Saved Search Function Parameters must be in the following format: param-name1:type1=default_value1 OR param-name1:type1 OR param-name1:string='string goes here'",
				),
			},
		},

		"tags": commonschema.TagsForceNew(),
	}
}
