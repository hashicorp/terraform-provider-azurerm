// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package function_app_function

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the function.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FunctionAppFunctionName,
		},

		"function_app_id": {
			Description:  "The ID of the Function App in which this function should reside.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateFunctionAppID,
		},

		"config_json": {
			Description:  "The config for this Function in JSON format.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsJSON,
		},

		"enabled": {
			Description: "Whether this function should be enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"language": {
			Description: "The language the Function is written in.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"CSharp",
				"Custom",
				"Java",
				"Javascript",
				"Python",
				"PowerShell",
				"TypeScript",
			}, false), // TODO - find the valida list of strings
		},

		"file": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The filename of the file to be uploaded.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"content": {
						Description:  "The content of the file.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"test_data": {
			Description:  "The test data for the function.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"config_url": {
			Description: "The URL of the configuration JSON.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"url": {
			Description: "The function URL.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"invocation_url": {
			Description: "The invocation URL.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"script_url": {
			Description: "The script URL.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"script_root_path_url": {
			Description: "The Script root path URL.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"test_data_url": {
			Description: "The Test data URL.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"secrets_file_url": {
			Description: "The URL for the Secrets File.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
