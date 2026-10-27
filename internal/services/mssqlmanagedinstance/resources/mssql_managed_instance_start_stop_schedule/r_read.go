// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_instance_start_stop_schedule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/startstopmanagedinstanceschedules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/parse"
)

func (r Resource) Read() sdk.ResourceFunc {
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

func flattenScheduleItemModelArray(inputList []startstopmanagedinstanceschedules.ScheduleItem) []ScheduleItemModel {
	outputList := make([]ScheduleItemModel, 0, len(inputList))

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
