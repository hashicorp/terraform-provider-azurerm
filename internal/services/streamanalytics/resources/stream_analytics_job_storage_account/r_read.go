// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_job_storage_account

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/streamingjobs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	streamanalyticsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/client"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := streamanalyticsClient.From(metadata.Client).JobsClient

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
						state.AuthenticationMode = pointer.FromEnum(jobStorage.AuthenticationMode)
						state.StorageAccountKey = metadata.ResourceData.Get("storage_account_key").(string)
						state.StorageAccountName = pointer.From(jobStorage.AccountName)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}
