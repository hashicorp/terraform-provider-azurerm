// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/datafactory/2018-06-01/datafactory" // nolint: staticcheck
)

type TriggerScheduleDataSource struct{}

type TriggerScheduleDataSourceModel struct {
	Name          string            `tfschema:"name"`
	DataFactoryID string            `tfschema:"data_factory_id"`
	Description   string            `tfschema:"description"`
	Schedule      []TriggerSchedule `tfschema:"schedule"`
	StartTime     string            `tfschema:"start_time"`
	EndTime       string            `tfschema:"end_time"`
	TimeZone      string            `tfschema:"time_zone"`
	Frequency     string            `tfschema:"frequency"`
	Interval      int64             `tfschema:"interval"`
	Activated     bool              `tfschema:"activated"`
	PipelineName  string            `tfschema:"pipeline_name"`
	Annotations   []string          `tfschema:"annotations"`
}

type TriggerSchedule struct {
	DaysOfMonth []int64                          `tfschema:"days_of_month"`
	DaysOfWeek  []string                         `tfschema:"days_of_week"`
	Hours       []int64                          `tfschema:"hours"`
	Minutes     []int64                          `tfschema:"minutes"`
	Monthly     []TriggerScheduleScheduleMonthly `tfschema:"monthly"`
}

type TriggerScheduleScheduleMonthly struct {
	Weekday string `tfschema:"weekday"`
	Week    int64  `tfschema:"week"`
}

var _ sdk.DataSource = TriggerScheduleDataSource{}

func (d TriggerScheduleDataSource) ModelObject() interface{} {
	return &TriggerScheduleDataSourceModel{}
}

func (d TriggerScheduleDataSource) ResourceType() string {
	return "azurerm_data_factory_trigger_schedule"
}

func (d TriggerScheduleDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.DataFactoryPipelineAndTriggerName(),
		},

		"data_factory_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: factories.ValidateFactoryID,
		},
	}
}

func (d TriggerScheduleDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"description": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"schedule": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"days_of_month": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeInt,
						},
					},

					"days_of_week": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"hours": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeInt,
						},
					},

					"minutes": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeInt,
						},
					},

					"monthly": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"weekday": {
									Type:     pluginsdk.TypeString,
									Computed: true,
								},

								"week": {
									Type:     pluginsdk.TypeInt,
									Computed: true,
								},
							},
						},
					},
				},
			},
		},

		"start_time": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"end_time": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"time_zone": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"frequency": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"interval": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"activated": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"pipeline_name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"annotations": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (d TriggerScheduleDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model TriggerScheduleDataSourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.DataFactory.TriggersClient

			dataFactoryId, err := factories.ParseFactoryID(model.DataFactoryID)
			if err != nil {
				return err
			}

			id := parse.NewTriggerID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, model.Name)
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.Name, "")
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retreiving %s: %+v", id, err)
			}

			metadata.SetID(id)

			model.Name = *existing.Name
			model.DataFactoryID = dataFactoryId.ID()

			scheduleTriggerProps, ok := existing.Properties.AsScheduleTrigger()
			if !ok {
				return fmt.Errorf("classifying %s: Expected: %q Received: %q", id.ID(), datafactory.TypeBasicTriggerTypeScheduleTrigger, *existing.Type)
			}

			if scheduleTriggerProps != nil {
				model.Activated = scheduleTriggerProps.RuntimeState == datafactory.TriggerRuntimeStateStarted

				if recurrence := scheduleTriggerProps.Recurrence; recurrence != nil {
					if v := recurrence.StartTime; v != nil {
						model.StartTime = v.Format(time.RFC3339)
					}
					if v := recurrence.EndTime; v != nil {
						model.EndTime = v.Format(time.RFC3339)
					}
					model.Frequency = string(recurrence.Frequency)
					model.Interval = int64(*recurrence.Interval)
					model.TimeZone = *recurrence.TimeZone

					if schedule := recurrence.Schedule; schedule != nil {
						model.Schedule = flattenDataFactoryScheduleModel(schedule)
					}
				}

				if pipelines := scheduleTriggerProps.Pipelines; pipelines != nil {
					if len(*pipelines) > 0 {
						pipeline := *pipelines
						if reference := pipeline[0].PipelineReference; reference != nil {
							model.PipelineName = *reference.ReferenceName
						}
					}
				}

				model.Annotations = flattenDataFactoryAnnotations(scheduleTriggerProps.Annotations)

				if scheduleTriggerProps.Description != nil {
					model.Description = *scheduleTriggerProps.Description
				}

			}
			if err := metadata.Encode(&model); err != nil {
				return fmt.Errorf("encoding: %+v", err)
			}

			return metadata.Encode(&model)
		},
	}
}

func flattenDataFactoryScheduleModel(schedule *datafactory.RecurrenceSchedule) []TriggerSchedule {
	if schedule == nil {
		return []TriggerSchedule{}
	}

	result := TriggerSchedule{}
	results := []TriggerSchedule{}

	if schedule.Hours != nil {
		for _, v := range *schedule.Hours {
			result.Hours = append(result.Hours, int64(v))
		}
	}

	if schedule.Minutes != nil {
		for _, v := range *schedule.Minutes {
			result.Minutes = append(result.Minutes, int64(v))
		}
	}

	if schedule.MonthDays != nil {
		for _, v := range *schedule.MonthDays {
			result.DaysOfMonth = append(result.DaysOfMonth, int64(v))
		}
	}

	if schedule.WeekDays != nil {
		weekDays := make([]string, 0)
		for _, v := range *schedule.WeekDays {
			weekDays = append(weekDays, string(v))
		}
		result.DaysOfWeek = weekDays
	}

	if schedule.MonthlyOccurrences != nil {
		var monthlyOccurrences []TriggerScheduleScheduleMonthly
		for _, v := range *schedule.MonthlyOccurrences {
			occurrence := TriggerScheduleScheduleMonthly{}
			occurrence.Weekday = string(v.Day)
			if v.Occurrence != nil {
				occurrence.Week = int64(*v.Occurrence)
			}
			monthlyOccurrences = append(monthlyOccurrences, occurrence)
		}
		result.Monthly = monthlyOccurrences
	}
	results = append(results, result)
	return results
}
