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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/helper"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/validate"
	storageParsers "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parsers"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDataShareDataSetDataLakeGen2() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDataShareDataSetDataLakeGen2Create,
		Read:   resourceArmDataShareDataSetDataLakeGen2Read,
		Delete: resourceArmDataShareDataSetDataLakeGen2Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.DataShareDataSetID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DatashareDataSetName(),
			},

			"share_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataShareID,
			},

			"storage_account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.StorageAccountID,
			},

			"file_system_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"file_path": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringIsNotEmpty,
				ConflictsWith: []string{"folder_path"},
			},

			"folder_path": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringIsNotEmpty,
				ConflictsWith: []string{"file_path"},
			},

			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceArmDataShareDataSetDataLakeGen2Create(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.DataSetClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	shareId, err := parse.DataShareID(d.Get("share_id").(string))
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
		return tf.ImportAsExistsError("azurerm_data_share_dataset_data_lake_gen2", *existingId)
	}

	strId, err := storageParsers.ParseAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	var dataSet datashare.BasicDataSet

	if filePath, ok := d.GetOk("file_path"); ok {
		dataSet = datashare.ADLSGen2FileDataSet{
			Kind: datashare.KindAdlsGen2File,
			ADLSGen2FileProperties: &datashare.ADLSGen2FileProperties{
				StorageAccountName: utils.String(strId.Name),
				ResourceGroup:      utils.String(strId.ResourceGroup),
				SubscriptionID:     utils.String(strId.SubscriptionId),
				FileSystem:         utils.String(d.Get("file_system_name").(string)),
				FilePath:           utils.String(filePath.(string)),
			},
		}
	} else if folderPath, ok := d.GetOk("folder_path"); ok {
		dataSet = datashare.ADLSGen2FolderDataSet{
			Kind: datashare.KindAdlsGen2Folder,
			ADLSGen2FolderProperties: &datashare.ADLSGen2FolderProperties{
				StorageAccountName: utils.String(strId.Name),
				ResourceGroup:      utils.String(strId.ResourceGroup),
				SubscriptionID:     utils.String(strId.SubscriptionId),
				FileSystem:         utils.String(d.Get("file_system_name").(string)),
				FolderPath:         utils.String(folderPath.(string)),
			},
		}
	} else {
		dataSet = datashare.ADLSGen2FileSystemDataSet{
			Kind: datashare.KindAdlsGen2FileSystem,
			ADLSGen2FileSystemProperties: &datashare.ADLSGen2FileSystemProperties{
				StorageAccountName: utils.String(strId.Name),
				ResourceGroup:      utils.String(strId.ResourceGroup),
				SubscriptionID:     utils.String(strId.SubscriptionId),
				FileSystem:         utils.String(d.Get("file_system_name").(string)),
			},
		}
	}

	if _, err := client.Create(ctx, shareId.ResourceGroup, shareId.AccountName, shareId.Name, name, dataSet); err != nil {
		return fmt.Errorf("creating DataShare DataSet %q (Resource Group %q / accountName %q / shareName %q): %+v", name, shareId.ResourceGroup, shareId.AccountName, shareId.Name, err)
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
	return resourceArmDataShareDataSetDataLakeGen2Read(d, meta)
}

func resourceArmDataShareDataSetDataLakeGen2Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.DataSetClient
	shareClient := meta.(*clients.Client).DataShare.SharesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataShareDataSetID(d.Id())
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
		return fmt.Errorf("empty or nil ID returned for DataShare %q (Resource Group %q / accountName %q)", id.ShareName, id.ResourceGroup, id.AccountName)
	}
	d.Set("share_id", shareResp.ID)

	switch resp := resp.Value.(type) {
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
		return fmt.Errorf("data share dataset %q (Resource Group %q / accountName %q / shareName %q) is not a datalake store gen2 dataset", id.Name, id.ResourceGroup, id.AccountName, id.ShareName)
	}

	return nil
}

func resourceArmDataShareDataSetDataLakeGen2Delete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.DataSetClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataShareDataSetID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.AccountName, id.ShareName, id.Name); err != nil {
		return fmt.Errorf("deleting DataShare DataSet %q (Resource Group %q / accountName %q / shareName %q): %+v", id.Name, id.ResourceGroup, id.AccountName, id.ShareName, err)
	}
	return nil
}
