// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_share_dataset_data_lake_gen2

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/dataset"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/share"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	datashareClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Data Share Data Lake Gen2 Dataset.",
		Read:        dataSourceDataShareDatasetDataLakeGen2Read,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of this Data Share Data Lake Gen2 Dataset.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.DataSetName(),
			},

			"share_id": {
				Description:  "The resource ID of the Data Share where this Data Share Data Lake Gen2 Dataset should be created.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: share.ValidateShareID,
			},

			"storage_account_id": {
				Description: "The resource ID of the Storage Account of the Data Lake file system to be shared with the receiver.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"file_system_name": {
				Description: "The name of the Data Lake file system to be shared with the receiver.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"file_path": {
				Description: "The path of the file in the Data Lake file system to be shared with the receiver.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"folder_path": {
				Description: "The folder path in the Data Lake file system to be shared with the receiver.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"display_name": {
				Description: "The name of the Data Share Dataset.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceDataShareDatasetDataLakeGen2Read(d *pluginsdk.ResourceData, meta any) error {
	client := datashareClient.From(meta).DataSetClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	shareId, err := share.ParseShareID(d.Get("share_id").(string))
	if err != nil {
		return err
	}
	id := dataset.NewDataSetID(shareId.SubscriptionId, shareId.ResourceGroupName, shareId.AccountName, shareId.ShareName, d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.DataSetName)
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
			return fmt.Errorf("%s is not a datalake store gen2 dataset", id)
		}
	}

	return nil
}
