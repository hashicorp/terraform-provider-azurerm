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
				ValidateFunc: validate.ShareID,
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

	shareId, err := parse.ShareID(d.Get("share_id").(string))
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
		return fmt.Errorf("empty or nil ID returned for %s", id)
	}

	d.SetId(id.ID())
	d.Set("name", id.Name)
	d.Set("share_id", shareId.ID())

	switch resp := respModel.Value.(type) {
	case datashare.ADLSGen2FileDataSet:
		if props := resp.ADLSGen2FileProperties; props != nil {
			d.Set("storage_account_id", fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s", *props.SubscriptionID, *props.ResourceGroup, *props.StorageAccountName))
			d.Set("file_system_name", props.FileSystem)
			d.Set("file_path", props.FilePath)
			d.Set("display_name", props.DataSetID)
		}

	case datashare.ADLSGen2FolderDataSet:
		if props := resp.ADLSGen2FolderProperties; props != nil {
			d.Set("storage_account_id", fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s", *props.SubscriptionID, *props.ResourceGroup, *props.StorageAccountName))
			d.Set("file_system_name", props.FileSystem)
			d.Set("folder_path", props.FolderPath)
			d.Set("display_name", props.DataSetID)
		}

	case datashare.ADLSGen2FileSystemDataSet:
		if props := resp.ADLSGen2FileSystemProperties; props != nil {
			d.Set("storage_account_id", fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s", *props.SubscriptionID, *props.ResourceGroup, *props.StorageAccountName))
			d.Set("file_system_name", props.FileSystem)
			d.Set("display_name", props.DataSetID)
		}

	default:
		return fmt.Errorf("%s is not a datalake store gen2 dataset", id)
	}

	return nil
}
