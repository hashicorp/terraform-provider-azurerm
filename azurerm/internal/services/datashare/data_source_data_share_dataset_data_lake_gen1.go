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

func dataSourceDataShareDatasetDataLakeGen1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmDataShareDatasetDataLakeGen1Read,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.DataSetName(),
			},

			"data_share_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ShareID,
			},

			"data_lake_store_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"folder_path": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"file_name": {
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

func dataSourceArmDataShareDatasetDataLakeGen1Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.DataSetClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	shareID := d.Get("data_share_id").(string)
	shareId, err := parse.ShareID(shareID)
	if err != nil {
		return err
	}

	respModel, err := client.Get(ctx, shareId.ResourceGroup, shareId.AccountName, shareId.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving DataShare Data Lake Gen1 DataSet %q (Resource Group %q / accountName %q / shareName %q): %+v", name, shareId.ResourceGroup, shareId.AccountName, shareId.Name, err)
	}

	respId := helper.GetAzurermDataShareDataSetId(respModel.Value)
	if respId == nil || *respId == "" {
		return fmt.Errorf("empty or nil ID returned for DataShare Data Lake Gen1 DataSet %q (Resource Group %q / accountName %q / shareName %q)", name, shareId.ResourceGroup, shareId.AccountName, shareId.Name)
	}

	d.SetId(*respId)
	d.Set("name", name)
	d.Set("data_share_id", shareID)

	switch resp := respModel.Value.(type) {
	case datashare.ADLSGen1FileDataSet:
		if props := resp.ADLSGen1FileProperties; props != nil {
			if props.SubscriptionID != nil && props.ResourceGroup != nil && props.AccountName != nil {
				d.Set("data_lake_store_id", fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataLakeStore/accounts/%s", *props.SubscriptionID, *props.ResourceGroup, *props.AccountName))
			}
			d.Set("folder_path", props.FolderPath)
			d.Set("file_name", props.FileName)
			d.Set("display_name", props.DataSetID)
		}

	case datashare.ADLSGen1FolderDataSet:
		if props := resp.ADLSGen1FolderProperties; props != nil {
			if props.SubscriptionID != nil && props.ResourceGroup != nil && props.AccountName != nil {
				d.Set("data_lake_store_id", fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataLakeStore/accounts/%s", *props.SubscriptionID, *props.ResourceGroup, *props.AccountName))
			}
			d.Set("folder_path", props.FolderPath)
			d.Set("display_name", props.DataSetID)
		}

	default:
		return fmt.Errorf("data share dataset %q (Resource Group %q / accountName %q / shareName %q) is not a datalake store gen1 dataset", name, shareId.ResourceGroup, shareId.AccountName, shareId.Name)
	}

	return nil
}
