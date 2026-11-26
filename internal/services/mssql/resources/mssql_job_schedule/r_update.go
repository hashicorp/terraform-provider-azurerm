// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_job_schedule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/jobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.JobsClient

			jobId, err := jobs.ParseJobID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(jobId.ID())
			defer locks.UnlockByID(jobId.ID())

			var config MsSqlJobScheduleResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *jobId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", jobId, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", jobId)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", jobId)
			}

			if existing.Model.Properties.Schedule == nil {
				return fmt.Errorf("retrieving %s: `schedule` was nil", jobId)
			}

			schedule := existing.Model.Properties.Schedule
			if metadata.ResourceData.HasChange("enabled") {
				schedule.Enabled = pointer.To(config.Enabled)
			}

			if metadata.ResourceData.HasChange("end_time") {
				schedule.EndTime = pointer.To(config.EndTime)
			}

			if metadata.ResourceData.HasChange("interval") {
				schedule.Interval = pointer.To(config.Interval)
			}

			if metadata.ResourceData.HasChange("start_time") {
				schedule.StartTime = pointer.To(config.StartTime)
			}

			if metadata.ResourceData.HasChange("type") {
				schedule.Type = pointer.ToEnum[jobs.JobScheduleType](config.Type)
			}

			if _, err := client.CreateOrUpdate(ctx, *jobId, *existing.Model); err != nil {
				return fmt.Errorf("updating schedule for %s: %+v", jobId, err)
			}

			return nil
		},
	}
}
