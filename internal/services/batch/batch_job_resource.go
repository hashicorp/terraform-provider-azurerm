// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package batch

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
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type BatchJobResource struct{}

var _ sdk.ResourceWithUpdate = BatchJobResource{}

type BatchJobModel struct {
	Name                        string            `tfschema:"name"`
	BatchPoolId                 string            `tfschema:"batch_pool_id"`
	DisplayName                 string            `tfschema:"display_name"`
	Priority                    int64             `tfschema:"priority"`
	TaskRetryMaximum            int64             `tfschema:"task_retry_maximum"`
	CommonEnvironmentProperties map[string]string `tfschema:"common_environment_properties"`
}

func (r BatchJobResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.JobName,
		},
		"batch_pool_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: pool.ValidatePoolID,
		},
		"display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"common_environment_properties": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
		"priority": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      0,
			ValidateFunc: validation.IntBetween(-1000, 1000),
		},
		"task_retry_maximum": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(-1),
		},
	}
}

func (r BatchJobResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r BatchJobResource) ResourceType() string {
	return "azurerm_batch_job"
}

func (r BatchJobResource) ModelObject() interface{} {
	return &BatchJobModel{}
}

func (r BatchJobResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.JobID
}

func (r BatchJobResource) GetEndpoint(ctx context.Context, client *batchaccount.BatchAccountClient, accountID batchaccount.BatchAccountId) (string, error) {
	account, err := client.Get(ctx, accountID)
	if err != nil {
		return "", fmt.Errorf("retrieving %s: %v", accountID, err)
	}

	endpoint := ""
	if account.Model != nil && account.Model.Properties != nil {
		endpoint = "https://" + *account.Model.Properties.AccountEndpoint
	}

	if endpoint == "" {
		return "", fmt.Errorf("retrieving %s: unable to determine account data plane endpoint", accountID)
	}

	return endpoint, nil
}

func (r BatchJobResource) Create() sdk.ResourceFunc {
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
				Priority:    pointer.To(model.Priority),
				Constraints: &jobs.JobConstraints{
					MaxTaskRetryCount: pointer.To(model.TaskRetryMaximum),
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

func (r BatchJobResource) Read() sdk.ResourceFunc {
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

			state := BatchJobModel{
				Name:        id.Name,
				BatchPoolId: pool.NewPoolID(id.SubscriptionId, id.ResourceGroup, id.BatchAccountName, id.PoolName).ID(),
			}

			if model := resp.Model; model != nil {
				state.CommonEnvironmentProperties = r.flattenEnvironmentSettings(model.CommonEnvironmentSettings)
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

func (r BatchJobResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config BatchJobModel
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

func (r BatchJobResource) Delete() sdk.ResourceFunc {
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

func (r BatchJobResource) expandEnvironmentSettings(input map[string]string) *[]jobs.EnvironmentSetting {
	if len(input) == 0 {
		return nil
	}

	m := make([]jobs.EnvironmentSetting, 0, len(input))
	for k, v := range input {
		m = append(m, jobs.EnvironmentSetting{
			Name:  k,
			Value: pointer.To(v),
		})
	}
	return &m
}

func (r BatchJobResource) flattenEnvironmentSettings(input *[]jobs.EnvironmentSetting) map[string]string {
	if input == nil {
		return nil
	}

	m := make(map[string]string)
	for _, setting := range *input {
		m[setting.Name] = pointer.From(setting.Value)
	}
	return m
}
