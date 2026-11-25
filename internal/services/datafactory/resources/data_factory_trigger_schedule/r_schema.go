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
			Description:  "The name of the Data Factory Schedule Trigger. See the [Microsoft documentation](https://learn.microsoft.com/azure/data-factory/naming-rules) for all restrictions.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DataFactoryPipelineAndTriggerName(),
		},

		"data_factory_id": {
			Description:  "The Data Factory ID in which to associate the Linked Service with.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: factories.ValidateFactoryID,
		},

		"description": {
			Description:  "The Schedule Trigger's description.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"schedule": {
			Description: "A schedule is capable of limiting or increasing the number of trigger executions specified by the `frequency` and `interval` properties.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MinItems:    1,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"days_of_month": {
						Description: "Day(s) of the month on which the trigger is scheduled. This value can be specified with a monthly frequency only.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeInt,
							ValidateFunc: validation.Any(
								validation.IntBetween(1, 31),
								validation.IntBetween(-31, -1),
							),
						},
					},

					"days_of_week": {
						Description: "Days of the week on which the trigger is scheduled. This value can be specified only with a weekly frequency.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						MaxItems:    7,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsDayOfTheWeek(false),
						},
					},

					"hours": {
						Description: "Hours of the day on which the trigger is scheduled.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeInt,
							ValidateFunc: validation.IntBetween(0, 24),
						},
					},

					"minutes": {
						Description: "Minutes of the hour on which the trigger is scheduled.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
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
									Description:  "The day of the week on which the trigger runs. For example, a `monthly` property with a `weekday` value of `Sunday` means every Sunday of the month.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.IsDayOfTheWeek(false),
								},

								"week": {
									Description: "The occurrence of the specified day during the month. For example, a `monthly` property with `weekday` and `week` values of `Sunday, -1` means the last Sunday of the month.",
									Type:        pluginsdk.TypeInt,
									Optional:    true,
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
			Description:      "The time the Schedule Trigger will start. The time will be represented in UTC.",
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
			Description:      "The time the Schedule Trigger should end. The time will be represented in UTC.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			DiffSuppressFunc: suppress.RFC3339Time,
			ValidateFunc:     validation.IsRFC3339Time, // times in the past just start immediately
		},

		"time_zone": {
			Description:  "The timezone of the start/end time.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"frequency": {
			Description: "The trigger frequency.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     string(datafactory.RecurrenceFrequencyMinute),
			ValidateFunc: validation.StringInSlice([]string{
				string(datafactory.RecurrenceFrequencyMinute),
				string(datafactory.RecurrenceFrequencyHour),
				string(datafactory.RecurrenceFrequencyDay),
				string(datafactory.RecurrenceFrequencyWeek),
				string(datafactory.RecurrenceFrequencyMonth),
			}, false),
		},

		"interval": {
			Description:  "The interval for how often the trigger occurs.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      1,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"activated": {
			Description: "Whether the Data Factory Schedule Trigger is activated.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
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
						Description:  "Reference pipeline name.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.DataFactoryPipelineAndTriggerName(),
					},

					"parameters": {
						Description: "The pipeline parameters that the trigger will act upon.",
						Type:        pluginsdk.TypeMap,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"pipeline_name": {
			Description:  "The Data Factory Pipeline name that the trigger will act on.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ExactlyOneOf: []string{"pipeline", "pipeline_name"},
			ValidateFunc: validate.DataFactoryPipelineAndTriggerName(),
		},

		"pipeline_parameters": {
			Description:   "The pipeline parameters that the trigger will act upon.",
			Type:          pluginsdk.TypeMap,
			Optional:      true,
			Computed:      true, // azignore:AZS007 - pre-existing violation
			ConflictsWith: []string{"pipeline"},
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"annotations": {
			Description: "List of tags that can be used for describing the Data Factory Schedule Trigger.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}
