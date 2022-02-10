package datashare

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datashare/mgmt/2019-11-01/datashare"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	dataLakeParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/datalake/parse"
	dataLakeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/datalake/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/helper"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDataShareDataSetDataLakeGen1() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataShareDataSetDataLakeGen1Create,
		Read:   resourceDataShareDataSetDataLakeGen1Read,
		Delete: resourceDataShareDataSetDataLakeGen1Delete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DataSetID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataSetName(),
			},

			"data_share_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ShareID,
			},

			"data_lake_store_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: dataLakeValidate.AccountID,
			},

			"folder_path": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"file_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"display_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDataShareDataSetDataLakeGen1Create(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.DataSetClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	shareId, err := parse.ShareID(d.Get("data_share_id").(string))
	if err != nil {
		return err
	}
	id := parse.NewDataSetID(subscriptionId, shareId.ResourceGroup, shareId.AccountName, shareId.Name, d.Get("name").(string))

	existing, err := client.Get(ctx, id.ResourceGroup, id.AccountName, id.ShareName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for present of existing %s: %+v", id, err)
		}
	}
	existingId := helper.GetAzurermDataShareDataSetId(existing.Value)
	if existingId != nil && *existingId != "" {
		return tf.ImportAsExistsError("azurerm_data_share_dataset_data_lake_gen1", *existingId)
	}

	dataLakeStoreId, err := dataLakeParse.AccountID(d.Get("data_lake_store_id").(string))
	if err != nil {
		return err
	}

	var dataSet datashare.BasicDataSet

	if fileName, ok := d.GetOk("file_name"); ok {
		dataSet = datashare.ADLSGen1FileDataSet{
			Kind: datashare.KindAdlsGen1File,
			ADLSGen1FileProperties: &datashare.ADLSGen1FileProperties{
				AccountName:    utils.String(dataLakeStoreId.Name),
				ResourceGroup:  utils.String(dataLakeStoreId.ResourceGroup),
				SubscriptionID: utils.String(dataLakeStoreId.SubscriptionId),
				FolderPath:     utils.String(d.Get("folder_path").(string)),
				FileName:       utils.String(fileName.(string)),
			},
		}
	} else {
		dataSet = datashare.ADLSGen1FolderDataSet{
			Kind: datashare.KindAdlsGen1Folder,
			ADLSGen1FolderProperties: &datashare.ADLSGen1FolderProperties{
				AccountName:    utils.String(dataLakeStoreId.Name),
				ResourceGroup:  utils.String(dataLakeStoreId.ResourceGroup),
				SubscriptionID: utils.String(dataLakeStoreId.SubscriptionId),
				FolderPath:     utils.String(d.Get("folder_path").(string)),
			},
		}
	}

	if _, err := client.Create(ctx, id.ResourceGroup, id.AccountName, id.ShareName, id.Name, dataSet); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDataShareDataSetDataLakeGen1Read(d, meta)
}

func resourceDataShareDataSetDataLakeGen1Read(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.DataSetClient
	shareClient := meta.(*clients.Client).DataShare.SharesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataSetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.AccountName, id.ShareName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] DataShare %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving DataShare DataSet %q (Resource Group %q / accountName %q / shareName %q): %+v", id.Name, id.ResourceGroup, id.AccountName, id.ShareName, err)
	}
	d.Set("name", id.Name)
	shareResp, err := shareClient.Get(ctx, id.ResourceGroup, id.AccountName, id.ShareName)
	if err != nil {
		return fmt.Errorf("retrieving DataShare %q (Resource Group %q / accountName %q): %+v", id.ShareName, id.ResourceGroup, id.AccountName, err)
	}
	if shareResp.ID == nil || *shareResp.ID == "" {
		return fmt.Errorf("reading ID of DataShare %q (Resource Group %q / accountName %q): ID is empt", id.ShareName, id.ResourceGroup, id.AccountName)
	}
	d.Set("data_share_id", shareResp.ID)

	switch resp := resp.Value.(type) {
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
		return fmt.Errorf("data share dataset %q (Resource Group %q / accountName %q / shareName %q) is not a datalake store gen1 dataset", id.Name, id.ResourceGroup, id.AccountName, id.ShareName)
	}

	return nil
}

func resourceDataShareDataSetDataLakeGen1Delete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.DataSetClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataSetID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.AccountName, id.ShareName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting DataShare DataSet %q (Resource Group %q / accountName %q / shareName %q): %+v", id.Name, id.ResourceGroup, id.AccountName, id.ShareName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of DataShare Data Lake Gen1 DataSet  %q (Resource Group %q / accountName %q / shareName %q): %+v", id.Name, id.ResourceGroup, id.AccountName, id.ShareName, err)
	}

	return nil
}
