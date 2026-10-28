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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: savedsearches.ValidateWorkspaceID,
		},

		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"category": {
			Type:         pluginsdk.TypeString,
			ForceNew:     true,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"query": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"function_alias": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"function_parameters": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
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
