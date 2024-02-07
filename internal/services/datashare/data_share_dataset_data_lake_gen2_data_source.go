// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datashare

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/dataset"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/share"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceDataShareDatasetDataLakeGen2() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceDataShareDatasetDataLakeGen2Read,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.DataSetName(),
			},

			"share_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: share.ValidateShareID,
			},

			"storage_account_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"file_system_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
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

func dataSourceDataShareDatasetDataLakeGen2Read(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.DataSetClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
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
		m := *model
		if ds, ok := m.(dataset.ADLSGen2FileDataSet); ok {
			props := ds.Properties
			d.Set("storage_account_id", commonids.NewStorageAccountID(props.SubscriptionId, props.ResourceGroup, props.StorageAccountName).ID())
			d.Set("file_system_name", props.FileSystem)
			d.Set("file_path", props.FilePath)
			d.Set("display_name", props.DataSetId)
		} else if ds, ok := m.(dataset.ADLSGen2FolderDataSet); ok {
			props := ds.Properties
			d.Set("storage_account_id", commonids.NewStorageAccountID(props.SubscriptionId, props.ResourceGroup, props.StorageAccountName).ID())
			d.Set("file_system_name", props.FileSystem)
			d.Set("folder_path", props.FolderPath)
			d.Set("display_name", props.DataSetId)
		} else if ds, ok := m.(dataset.ADLSGen2FileSystemDataSet); ok {
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
