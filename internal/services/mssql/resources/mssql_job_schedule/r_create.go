// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_job_schedule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/jobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.JobsClient

			var config MsSqlJobScheduleResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			jobId, err := jobs.ParseJobID(config.JobID)
			if err != nil {
				return err
			}

			locks.ByID(jobId.ID())
			defer locks.UnlockByID(jobId.ID())

			existing, err := client.Get(ctx, *jobId)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s was not found: %+v", jobId, err)
				}

				return fmt.Errorf("checking for presence of existing %s: %+v", jobId, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", jobId)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `model.Properties` was nil", jobId)
			}

			if existing.Model.Properties.Schedule == nil {
				return fmt.Errorf("retrieving %s: `model.Properties.Schedule` was nil", jobId)
			}

			// Default schedule is disabled when created using the API
			// if schedule is enabled we can reasonably assume the schedule was modified outside of Terraform and should be imported.
			schedule := existing.Model.Properties.Schedule
			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				if pointer.From(schedule.Enabled) {
					return metadata.ResourceRequiresImport(r.ResourceType(), jobId)
				}
			}

			schedule.Enabled = pointer.To(config.Enabled)
			schedule.Type = pointer.ToEnum[jobs.JobScheduleType](config.Type)

			if config.EndTime != "" {
				schedule.EndTime = pointer.To(config.EndTime)
			}
			if config.Interval != "" && config.Type == string(jobs.JobScheduleTypeRecurring) {
				schedule.Interval = pointer.To(config.Interval)
			}
			if config.StartTime != "" {
				schedule.StartTime = pointer.To(config.StartTime)
			}

			if _, err := client.CreateOrUpdate(ctx, *jobId, *existing.Model); err != nil {
				return fmt.Errorf("creating schedule for %s: %+v", jobId, err)
			}

			metadata.SetID(jobId)
			return nil
		},
	}
}
