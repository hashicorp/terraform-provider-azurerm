// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package batch_job

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/data-plane/batch/2022-01-01-15-0/jobs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/batchaccount"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/parse"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config ResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

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
			existing, err := client.JobGet(ctx, idSDK, jobs.DefaultJobGetOperationOptions())
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}
			model := existing.Model

			payload := jobs.JobUpdateParameter{
				AllowTaskPreemption: model.AllowTaskPreemption,
				Constraints:         model.Constraints,
				MaxParallelTasks:    model.MaxParallelTasks,
				Metadata:            model.Metadata,
				OnAllTasksComplete:  model.OnAllTasksComplete,
				PoolInfo:            pointer.From(model.PoolInfo),
				Priority:            model.Priority,
			}

			if metadata.ResourceData.HasChange("priority") {
				payload.Priority = pointer.To(config.Priority)
			}

			if metadata.ResourceData.HasChange("task_retry_maximum") {
				if payload.Constraints == nil {
					payload.Constraints = new(jobs.JobConstraints)
				}
				payload.Constraints.MaxTaskRetryCount = pointer.To(config.TaskRetryMaximum)
			}

			if _, err := client.JobUpdate(ctx, idSDK, payload, jobs.DefaultJobUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}
