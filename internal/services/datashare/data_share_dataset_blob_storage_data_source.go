// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datashare

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

func dataSourceDataShareDatasetBlobStorage() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceDataShareDatasetBlobStorageRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.DataSetName(),
			},

			"data_share_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: share.ValidateShareID,
			},

			"container_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"storage_account": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"resource_group_name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"subscription_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"file_path": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"folder_path": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"display_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDataShareDatasetBlobStorageRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		m := *model
		if ds, ok := m.(dataset.BlobDataSet); ok {
			props := ds.Properties
			d.Set("container_name", props.ContainerName)
			if err := d.Set("storage_account", flattenAzureRmDataShareDataSetBlobStorageAccount(props.StorageAccountName, props.ResourceGroup, props.SubscriptionId)); err != nil {
				return fmt.Errorf("setting `storage_account`: %+v", err)
			}
			d.Set("file_path", props.FilePath)
			d.Set("display_name", props.DataSetId)

		} else if ds, ok := m.(dataset.BlobFolderDataSet); ok {
			props := ds.Properties
			d.Set("container_name", props.ContainerName)
			if err := d.Set("storage_account", flattenAzureRmDataShareDataSetBlobStorageAccount(props.StorageAccountName, props.ResourceGroup, props.SubscriptionId)); err != nil {
				return fmt.Errorf("setting `storage_account`: %+v", err)
			}
			d.Set("folder_path", props.Prefix)
			d.Set("display_name", props.DataSetId)
		} else if ds, ok := m.(dataset.BlobContainerDataSet); ok {
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
