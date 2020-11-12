package datashare

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datashare/mgmt/2019-11-01/datashare"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/helper"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func dataSourceDataShareDatasetDataLakeGen2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmDataShareDatasetDataLakeGen2Read,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.DatashareDataSetName(),
			},

			"share_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.DataShareID,
			},

			"storage_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"file_system_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"file_path": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"folder_path": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmDataShareDatasetDataLakeGen2Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.DataSetClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	shareID := d.Get("share_id").(string)
	shareId, err := parse.DataShareID(shareID)
	if err != nil {
		return err
	}

	respModel, err := client.Get(ctx, shareId.ResourceGroup, shareId.AccountName, shareId.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving DataShare Data Lake Gen2 DataSet %q (Resource Group %q / accountName %q / shareName %q): %+v", name, shareId.ResourceGroup, shareId.AccountName, shareId.Name, err)
	}

	respId := helper.GetAzurermDataShareDataSetId(respModel.Value)
	if respId == nil || *respId == "" {
		return fmt.Errorf("empty or nil ID returned for DataShare Data Lake Gen2 DataSet %q (Resource Group %q / accountName %q / shareName %q)", name, shareId.ResourceGroup, shareId.AccountName, shareId.Name)
	}

	d.SetId(*respId)
	d.Set("name", name)
	d.Set("share_id", shareID)

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
		return fmt.Errorf("data share dataset %q (Resource Group %q / accountName %q / shareName %q) is not a datalake store gen2 dataset", name, shareId.ResourceGroup, shareId.AccountName, shareId.Name)
	}

	return nil
}
