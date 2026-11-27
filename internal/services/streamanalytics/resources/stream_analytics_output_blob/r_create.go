// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_output_blob

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	streamanalyticsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceStreamAnalyticsOutputBlobCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := streamanalyticsClient.From(meta).OutputsClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	id := outputs.NewOutputID(subscriptionId, d.Get("resource_group_name").(string), d.Get("stream_analytics_job_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_stream_analytics_output_blob", id.ID())
			}
		}
	}

	containerName := d.Get("storage_container_name").(string)
	dateFormat := d.Get("date_format").(string)
	pathPattern := d.Get("path_pattern").(string)
	storageAccountName := d.Get("storage_account_name").(string)
	timeFormat := d.Get("time_format").(string)

	serializationRaw := d.Get("serialization").([]any)
	serialization, err := helpers.ExpandStreamAnalyticsOutputSerialization(serializationRaw)
	if err != nil {
		return fmt.Errorf("expanding `serialization`: %+v", err)
	}

	props := outputs.Output{
		Name: pointer.To(id.OutputName),
		Properties: &outputs.OutputProperties{
			Datasource: &outputs.BlobOutputDataSource{
				Properties: &outputs.BlobOutputDataSourceProperties{
					StorageAccounts: &[]outputs.StorageAccount{
						{
							AccountKey:  getStorageAccountKey(d.Get("storage_account_key").(string)),
							AccountName: pointer.To(storageAccountName),
						},
					},
					Container:          pointer.To(containerName),
					DateFormat:         pointer.To(dateFormat),
					PathPattern:        pointer.To(pathPattern),
					TimeFormat:         pointer.To(timeFormat),
					AuthenticationMode: pointer.ToEnum[outputs.AuthenticationMode](d.Get("authentication_mode").(string)),
					BlobWriteMode:      pointer.ToEnum[outputs.BlobWriteMode](d.Get("blob_write_mode").(string)),
				},
			},
			Serialization: serialization,
		},
	}

	if batchMaxWaitTime, ok := d.GetOk("batch_max_wait_time"); ok {
		props.Properties.TimeWindow = pointer.To(batchMaxWaitTime.(string))
	}

	if batchMinRows, ok := d.GetOk("batch_min_rows"); ok {
		props.Properties.SizeWindow = pointer.To(int64(batchMinRows.(int)))
	}

	// timeWindow and sizeWindow must be set for Parquet serialization
	_, isParquet := serialization.(outputs.ParquetSerialization)
	if isParquet && (props.Properties.TimeWindow == nil || props.Properties.SizeWindow == nil) {
		return fmt.Errorf("cannot create %s: batch_min_rows and batch_time_window must be set for Parquet serialization", id)
	}

	var createOpts outputs.CreateOrReplaceOperationOptions
	var updateOpts outputs.UpdateOperationOptions
	if d.IsNewResource() {
		if _, err := client.CreateOrReplace(ctx, id, props, createOpts); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}

		d.SetId(id.ID())
	} else if _, err := client.Update(ctx, id, props, updateOpts); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceStreamAnalyticsOutputBlobRead(d, meta)
}
