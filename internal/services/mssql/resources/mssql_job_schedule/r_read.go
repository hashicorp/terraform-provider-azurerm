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
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQL.JobsClient

			jobId, err := jobs.ParseJobID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *jobId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(jobId)
				}

				return fmt.Errorf("retrieving %s: %+v", jobId, err)
			}

			state := MsSqlJobScheduleResourceModel{
				JobID: jobId.ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					if schedule := props.Schedule; schedule != nil {
						state.Enabled = pointer.From(schedule.Enabled)
						state.EndTime = pointer.From(schedule.EndTime)
						state.Interval = pointer.From(schedule.Interval)
						state.StartTime = pointer.From(schedule.StartTime)
						state.Type = pointer.FromEnum(schedule.Type)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}
