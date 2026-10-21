// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_share_dataset_data_lake_gen2

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/dataset"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/share"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataShareDataSetDataLakeGen2Create(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataShare.DataSetClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	shareId, err := share.ParseShareID(d.Get("share_id").(string))
	if err != nil {
		return err
	}
	id := dataset.NewDataSetID(shareId.SubscriptionId, shareId.ResourceGroupName, shareId.AccountName, shareId.ShareName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_data_share_dataset_data_lake_gen2", id.ID())
		}
	}

	strId, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	var dataSet dataset.DataSet

	if filePath, ok := d.GetOk("file_path"); ok {
		dataSet = dataset.ADLSGen2FileDataSet{
			Properties: dataset.ADLSGen2FileProperties{
				StorageAccountName: strId.StorageAccountName,
				ResourceGroup:      strId.ResourceGroupName,
				SubscriptionId:     strId.SubscriptionId,
				FileSystem:         d.Get("file_system_name").(string),
				FilePath:           filePath.(string),
			},
		}
	} else if folderPath, ok := d.GetOk("folder_path"); ok {
		dataSet = dataset.ADLSGen2FolderDataSet{
			Properties: dataset.ADLSGen2FolderProperties{
				StorageAccountName: strId.StorageAccountName,
				ResourceGroup:      strId.ResourceGroupName,
				SubscriptionId:     strId.SubscriptionId,
				FileSystem:         d.Get("file_system_name").(string),
				FolderPath:         folderPath.(string),
			},
		}
	} else {
		dataSet = dataset.ADLSGen2FileSystemDataSet{
			Properties: dataset.ADLSGen2FileSystemProperties{
				StorageAccountName: strId.StorageAccountName,
				ResourceGroup:      strId.ResourceGroupName,
				SubscriptionId:     strId.SubscriptionId,
				FileSystem:         d.Get("file_system_name").(string),
			},
		}
	}

	if _, err := client.Create(ctx, id, dataSet); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDataShareDataSetDataLakeGen2Read(d, meta)
}
