// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package job

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/data-plane/batch/2022-01-01-15-0/jobs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/batchaccount"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/pool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/parse"
)

func (r Resource) Read() sdk.ResourceFunc {
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
			resp, err := client.JobGet(ctx, idSDK, jobs.DefaultJobGetOperationOptions())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := ResourceModel{
				Name:        id.Name,
				BatchPoolId: pool.NewPoolID(id.SubscriptionId, id.ResourceGroup, id.BatchAccountName, id.PoolName).ID(),
			}

			if model := resp.Model; model != nil {
				state.CommonEnvironmentProperties = flattenEnvironmentSettings(model.CommonEnvironmentSettings)
				state.DisplayName = pointer.From(model.DisplayName)
				state.Priority = pointer.From(model.Priority)

				if constraints := model.Constraints; constraints != nil {
					state.TaskRetryMaximum = pointer.From(constraints.MaxTaskRetryCount)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func flattenEnvironmentSettings(input *[]jobs.EnvironmentSetting) map[string]string {
	if input == nil {
		return map[string]string{}
	}

	m := make(map[string]string)
	for _, setting := range *input {
		m[setting.Name] = pointer.From(setting.Value)
	}
	return m
}
