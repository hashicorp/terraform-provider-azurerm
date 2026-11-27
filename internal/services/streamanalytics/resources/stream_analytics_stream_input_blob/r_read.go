// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_stream_input_blob

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/inputs"
	clientsAccount "github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	streamanalyticsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStreamAnalyticsStreamInputBlobRead(d *pluginsdk.ResourceData, meta any) error {
	client := streamanalyticsClient.From(meta).InputsClient
	ctx, cancel := timeouts.ForRead(clientsAccount.StopContext(meta), d)
	defer cancel()

	id, err := inputs.ParseInputID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.InputName)
	d.Set("stream_analytics_job_name", id.StreamingJobName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			streamInput, ok := props.(inputs.StreamInputProperties)
			if !ok {
				return fmt.Errorf("converting %s to a Stream Input", *id)
			}

			streamBlobInput, ok := streamInput.Datasource.(inputs.BlobStreamInputDataSource)
			if !ok {
				return fmt.Errorf("converting Stream Input Blob to an Stream Input: %+v", err)
			}

			if streamBlobInputProps := streamBlobInput.Properties; streamBlobInputProps != nil {
				d.Set("date_format", pointer.From(streamBlobInput.Properties.DateFormat))

				d.Set("path_pattern", pointer.From(streamBlobInputProps.PathPattern))

				d.Set("storage_container_name", pointer.From(streamBlobInputProps.Container))

				d.Set("time_format", pointer.From(streamBlobInputProps.TimeFormat))

				authMode := ""
				if v := streamBlobInputProps.AuthenticationMode; v != nil {
					authMode = string(*v)
				}
				d.Set("authentication_mode", authMode)

				if accounts := streamBlobInputProps.StorageAccounts; accounts != nil && len(*accounts) > 0 {
					account := (*accounts)[0]
					d.Set("storage_account_name", account.AccountName)
				}

				if err := d.Set("serialization", helpers.FlattenStreamAnalyticsStreamInputSerialization(streamInput.Serialization)); err != nil {
					return fmt.Errorf("setting `serialization`: %+v", err)
				}
			}
		}
	}

	return nil
}
