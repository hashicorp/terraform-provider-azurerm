package datashare

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datashare/mgmt/2019-11-01/datashare"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	dataLakeParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datalake/parse"
	dataLakeValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datalake/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/helper"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDataShareDataSetDataLakeGen1() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataShareDataSetDataLakeGen1Create,
		Read:   resourceDataShareDataSetDataLakeGen1Read,
		Delete: resourceDataShareDataSetDataLakeGen1Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.DataSetID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataSetName(),
			},

			"data_share_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ShareID,
			},

			"data_lake_store_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: dataLakeValidate.AccountID,
			},

			"folder_path": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"file_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDataShareDataSetDataLakeGen1Create(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.DataSetClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	shareId, err := parse.ShareID(d.Get("data_share_id").(string))
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, shareId.ResourceGroup, shareId.AccountName, shareId.Name, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for present of existing  DataShare DataSet %q (Resource Group %q / accountName %q / shareName %q): %+v", name, shareId.ResourceGroup, shareId.AccountName, shareId.Name, err)
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

	if _, err := client.Create(ctx, shareId.ResourceGroup, shareId.AccountName, shareId.Name, name, dataSet); err != nil {
		return fmt.Errorf("creating/updating DataShare DataSet %q (Resource Group %q / accountName %q / shareName %q): %+v", name, shareId.ResourceGroup, shareId.AccountName, shareId.Name, err)
	}

	resp, err := client.Get(ctx, shareId.ResourceGroup, shareId.AccountName, shareId.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving DataShare DataSet %q (Resource Group %q / accountName %q / shareName %q): %+v", name, shareId.ResourceGroup, shareId.AccountName, shareId.Name, err)
	}

	respId := helper.GetAzurermDataShareDataSetId(resp.Value)
	if respId == nil || *respId == "" {
		return fmt.Errorf("empty or nil ID returned for DataShare DataSet %q (Resource Group %q / accountName %q / shareName %q)", name, shareId.ResourceGroup, shareId.AccountName, shareId.Name)
	}

	d.SetId(*respId)
	return resourceDataShareDataSetDataLakeGen1Read(d, meta)
}

func resourceDataShareDataSetDataLakeGen1Read(d *schema.ResourceData, meta interface{}) error {
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

func resourceDataShareDataSetDataLakeGen1Delete(d *schema.ResourceData, meta interface{}) error {
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
