// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_reference_input_blob

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/inputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStreamAnalyticsReferenceInputBlobRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).StreamAnalytics.InputsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := inputs.ParseInputID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
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
			dataSource, ok := props.(inputs.ReferenceInputProperties)
			if !ok {
				return fmt.Errorf("converting %s to a Reference Input", *id)
			}

			referenceInputBlob, ok := dataSource.Datasource.(inputs.BlobReferenceInputDataSource)
			if !ok {
				return fmt.Errorf("converting %s to a Blob Reference Input", *id)
			}

			if referenceInputBlob.Properties != nil {
				d.Set("date_format", pointer.From(referenceInputBlob.Properties.DateFormat))

				d.Set("path_pattern", pointer.From(referenceInputBlob.Properties.PathPattern))

				d.Set("storage_container_name", pointer.From(referenceInputBlob.Properties.Container))

				d.Set("time_format", pointer.From(referenceInputBlob.Properties.TimeFormat))

				authMode := ""
				if v := referenceInputBlob.Properties.AuthenticationMode; v != nil {
					authMode = string(*v)
				}
				d.Set("authentication_mode", authMode)

				if accounts := referenceInputBlob.Properties.StorageAccounts; accounts != nil && len(*accounts) > 0 {
					account := (*accounts)[0]
					d.Set("storage_account_name", account.AccountName)
				}
			}
			if err := d.Set("serialization", helpers.FlattenStreamAnalyticsStreamInputSerialization(dataSource.Serialization)); err != nil {
				return fmt.Errorf("setting `serialization`: %+v", err)
			}
		}
	}

	return nil
}

func normalizeAccountKey(accountKey string) *string {
	if accountKey != "" {
		return pointer.To(accountKey)
	}

	return nil
}
