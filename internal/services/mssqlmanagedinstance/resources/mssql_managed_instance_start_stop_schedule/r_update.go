// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_instance_start_stop_schedule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/parse"
)

func (r Resource) Update() sdk.ResourceFunc {
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
				return fmt.Errorf("retrieving %s: `model` was nil", managedInstanceID)
			}
			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", managedInstanceID)
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
