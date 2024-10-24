package mssqlmanagedinstance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	schedule "github.com/hashicorp/go-azure-sdk/resource-manager/sql/2023-08-01-preview/startstopmanagedinstanceschedules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SqlManagedInstanceStartStopScheduleModel struct {
	SqlManagedInstanceId string              `tfschema:"managed_instance_id"`
	Description          string              `tfschema:"description"`
	Schedule             []ScheduleItemModel `tfschema:"schedule"`
	TimeZoneId           string              `tfschema:"timezone_id"`
	NextExecutionTime    string              `tfschema:"next_execution_time"`
	NextRunAction        string              `tfschema:"next_run_action"`
}

type ScheduleItemModel struct {
	StartDay  schedule.DayOfWeek `tfschema:"start_day"`
	StartTime string             `tfschema:"start_time"`
	StopDay   schedule.DayOfWeek `tfschema:"stop_day"`
	StopTime  string             `tfschema:"stop_time"`
}

type MsSqlManagedInstanceStartStopScheduleResource struct{}

var _ sdk.ResourceWithUpdate = MsSqlManagedInstanceStartStopScheduleResource{}

func (r MsSqlManagedInstanceStartStopScheduleResource) ResourceType() string {
	return "azurerm_mssql_managed_instance_start_stop_schedule"
}

func (r MsSqlManagedInstanceStartStopScheduleResource) ModelObject() interface{} {
	return &SqlManagedInstanceStartStopScheduleModel{}
}

func (r MsSqlManagedInstanceStartStopScheduleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.ManagedInstanceStartStopScheduleID
}

func (r MsSqlManagedInstanceStartStopScheduleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"managed_instance_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateSqlManagedInstanceID,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"schedule": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"start_day": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(schedule.PossibleValuesForDayOfWeek(), false),
					},

					"start_time": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"stop_day": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(schedule.PossibleValuesForDayOfWeek(), false),
					},

					"stop_time": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"timezone_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "UTC",
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r MsSqlManagedInstanceStartStopScheduleResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"next_execution_time": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"next_run_action": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r MsSqlManagedInstanceStartStopScheduleResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SqlManagedInstanceStartStopScheduleModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MSSQLManagedInstance.ManagedInstanceStartStopSchedulesClient

			managedInstanceId, err := commonids.ParseSqlManagedInstanceID(model.SqlManagedInstanceId)
			if err != nil {
				return err
			}

			id := *managedInstanceId

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &schedule.StartStopManagedInstanceSchedule{
				Properties: &schedule.StartStopManagedInstanceScheduleProperties{},
			}

			if model.Description != "" {
				properties.Properties.Description = &model.Description
			}

			properties.Properties.ScheduleList = expandScheduleItemModelArray(model.Schedule)

			if model.TimeZoneId != "" {
				properties.Properties.TimeZoneId = &model.TimeZoneId
			}

			if _, err := client.CreateOrUpdate(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			scheduleID := parse.NewManagedInstanceStartStopScheduleID(id.SubscriptionId, id.ResourceGroupName, id.ManagedInstanceName, "default")
			metadata.SetID(scheduleID)

			return nil
		},
	}
}

func (r MsSqlManagedInstanceStartStopScheduleResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQLManagedInstance.ManagedInstanceStartStopSchedulesClient

			id, err := parse.ManagedInstanceStartStopScheduleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			managedInstanceID := commonids.NewSqlManagedInstanceID(id.SubscriptionId, id.ResourceGroup, id.ManagedInstanceName)

			var model SqlManagedInstanceStartStopScheduleModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, managedInstanceID)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", managedInstanceID, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s:`model` was nil", managedInstanceID)
			}
			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s:`properties` was nil", managedInstanceID)
			}

			properties := resp.Model

			if metadata.ResourceData.HasChange("description") {
				properties.Properties.Description = pointer.To(model.Description)
			}

			if metadata.ResourceData.HasChange("schedule") {
				properties.Properties.ScheduleList = expandScheduleItemModelArray(model.Schedule)
			}

			if metadata.ResourceData.HasChange("timezone_id") {
				properties.Properties.TimeZoneId = pointer.To(model.TimeZoneId)
			}

			if _, err := client.CreateOrUpdate(ctx, managedInstanceID, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", managedInstanceID, err)
			}

			return nil
		},
	}
}

func (r MsSqlManagedInstanceStartStopScheduleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQLManagedInstance.ManagedInstanceStartStopSchedulesClient

			id, err := parse.ManagedInstanceStartStopScheduleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			managedInstanceID := commonids.NewSqlManagedInstanceID(id.SubscriptionId, id.ResourceGroup, id.ManagedInstanceName)

			resp, err := client.Get(ctx, managedInstanceID)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(managedInstanceID)
				}

				return fmt.Errorf("retrieving %s: %+v", managedInstanceID, err)
			}

			state := SqlManagedInstanceStartStopScheduleModel{
				SqlManagedInstanceId: managedInstanceID.ID(),
			}

			if model := resp.Model; model != nil {
				if properties := model.Properties; properties != nil {
					state.Description = pointer.From(properties.Description)

					state.NextExecutionTime = pointer.From(properties.NextExecutionTime)

					state.NextRunAction = pointer.From(properties.NextRunAction)

					if properties.ScheduleList != nil {
						state.Schedule = flattenScheduleItemModelArray(properties.ScheduleList)
					}

					state.TimeZoneId = pointer.From(properties.TimeZoneId)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r MsSqlManagedInstanceStartStopScheduleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQLManagedInstance.ManagedInstanceStartStopSchedulesClient

			id, err := parse.ManagedInstanceStartStopScheduleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			managedInstanceID := commonids.NewSqlManagedInstanceID(id.SubscriptionId, id.ResourceGroup, id.ManagedInstanceName)

			if _, err := client.Delete(ctx, managedInstanceID); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandScheduleItemModelArray(inputList []ScheduleItemModel) []schedule.ScheduleItem {
	var outputList []schedule.ScheduleItem
	for _, v := range inputList {
		input := v
		output := schedule.ScheduleItem{
			StartDay:  input.StartDay,
			StartTime: input.StartTime,
			StopDay:   input.StopDay,
			StopTime:  input.StopTime,
		}

		outputList = append(outputList, output)
	}
	return outputList
}

func flattenScheduleItemModelArray(inputList []schedule.ScheduleItem) []ScheduleItemModel {
	var outputList []ScheduleItemModel
	if inputList == nil {
		return outputList
	}
	for _, input := range inputList {
		output := ScheduleItemModel{
			StartDay:  input.StartDay,
			StartTime: input.StartTime,
			StopDay:   input.StopDay,
			StopTime:  input.StopTime,
		}

		outputList = append(outputList, output)
	}
	return outputList
}
