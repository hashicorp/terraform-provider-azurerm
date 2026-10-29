// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_share_dataset_blob_storage

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/dataset"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/share"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataShareDataSetBlobStorageRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataShare.DataSetClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dataset.ParseDataSetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] DataShare %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.DataSetName)

	shareId := share.NewShareID(id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.ShareName)
	d.Set("data_share_id", shareId.ID())

	if model := resp.Model; model != nil {
		if ds, ok := model.(dataset.BlobDataSet); ok {
			props := ds.Properties
			d.Set("container_name", props.ContainerName)
			if err := d.Set("storage_account", flattenAzureRmDataShareDataSetBlobStorageAccount(props.StorageAccountName, props.ResourceGroup, props.SubscriptionId)); err != nil {
				return fmt.Errorf("setting `storage_account`: %+v", err)
			}
			d.Set("file_path", props.FilePath)
			d.Set("display_name", props.DataSetId)
		} else if ds, ok := model.(dataset.BlobFolderDataSet); ok {
			props := ds.Properties
			d.Set("container_name", props.ContainerName)
			if err := d.Set("storage_account", flattenAzureRmDataShareDataSetBlobStorageAccount(props.StorageAccountName, props.ResourceGroup, props.SubscriptionId)); err != nil {
				return fmt.Errorf("setting `storage_account`: %+v", err)
			}
			d.Set("folder_path", props.Prefix)
			d.Set("display_name", props.DataSetId)
		} else if ds, ok := model.(dataset.BlobContainerDataSet); ok {
			props := ds.Properties
			d.Set("container_name", props.ContainerName)
			if err := d.Set("storage_account", flattenAzureRmDataShareDataSetBlobStorageAccount(props.StorageAccountName, props.ResourceGroup, props.SubscriptionId)); err != nil {
				return fmt.Errorf("setting `storage_account`: %+v", err)
			}
			d.Set("display_name", props.DataSetId)
		} else {
			return fmt.Errorf("%s is not a blob storage dataset", *id)
		}
	}

	return nil
}

func flattenAzureRmDataShareDataSetBlobStorageAccount(name, rg, subs string) []any {
	return []any{
		map[string]any{
			"name":                name,
			"resource_group_name": rg,
			"subscription_id":     subs,
		},
	}
}
