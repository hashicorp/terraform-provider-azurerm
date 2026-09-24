// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package batch_job

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/data-plane/batch/2022-01-01-15-0/jobs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/batchaccount"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/parse"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.JobID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			endpoint, err := r.GetEndpoint(ctx, metadata.Client.Batch.AccountClient, batchaccount.NewBatchAccountID(id.SubscriptionId, id.ResourceGroup, id.BatchAccountName))
			if err != nil {
				return err
			}

			client := metadata.Client.Batch.JobsClient.Clone(endpoint)

			idSDK := jobs.NewJobID(endpoint, id.Name)
			if _, err := client.JobDelete(ctx, idSDK, jobs.DefaultJobDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
