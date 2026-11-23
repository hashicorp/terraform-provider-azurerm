// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_output_blob

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStreamAnalyticsOutputBlobRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).StreamAnalytics.OutputsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := outputs.ParseOutputID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.OutputName)
	d.Set("stream_analytics_job_name", id.StreamingJobName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			output, ok := props.Datasource.(outputs.BlobOutputDataSource)
			if !ok {
				return fmt.Errorf("converting %s to a Blob Output", *id)
			}

			d.Set("date_format", pointer.From(output.Properties.DateFormat))

			d.Set("path_pattern", pointer.From(output.Properties.PathPattern))

			d.Set("storage_container_name", pointer.From(output.Properties.Container))

			d.Set("time_format", pointer.From(output.Properties.TimeFormat))

			authenticationMode := ""
			if v := output.Properties.AuthenticationMode; v != nil {
				authenticationMode = string(*v)
			}
			d.Set("authentication_mode", authenticationMode)

			if accounts := output.Properties.StorageAccounts; accounts != nil && len(*accounts) > 0 {
				account := (*accounts)[0]
				d.Set("storage_account_name", account.AccountName)
			}

			if err := d.Set("serialization", helpers.FlattenStreamAnalyticsOutputSerialization(props.Serialization)); err != nil {
				return fmt.Errorf("setting `serialization`: %+v", err)
			}
			d.Set("batch_max_wait_time", props.TimeWindow)
			d.Set("batch_min_rows", props.SizeWindow)

			blobWriteMode := ""
			if v := output.Properties.BlobWriteMode; v != nil {
				blobWriteMode = string(*v)
			}
			d.Set("blob_write_mode", blobWriteMode)
		}
	}
	return nil
}

func getStorageAccountKey(input string) *string {
	if input == "" {
		return nil
	}

	return pointer.To(input)
}
