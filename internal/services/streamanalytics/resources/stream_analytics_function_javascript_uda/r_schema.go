// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_function_javascript_uda

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/streamingjobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func streamAnalyticsFunctionJavascriptUdaSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FunctionName,
		},

		"stream_analytics_job_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: streamingjobs.ValidateStreamingJobID,
		},

		"input": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"any",
							"array",
							"bigint",
							"datetime",
							"float",
							"nvarchar(max)",
							"record",
						}, false),
					},

					"configuration_parameter": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"output": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"any",
							"array",
							"bigint",
							"datetime",
							"float",
							"nvarchar(max)",
							"record",
						}, false),
					},
				},
			},
		},

		"script": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}
