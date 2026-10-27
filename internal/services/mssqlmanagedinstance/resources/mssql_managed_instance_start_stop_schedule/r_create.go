// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_instance_start_stop_schedule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/startstopmanagedinstanceschedules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/parse"
)

func (r Resource) Create() sdk.ResourceFunc {
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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			properties := &startstopmanagedinstanceschedules.StartStopManagedInstanceSchedule{
				Properties: &startstopmanagedinstanceschedules.StartStopManagedInstanceScheduleProperties{},
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

			metadata.SetID(parse.NewManagedInstanceStartStopScheduleID(id.SubscriptionId, id.ResourceGroupName, id.ManagedInstanceName, "default"))

			return nil
		},
	}
}

func expandScheduleItemModelArray(inputList []ScheduleItemModel) []startstopmanagedinstanceschedules.ScheduleItem {
	outputList := make([]startstopmanagedinstanceschedules.ScheduleItem, 0, len(inputList))

	for _, input := range inputList {
		output := startstopmanagedinstanceschedules.ScheduleItem{
			StartDay:  input.StartDay,
			StartTime: input.StartTime,
			StopDay:   input.StopDay,
			StopTime:  input.StopTime,
		}

		outputList = append(outputList, output)
	}
	return outputList
}
