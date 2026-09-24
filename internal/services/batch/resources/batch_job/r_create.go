// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package batch_job

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/data-plane/batch/2022-01-01-15-0/jobs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/batchaccount"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/pool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/parse"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model BatchJobModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			poolId, err := pool.ParsePoolID(model.BatchPoolId)
			if err != nil {
				return err
			}

			endpoint, err := r.GetEndpoint(ctx, metadata.Client.Batch.AccountClient, batchaccount.NewBatchAccountID(poolId.SubscriptionId, poolId.ResourceGroupName, poolId.BatchAccountName))
			if err != nil {
				return err
			}
			client := metadata.Client.Batch.JobsClient.Clone(endpoint)

			id := parse.NewJobID(poolId.SubscriptionId, poolId.ResourceGroupName, poolId.BatchAccountName, poolId.PoolName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				idSDK := jobs.NewJobID(endpoint, model.Name)
				existing, err := client.JobGet(ctx, idSDK, jobs.DefaultJobGetOperationOptions())
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			params := jobs.JobAddParameter{
				Id:          model.Name,
				DisplayName: &model.DisplayName,
				Priority:    new(model.Priority),
				Constraints: &jobs.JobConstraints{
					MaxTaskRetryCount: new(model.TaskRetryMaximum),
				},
				CommonEnvironmentSettings: r.expandEnvironmentSettings(model.CommonEnvironmentProperties),
				PoolInfo: jobs.PoolInformation{
					PoolId: &poolId.PoolName,
				},
			}

			if _, err := client.JobAdd(ctx, params, jobs.DefaultJobAddOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			// TODO: should this be migrated to the data plane ID?
			metadata.SetID(id)
			return nil
		},
	}
}

func (r Resource) expandEnvironmentSettings(input map[string]string) *[]jobs.EnvironmentSetting {
	if len(input) == 0 {
		return nil
	}

	m := make([]jobs.EnvironmentSetting, 0, len(input))
	for k, v := range input {
		m = append(m, jobs.EnvironmentSetting{
			Name:  k,
			Value: new(v),
		})
	}
	return &m
}
