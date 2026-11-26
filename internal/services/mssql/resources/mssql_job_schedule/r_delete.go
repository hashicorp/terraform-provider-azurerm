// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_job_schedule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/jobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (Resource) Delete() sdk.ResourceFunc {
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

			existing, err := client.Get(ctx, *jobId)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return nil
				}

				return fmt.Errorf("retrieving %s: %+v", jobId, err)
			}

			if model := existing.Model; model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", jobId)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", jobId)
			}

			// Set Schedule to nil allowing Azure to repopulate with default values
			existing.Model.Properties.Schedule = nil
			if _, err := client.CreateOrUpdate(ctx, *jobId, *existing.Model); err != nil {
				return fmt.Errorf("deleting schedule from %s, %+v", jobId, err)
			}

			return nil
		},
	}
}
