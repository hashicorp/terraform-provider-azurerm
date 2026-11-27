// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_job_storage_account

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/streamingjobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	streamanalyticsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/client"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := streamanalyticsClient.From(metadata.Client).JobsClient

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
				payload.Properties.JobStorageAccount.AuthenticationMode = pointer.ToEnum[streamingjobs.AuthenticationMode](model.AuthenticationMode)
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
