// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_trigger_schedule

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/jackofallops/kermit/sdk/datafactory/2018-06-01/datafactory"
)

func dataFactoryTriggerScheduleSchema() map[string]*pluginsdk.Schema {
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

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"schedule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"days_of_month": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeInt,
							ValidateFunc: validation.Any(
								validation.IntBetween(1, 31),
								validation.IntBetween(-31, -1),
							),
						},
					},

					"days_of_week": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 7,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsDayOfTheWeek(false),
						},
					},

					"hours": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeInt,
							ValidateFunc: validation.IntBetween(0, 24),
						},
					},

					"minutes": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeInt,
							ValidateFunc: validation.IntBetween(0, 60),
						},
					},

					"monthly": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MinItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"weekday": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.IsDayOfTheWeek(false),
								},

								"week": {
									Type:     pluginsdk.TypeInt,
									Optional: true,
									ValidateFunc: validation.Any(
										validation.IntBetween(1, 5),
										validation.IntBetween(-5, -1),
									),
								},
							},
						},
					},
				},
			},
		},

		// This time can only be  represented in UTC.
		// An issue has been filed in the SDK for the timezone attribute that doesn't seem to work
		// https://github.com/Azure/azure-sdk-for-go/issues/6244
		"start_time": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Computed:         true, // azignore:AZS007 - pre-existing violation
			DiffSuppressFunc: suppress.RFC3339Time,
			ValidateFunc:     validation.IsRFC3339Time, // times in the past just start immediately
		},

		// This time can only be  represented in UTC.
		// An issue has been filed in the SDK for the timezone attribute that doesn't seem to work
		// https://github.com/Azure/azure-sdk-for-go/issues/6244
		"end_time": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			DiffSuppressFunc: suppress.RFC3339Time,
			ValidateFunc:     validation.IsRFC3339Time, // times in the past just start immediately
		},

		"time_zone": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"frequency": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(datafactory.RecurrenceFrequencyMinute),
			ValidateFunc: validation.StringInSlice([]string{
				string(datafactory.RecurrenceFrequencyMinute),
				string(datafactory.RecurrenceFrequencyHour),
				string(datafactory.RecurrenceFrequencyDay),
				string(datafactory.RecurrenceFrequencyWeek),
				string(datafactory.RecurrenceFrequencyMonth),
			}, false),
		},

		"interval": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"activated": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"pipeline": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			Computed:      true, // azignore:AZS007 - pre-existing violation
			ConflictsWith: []string{"pipeline_parameters"},
			ExactlyOneOf:  []string{"pipeline", "pipeline_name"},
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

		"pipeline_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ExactlyOneOf: []string{"pipeline", "pipeline_name"},
			ValidateFunc: validate.DataFactoryPipelineAndTriggerName(),
		},

		"pipeline_parameters": {
			Type:          pluginsdk.TypeMap,
			Optional:      true,
			Computed:      true, // azignore:AZS007 - pre-existing violation
			ConflictsWith: []string{"pipeline"},
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
	}
}
