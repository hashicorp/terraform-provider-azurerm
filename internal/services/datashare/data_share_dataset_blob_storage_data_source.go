package datashare

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datashare/mgmt/2019-11-01/datashare"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/parse"
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
				ValidateFunc: validate.ShareID,
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

	shareId, err := parse.ShareID(d.Get("data_share_id").(string))
	if err != nil {
		return err
	}
	id := parse.NewDataSetID(shareId.SubscriptionId, shareId.ResourceGroup, shareId.AccountName, shareId.Name, d.Get("name").(string))

	respModel, err := client.Get(ctx, id.ResourceGroup, id.AccountName, id.ShareName, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	respId := helper.GetAzurermDataShareDataSetId(respModel.Value)
	if respId == nil || *respId == "" {
		return fmt.Errorf("empty or nil ID returned for reading %s", id)
	}

	d.SetId(id.ID())
	d.Set("name", id.Name)
	d.Set("data_share_id", shareId.ID())

	switch resp := respModel.Value.(type) {
	case datashare.BlobDataSet:
		if props := resp.BlobProperties; props != nil {
			d.Set("container_name", props.ContainerName)
			if err := d.Set("storage_account", flattenAzureRmDataShareDataSetBlobStorageAccount(props.StorageAccountName, props.ResourceGroup, props.SubscriptionID)); err != nil {
				return fmt.Errorf("setting `storage_account`: %+v", err)
			}
			d.Set("file_path", props.FilePath)
			d.Set("display_name", props.DataSetID)
		}

	case datashare.BlobFolderDataSet:
		if props := resp.BlobFolderProperties; props != nil {
			d.Set("container_name", props.ContainerName)
			if err := d.Set("storage_account", flattenAzureRmDataShareDataSetBlobStorageAccount(props.StorageAccountName, props.ResourceGroup, props.SubscriptionID)); err != nil {
				return fmt.Errorf("setting `storage_account`: %+v", err)
			}
			d.Set("folder_path", props.Prefix)
			d.Set("display_name", props.DataSetID)
		}

	case datashare.BlobContainerDataSet:
		if props := resp.BlobContainerProperties; props != nil {
			d.Set("container_name", props.ContainerName)
			if err := d.Set("storage_account", flattenAzureRmDataShareDataSetBlobStorageAccount(props.StorageAccountName, props.ResourceGroup, props.SubscriptionID)); err != nil {
				return fmt.Errorf("setting `storage_account`: %+v", err)
			}
			d.Set("display_name", props.DataSetID)
		}

	default:
		return fmt.Errorf("%s is not a blob storage dataset", id)
	}

	return nil
}
