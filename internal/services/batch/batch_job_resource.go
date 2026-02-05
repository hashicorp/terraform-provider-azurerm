// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package batch

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/data-plane/batch/2022-01-01-15-0/jobs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/batchaccount"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/pool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	batchDataplane "github.com/jackofallops/kermit/sdk/batch/2022-01.15.0/batch"
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
	return jobs.ValidateJobID
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

			accountId := batchaccount.NewBatchAccountID(poolId.SubscriptionId, poolId.ResourceGroupName, poolId.BatchAccountName)

			account, err := metadata.Client.Batch.AccountClient.Get(ctx, accountId)
			if err != nil || account.Model == nil {
				return err
			}

			loc := location.Normalize(account.Model.Location)

			id := jobs.NewJobID(fmt.Sprintf(dataPlaneEndpointFmt, poolId.BatchAccountName, loc), model.Name)
			client := metadata.Client.Batch.JobsDataPlaneClient.Clone(id.BaseURI)

			client.JobsClientSetEndpoint(id.BaseURI)

			existing, err := client.JobGet(ctx, id, jobs.DefaultJobGetOperationOptions())
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
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

			if _, err = client.JobAdd(ctx, params, jobs.DefaultJobAddOperationOptions()); err != nil {
				return err
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r BatchJobResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			data := BatchJobModel{}

			if err := metadata.Decode(&data); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			id, err := jobs.ParseJobID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// client, err := metadata.Client.Batch.JobClient(ctx, accountId)
			client := metadata.Client.Batch.JobsDataPlaneClient.Clone(id.BaseURI)

			resp, err := client.JobGet(ctx, *id, jobs.DefaultJobGetOperationOptions())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			data.Name = id.JobId

			if model := resp.Model; model != nil {
				data.Priority = pointer.From(model.Priority)
				data.DisplayName = pointer.From(model.DisplayName)
				if constraints := model.Constraints; constraints != nil {
					data.TaskRetryMaximum = pointer.From(constraints.MaxTaskRetryCount)
				}

				data.CommonEnvironmentProperties = r.flattenEnvironmentSettings(model.CommonEnvironmentSettings)

				// retrieve `batch_pool_id` from config if we can, will fail for import
				data.BatchPoolId = metadata.ResourceData.Get("batch_pool_id").(string)
			}

			return metadata.Encode(&data)
		},
	}
}

func (r BatchJobResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := jobs.ParseJobID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}
			client := metadata.Client.Batch.JobsDataPlaneClient.Clone(id.BaseURI)

			var model BatchJobModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			patch := jobs.JobPatchParameter{}

			if metadata.ResourceData.HasChange("priority") {
				patch.Priority = pointer.To(model.Priority)
			}

			if metadata.ResourceData.HasChange("task_retry_maximum") {
				if patch.Constraints == nil {
					patch.Constraints = new(jobs.JobConstraints)
				}
				patch.Constraints.MaxTaskRetryCount = pointer.To(model.TaskRetryMaximum)
			}

			if _, err = client.JobPatch(ctx, *id, patch, jobs.DefaultJobPatchOperationOptions()); err != nil {
				return err
			}

			return nil
		},
	}
}

func (r BatchJobResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := jobs.ParseJobID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			client := metadata.Client.Batch.JobsDataPlaneClient.Clone(id.BaseURI)

			if _, err := client.JobDelete(ctx, *id, jobs.DefaultJobDeleteOperationOptions()); err != nil {
				return err
			}

			return nil
		},
	}
}

func (r BatchJobResource) addJob(ctx context.Context, client *batchDataplane.JobClient, id parse.JobId, job batchDataplane.JobAddParameter) error {
	deadline, _ := ctx.Deadline()
	now := time.Now()
	timeout := deadline.Sub(now)
	_, err := client.Add(ctx, job, pointer.To(int32(timeout.Seconds())), nil, nil, &date.TimeRFC1123{Time: now})
	if err != nil {
		return fmt.Errorf("creating %s: %v", id, err)
	}
	return nil
}

func (r BatchJobResource) getJob(ctx context.Context, client *batchDataplane.JobClient, id parse.JobId) (batchDataplane.CloudJob, error) {
	deadline, _ := ctx.Deadline()
	now := time.Now()
	timeout := deadline.Sub(now)
	return client.Get(ctx, id.Name, "", "", pointer.To(int32(timeout.Seconds())), nil, nil, &date.TimeRFC1123{Time: now}, "", "", nil, nil)
}

func (r BatchJobResource) patchJob(ctx context.Context, client *batchDataplane.JobClient, id parse.JobId, job batchDataplane.JobPatchParameter) error {
	deadline, _ := ctx.Deadline()
	now := time.Now()
	timeout := deadline.Sub(now)
	_, err := client.Patch(ctx, id.Name, job, pointer.To(int32(timeout.Seconds())), nil, nil, &date.TimeRFC1123{Time: now}, "", "", nil, nil)
	return err
}

func (r BatchJobResource) deleteJob(ctx context.Context, client *batchDataplane.JobClient, id parse.JobId) error {
	deadline, _ := ctx.Deadline()
	now := time.Now()
	timeout := deadline.Sub(now)
	_, err := client.Delete(ctx, id.Name, pointer.To(int32(timeout.Seconds())), nil, nil, &date.TimeRFC1123{Time: now}, "", "", nil, nil)
	return err
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
		if setting.Value == nil {
			continue
		}
		m[setting.Name] = *setting.Value
	}
	return m
}
