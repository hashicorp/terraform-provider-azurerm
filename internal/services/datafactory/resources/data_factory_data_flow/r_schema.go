// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_data_flow

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dataFactoryDataFlowSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the Data Factory Data Flow.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"data_factory_id": {
			Description:  "The ID of Data Factory in which to associate the Data Flow with.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: factories.ValidateFactoryID,
		},

		"script": {
			Description:  "The script for the Data Factory Data Flow.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			AtLeastOneOf: []string{"script", "script_lines"},
		},

		"script_lines": {
			Description:  "The script lines for the Data Factory Data Flow.",
			Type:         pluginsdk.TypeList,
			Optional:     true,
			AtLeastOneOf: []string{"script", "script_lines"},
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"source": helpers.SchemaForDataFlowSourceAndSink(),

		"sink": helpers.SchemaForDataFlowSourceAndSink(),

		"transformation": helpers.SchemaForDataFlowSourceTransformation(),

		"annotations": {
			Description: "List of tags that can be used for describing the Data Factory Data Flow.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"description": {
			Description:  "The description for the Data Factory Data Flow.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"folder": {
			Description:  "The folder that this Data Flow is in. If not specified, the Data Flow will appear at the root level.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}
