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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DataFactoryPipelineAndTriggerName(),
		},

		"data_factory_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: factories.ValidateFactoryID,
		},

		"frequency": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInEnumSlice(datafactory.PossibleTumblingWindowFrequencyValues(), false),
		},

		"interval": {
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.DataFactoryPipelineAndTriggerName(),
					},

					"parameters": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"start_time": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			DiffSuppressFunc: suppress.RFC3339Time,
			ValidateFunc:     validation.IsRFC3339Time,
		},

		"activated": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"additional_properties": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"annotations": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"delay": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.TriggerTimespan,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"end_time": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			DiffSuppressFunc: suppress.RFC3339Time,
			ValidateFunc:     validation.IsRFC3339Time,
		},

		"max_concurrency": {
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
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(1),
					},

					"interval": {
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
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.TriggerTimespan,
					},

					"size": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.TriggerTimespan,
					},

					"trigger_name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},
	}
}
