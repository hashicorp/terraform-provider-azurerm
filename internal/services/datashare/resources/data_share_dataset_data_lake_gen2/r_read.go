// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_share_dataset_data_lake_gen2

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/dataset"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/share"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataShareDataSetDataLakeGen2Read(d *pluginsdk.ResourceData, meta any) error {
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
	d.Set("share_id", shareId.ID())

	if model := resp.Model; model != nil {
		if ds, ok := model.(dataset.ADLSGen2FileDataSet); ok {
			props := ds.Properties
			d.Set("storage_account_id", commonids.NewStorageAccountID(props.SubscriptionId, props.ResourceGroup, props.StorageAccountName).ID())
			d.Set("file_system_name", props.FileSystem)
			d.Set("file_path", props.FilePath)
			d.Set("display_name", props.DataSetId)
		} else if ds, ok := model.(dataset.ADLSGen2FolderDataSet); ok {
			props := ds.Properties
			d.Set("storage_account_id", commonids.NewStorageAccountID(props.SubscriptionId, props.ResourceGroup, props.StorageAccountName).ID())
			d.Set("file_system_name", props.FileSystem)
			d.Set("folder_path", props.FolderPath)
			d.Set("display_name", props.DataSetId)
		} else if ds, ok := model.(dataset.ADLSGen2FileSystemDataSet); ok {
			props := ds.Properties
			d.Set("storage_account_id", commonids.NewStorageAccountID(props.SubscriptionId, props.ResourceGroup, props.StorageAccountName).ID())
			d.Set("file_system_name", props.FileSystem)
			d.Set("display_name", props.DataSetId)
		} else {
			return fmt.Errorf("%s is not a datalake store gen2 dataset", *id)
		}
	}
	return nil
}
