// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic_app_integration_account_batch_configuration

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountbatchconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func logicAppIntegrationAccountBatchConfigurationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Logic App Integration Account Batch Configuration. Only Alphanumeric characters allowed.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IntegrationAccountBatchConfigurationName(),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"integration_account_name": {
			Description:  "The name of the Logic App Integration Account.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IntegrationAccountName(),
		},

		"batch_group_name": {
			Description:  "The batch group name of the Logic App Integration Batch Configuration.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IntegrationAccountBatchConfigurationBatchGroupName(),
		},

		"release_criteria": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"batch_size": {
						Description:  "The batch size in bytes for the Logic App Integration Batch Configuration.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(1, 83886080),
						AtLeastOneOf: []string{"release_criteria.0.batch_size", "release_criteria.0.message_count", "release_criteria.0.recurrence"},
					},

					"message_count": {
						Description:  "The message count for the Logic App Integration Batch Configuration.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(1, 8000),
						AtLeastOneOf: []string{"release_criteria.0.batch_size", "release_criteria.0.message_count", "release_criteria.0.recurrence"},
					},

					"recurrence": {
						Type:         pluginsdk.TypeList,
						Optional:     true,
						MaxItems:     1,
						AtLeastOneOf: []string{"release_criteria.0.batch_size", "release_criteria.0.message_count", "release_criteria.0.recurrence"},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"frequency": {
									Description:  "The frequency of the schedule.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(integrationaccountbatchconfigurations.PossibleValuesForRecurrenceFrequency(), false),
								},

								"interval": {
									Description:  "The number of `frequency`s between runs.",
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 100),
								},

								"end_time": {
									Description:      "The end time of the schedule, formatted as an RFC3339 string.",
									Type:             pluginsdk.TypeString,
									Optional:         true,
									ValidateFunc:     validation.IsRFC3339Time,
									DiffSuppressFunc: suppress.RFC3339Time,
								},

								"schedule": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"hours": {
												Description: "A list containing a single item, which specifies the Hour interval at which this recurrence should be triggered.",
												Type:        pluginsdk.TypeSet,
												Optional:    true,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeInt,
													ValidateFunc: validation.IntBetween(0, 23),
												},
											},

											"minutes": {
												Description: "A list containing a single item which specifies the Minute interval at which this recurrence should be triggered.",
												Type:        pluginsdk.TypeSet,
												Optional:    true,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeInt,
													ValidateFunc: validation.IntBetween(0, 59),
												},
											},

											"month_days": {
												Description:   "A list of days of the month that the job should execute on.",
												Type:          pluginsdk.TypeSet,
												Optional:      true,
												ConflictsWith: []string{"release_criteria.0.recurrence.0.schedule.0.week_days"},
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeInt,
													ValidateFunc: validation.All(
														validation.IntBetween(-31, 31),
														validation.IntNotInSlice([]int{0}),
													),
												},
											},

											"monthly": {
												Type:          pluginsdk.TypeSet,
												Optional:      true,
												ConflictsWith: []string{"release_criteria.0.recurrence.0.schedule.0.week_days"},
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"weekday": {
															Description:  "The day of the occurrence.",
															Type:         pluginsdk.TypeString,
															Required:     true,
															ValidateFunc: validation.StringInSlice(integrationaccountbatchconfigurations.PossibleValuesForDayOfWeek(), false),
														},

														"week": {
															Description: "The occurrence of the week within the month.",
															Type:        pluginsdk.TypeInt,
															Required:    true,
															ValidateFunc: validation.All(
																validation.IntBetween(-5, 5),
																validation.IntNotInSlice([]int{0}),
															),
														},
													},
												},
											},

											"week_days": {
												Description:   "A list of days of the week that the job should execute on.",
												Type:          pluginsdk.TypeSet,
												Optional:      true,
												ConflictsWith: []string{"release_criteria.0.recurrence.0.schedule.0.month_days", "release_criteria.0.recurrence.0.schedule.0.monthly"},
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeString,
													ValidateFunc: validation.StringInSlice(integrationaccountbatchconfigurations.PossibleValuesForDaysOfWeek(), false),
												},
											},
										},
									},
								},

								"start_time": {
									Description:      "The start time of the schedule, formatted as an RFC3339 string.",
									Type:             pluginsdk.TypeString,
									Optional:         true,
									ValidateFunc:     validation.IsRFC3339Time,
									DiffSuppressFunc: suppress.RFC3339Time,
								},

								"time_zone": {
									Description:  "The timezone of the start/end time.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validate.BatchConfigurationRecurrenceTimeZone(),
								},
							},
						},
					},
				},
			},
		},

		"metadata": {
			Description: "A JSON mapping of any Metadata for this Logic App Integration Account Batch Configuration.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}
