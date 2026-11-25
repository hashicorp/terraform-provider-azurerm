// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_trigger_tumbling_window

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func dataFactoryTriggerTumblingWindowSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Data Factory Tumbling Window Trigger.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DataFactoryPipelineAndTriggerName(),
		},

		"data_factory_id": {
			Description:  "The ID of Data Factory in which to associate the Trigger with.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: factories.ValidateFactoryID,
		},

		"frequency": {
			Description:  "The frequency of Tumbling Window.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInEnumSlice(datafactory.PossibleTumblingWindowFrequencyValues(), false),
		},

		"interval": {
			Description:  "The interval of Tumbling Window.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"pipeline": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The Data Factory Pipeline name that the trigger will act on.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.DataFactoryPipelineAndTriggerName(),
					},

					"parameters": {
						Description: "The Data Factory Pipeline parameters that the trigger will act on.",
						Type:        pluginsdk.TypeMap,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"start_time": {
			Description:      "The start time of Tumbling Window, formatted as an RFC3339 string.",
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			DiffSuppressFunc: suppress.RFC3339Time,
			ValidateFunc:     validation.IsRFC3339Time,
		},

		"activated": {
			Description: "Whether the Data Factory Tumbling Window Trigger is activated.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"additional_properties": {
			Description: "A map of additional properties to associate with the Data Factory Tumbling Window Trigger.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"annotations": {
			Description: "List of tags that can be used for describing the Data Factory Tumbling Window Trigger.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"delay": {
			Description:  "How long the trigger waits before triggering new run. Formatted as an `D.HH:MM:SS`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.TriggerTimespan,
		},

		"description": {
			Description:  "The description for the Data Factory Tumbling Window Trigger.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"end_time": {
			Description:      "The end time of Tumbling Window, formatted as an RFC3339 string.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			DiffSuppressFunc: suppress.RFC3339Time,
			ValidateFunc:     validation.IsRFC3339Time,
		},

		"max_concurrency": {
			Description:  "The max number for simultaneous trigger run fired by Tumbling Window.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      50,
			ValidateFunc: validation.IntBetween(1, 50),
		},

		"retry": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"count": {
						Description:  "The maximum retry attempts if the pipeline run failed.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(1),
					},

					"interval": {
						Description:  "The Interval in seconds between each retry if the pipeline run failed.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      30,
						ValidateFunc: validation.IntAtLeast(0),
					},
				},
			},
		},

		"trigger_dependency": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"offset": {
						Description:  "The offset of the dependency trigger.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.TriggerTimespan,
					},

					"size": {
						Description:  "The size of the dependency tumbling window.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.TriggerTimespan,
					},

					"trigger_name": {
						Description: "The dependency trigger name. If not specified, it will use self dependency.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},
				},
			},
		},
	}
}
