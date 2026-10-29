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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IntegrationAccountBatchConfigurationName(),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"integration_account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IntegrationAccountName(),
		},

		"batch_group_name": {
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
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(1, 83886080),
						AtLeastOneOf: []string{"release_criteria.0.batch_size", "release_criteria.0.message_count", "release_criteria.0.recurrence"},
					},

					"message_count": {
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
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(integrationaccountbatchconfigurations.PossibleValuesForRecurrenceFrequency(), false),
								},

								"interval": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 100),
								},

								"end_time": {
									Type:             pluginsdk.TypeString,
									Optional:         true,
									DiffSuppressFunc: suppress.RFC3339Time,
									ValidateFunc:     validation.IsRFC3339Time,
								},

								"schedule": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"hours": {
												Type:     pluginsdk.TypeSet,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeInt,
													ValidateFunc: validation.IntBetween(0, 23),
												},
											},

											"minutes": {
												Type:     pluginsdk.TypeSet,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeInt,
													ValidateFunc: validation.IntBetween(0, 59),
												},
											},

											"month_days": {
												Type:     pluginsdk.TypeSet,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type: pluginsdk.TypeInt,
													ValidateFunc: validation.All(
														validation.IntBetween(-31, 31),
														validation.IntNotInSlice([]int{0}),
													),
												},
												ConflictsWith: []string{"release_criteria.0.recurrence.0.schedule.0.week_days"},
											},

											"monthly": {
												Type:     pluginsdk.TypeSet,
												Optional: true,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"weekday": {
															Type:         pluginsdk.TypeString,
															Required:     true,
															ValidateFunc: validation.StringInSlice(integrationaccountbatchconfigurations.PossibleValuesForDayOfWeek(), false),
														},

														"week": {
															Type:     pluginsdk.TypeInt,
															Required: true,
															ValidateFunc: validation.All(
																validation.IntBetween(-5, 5),
																validation.IntNotInSlice([]int{0}),
															),
														},
													},
												},
												ConflictsWith: []string{"release_criteria.0.recurrence.0.schedule.0.week_days"},
											},

											"week_days": {
												Type:     pluginsdk.TypeSet,
												Optional: true,
												Elem: &pluginsdk.Schema{
													Type:         pluginsdk.TypeString,
													ValidateFunc: validation.StringInSlice(integrationaccountbatchconfigurations.PossibleValuesForDaysOfWeek(), false),
												},
												ConflictsWith: []string{"release_criteria.0.recurrence.0.schedule.0.month_days", "release_criteria.0.recurrence.0.schedule.0.monthly"},
											},
										},
									},
								},

								"start_time": {
									Type:             pluginsdk.TypeString,
									Optional:         true,
									DiffSuppressFunc: suppress.RFC3339Time,
									ValidateFunc:     validation.IsRFC3339Time,
								},

								"time_zone": {
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
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}
