// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_job_storage_account

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/streamingjobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	streamanalyticsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/client"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := streamanalyticsClient.From(metadata.Client).JobsClient

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

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, *id, streamingjobs.DefaultGetOperationOptions())
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				jobStorageAccountExists := existing.Model != nil && existing.Model.Properties != nil && existing.Model.Properties.JobStorageAccount != nil

				if !response.WasNotFound(existing.HttpResponse) && jobStorageAccountExists {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			payload := streamingjobs.StreamingJob{
				Properties: &streamingjobs.StreamingJobProperties{
					JobStorageAccount: &streamingjobs.JobStorageAccount{
						AccountName:        pointer.To(model.StorageAccountName),
						AuthenticationMode: pointer.ToEnum[streamingjobs.AuthenticationMode](model.AuthenticationMode),
					},
				},
			}

			if model.AuthenticationMode == string(streamingjobs.AuthenticationModeMsi) && model.StorageAccountKey != "" {
				return errors.New("`storage_account_key` cannot be set if `authentication_mode` is `Msi`")
			}

			if model.AuthenticationMode == string(streamingjobs.AuthenticationModeConnectionString) && model.StorageAccountKey == "" {
				return errors.New("`storage_account_key` cannot be empty if `authentication_mode` is `ConnectionString`")
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
