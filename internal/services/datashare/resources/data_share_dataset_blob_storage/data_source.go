// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_share_dataset_blob_storage

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/dataset"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/share"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Data Share Blob Storage Dataset.",
		Read:        dataSourceDataShareDatasetBlobStorageRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of this Data Share Blob Storage Dataset.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.DataSetName(),
			},

			"data_share_id": {
				Description:  "The ID of the Data Share in which this Data Share Blob Storage Dataset should be created.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: share.ValidateShareID,
			},

			"container_name": {
				Description: "The name of the Storage Account container to be shared with the receiver.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"storage_account": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Description: "The name of the Storage Account to be shared with the receiver.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"resource_group_name": {
							Description: "The Resource Group name of the Storage Account to be shared with the receiver.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},

						"subscription_id": {
							Description: "The subscription ID of the Storage Account to be shared with the receiver.",
							Type:        pluginsdk.TypeString,
							Computed:    true,
						},
					},
				},
			},

			"file_path": {
				Description: "The path of the file in the Storage Container to be shared with the receiver.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"folder_path": {
				Description: "The folder path of the file in the Storage Container to be shared with the receiver.",
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

func dataSourceDataShareDatasetBlobStorageRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataShare.DataSetClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	shareId, err := share.ParseShareID(d.Get("data_share_id").(string))
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
			return fmt.Errorf("%s is not a blob storage dataset", id)
		}
	}

	return nil
}
