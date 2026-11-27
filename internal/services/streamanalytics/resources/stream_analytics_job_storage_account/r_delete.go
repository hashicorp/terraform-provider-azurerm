// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_job_storage_account

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/streamingjobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	streamanalyticsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/client"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := streamanalyticsClient.From(metadata.Client).JobsClient

			id, err := streamingjobs.ParseStreamingJobID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			existing, err := client.Get(ctx, *id, streamingjobs.DefaultGetOperationOptions())
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}
			if existing.Model.Properties.JobStorageAccount == nil {
				return fmt.Errorf("retrieving %s: `properties.JobStorage` was nil", id)
			}

			payload := existing.Model

			// We're unable to remove the Job Storage Account using the PATCH call, the only way to remove it is by using the PUT
			payload.Properties.JobStorageAccount = nil

			if err := client.CreateOrReplaceThenPoll(ctx, *id, *payload, streamingjobs.DefaultCreateOrReplaceOperationOptions()); err != nil {
				return fmt.Errorf("deleting Job Storage Account for %s: %+v", *id, err)
			}

			return nil
		},
	}
}
