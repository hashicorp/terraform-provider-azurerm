// Copyright IBM Corp. 2014, 2025
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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FunctionAppFunctionName,
			Description:  "The name of the function.",
		},

		"function_app_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: commonids.ValidateFunctionAppID,
			ForceNew:     true,
			Description:  "The ID of the Function App in which this function should reside.",
		},

		"config_json": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsJSON,
			Description:  "The config for this Function in JSON format.",
		},

		"enabled": {
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Should this function be enabled. Defaults to `true`.",
		},

		"language": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"CSharp",
				"Custom",
				"Java",
				"Javascript",
				"Python",
				"PowerShell",
				"TypeScript",
			}, false), // TODO - find the valida list of strings
			Description: "The language the Function is written in.",
		},

		"file": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			ForceNew: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						Description:  "The filename of the file to be uploaded.",
					},
					"content": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						Description:  "The content of the file.",
					},
				},
			},
		},

		"test_data": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The test data for the function.",
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"config_url": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The URL of the configuration JSON.",
		},

		"url": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The function URL.",
		},

		"invocation_url": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The invocation URL.",
		},

		"script_url": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The script URL.",
		},

		"script_root_path_url": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The Script root path URL.",
		},

		"test_data_url": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The Test data URL.",
		},

		"secrets_file_url": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The URL for the Secrets File.",
		},
	}
}
