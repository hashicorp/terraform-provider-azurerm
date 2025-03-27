// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package streamanalytics

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/streamingjobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type JobStorageAccountResource struct{}

type JobStorageAccountModel struct {
	JobId              string `tfschema:"stream_analytics_job_id"`
	AuthenticationMode string `tfschema:"authentication_mode"`
	StorageAccountKey  string `tfschema:"storage_account_key"`
	StorageAccountName string `tfschema:"storage_account_name"`
}

var _ sdk.ResourceWithUpdate = JobStorageAccountResource{}

func (r JobStorageAccountResource) ModelObject() interface{} {
	return &JobStorageAccountModel{}
}

func (r JobStorageAccountResource) ResourceType() string {
	return "azurerm_stream_analytics_job_storage_account"
}

func (r JobStorageAccountResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return streamingjobs.ValidateStreamingJobID
}

func (r JobStorageAccountResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"stream_analytics_job_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: streamingjobs.ValidateStreamingJobID,
		},

		"authentication_mode": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(streamingjobs.AuthenticationModeMsi),
				string(streamingjobs.AuthenticationModeConnectionString),
				// auth mode `UserToken` is not supported
			}, false),
		},

		"storage_account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"storage_account_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r JobStorageAccountResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r JobStorageAccountResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.JobsClient

			var model JobStorageAccountModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := streamingjobs.ParseStreamingJobID(model.JobId)
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			existing, err := client.Get(ctx, *id, streamingjobs.DefaultGetOperationOptions())
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			jobStorageAccountExists := existing.Model != nil && existing.Model.Properties != nil && existing.Model.Properties.JobStorageAccount != nil

			if !response.WasNotFound(existing.HttpResponse) && jobStorageAccountExists {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := streamingjobs.StreamingJob{
				Properties: &streamingjobs.StreamingJobProperties{
					JobStorageAccount: &streamingjobs.JobStorageAccount{
						AccountName:        pointer.To(model.StorageAccountName),
						AuthenticationMode: pointer.To(streamingjobs.AuthenticationMode(model.AuthenticationMode)),
					},
				},
			}

			if model.AuthenticationMode == string(streamingjobs.AuthenticationModeMsi) && model.StorageAccountKey != "" {
				return fmt.Errorf("`storage_account_key` cannot be set if `authentication_mode` is `Msi`")
			}

			if model.AuthenticationMode == string(streamingjobs.AuthenticationModeConnectionString) && model.StorageAccountKey == "" {
				return fmt.Errorf("`storage_account_key` cannot be empty if `authentication_mode` is `ConnectionString`")
			}

			if model.StorageAccountKey != "" {
				payload.Properties.JobStorageAccount.AccountKey = pointer.To(model.StorageAccountKey)
			}

			if _, err := client.Update(ctx, *id, payload, streamingjobs.DefaultUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating Job Storage Account for %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r JobStorageAccountResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.JobsClient

			id, err := streamingjobs.ParseStreamingJobID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id, streamingjobs.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) || resp.Model == nil || resp.Model.Properties == nil || resp.Model.Properties.JobStorageAccount == nil {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := JobStorageAccountModel{
				JobId: id.ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					if jobStorage := props.JobStorageAccount; jobStorage != nil {
						state.AuthenticationMode = string(pointer.From(jobStorage.AuthenticationMode))
						state.StorageAccountKey = metadata.ResourceData.Get("storage_account_key").(string)
						state.StorageAccountName = pointer.From(jobStorage.AccountName)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r JobStorageAccountResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.JobsClient

			var model JobStorageAccountModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := streamingjobs.ParseStreamingJobID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			payload := streamingjobs.StreamingJob{
				Properties: &streamingjobs.StreamingJobProperties{
					JobStorageAccount: &streamingjobs.JobStorageAccount{},
				},
			}

			if metadata.ResourceData.HasChange("authentication_mode") {
				payload.Properties.JobStorageAccount.AuthenticationMode = pointer.To(streamingjobs.AuthenticationMode(model.AuthenticationMode))
			}

			if metadata.ResourceData.HasChange("storage_account_name") {
				payload.Properties.JobStorageAccount.AccountName = pointer.To(model.StorageAccountName)
			}

			if metadata.ResourceData.HasChange("storage_account_key") {
				if model.StorageAccountKey != "" {
					payload.Properties.JobStorageAccount.AccountKey = pointer.To(model.StorageAccountKey)
				}
			}

			if _, err := client.Update(ctx, *id, payload, streamingjobs.DefaultUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating Job Storage Account for %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r JobStorageAccountResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StreamAnalytics.JobsClient

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
