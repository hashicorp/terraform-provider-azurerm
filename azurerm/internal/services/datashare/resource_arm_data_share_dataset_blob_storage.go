package datashare

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datashare/mgmt/2019-11-01/datashare"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	azValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/helper"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/datashare/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDataShareDataSetBlobStorage() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDataShareDataSetBlobStorageCreate,
		Read:   resourceArmDataShareDataSetBlobStorageRead,
		Delete: resourceArmDataShareDataSetBlobStorageDelete,

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

			"data_share_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataShareID,
			},

			"container_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azValidate.StorageContainerName,
			},

			"storage_account": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: storage.ValidateArmStorageAccountName,
						},

						"resource_group_name": azure.SchemaResourceGroupName(),

						"subscription_id": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.IsUUID,
						},
					},
				},
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

func resourceArmDataShareDataSetBlobStorageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.DataSetClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	shareId, err := parse.DataShareID(d.Get("data_share_id").(string))
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, shareId.ResourceGroup, shareId.AccountName, shareId.Name, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for present of existing DataShare Blob Storage DataSet %q (Resource Group %q / accountName %q / shareName %q): %+v", name, shareId.ResourceGroup, shareId.AccountName, shareId.Name, err)
		}
	}
	existingId := helper.GetAzurermDataShareDataSetId(existing.Value)
	if existingId != nil && *existingId != "" {
		return tf.ImportAsExistsError("azurerm_data_share_dataset_blob_storage", *existingId)
	}

	var dataSet datashare.BasicDataSet
	if filePath, ok := d.GetOk("file_path"); ok {
		dataSet = datashare.BlobDataSet{
			Kind: datashare.KindBlob,
			BlobProperties: &datashare.BlobProperties{
				ContainerName:      utils.String(d.Get("container_name").(string)),
				StorageAccountName: utils.String(d.Get("storage_account.0.name").(string)),
				ResourceGroup:      utils.String(d.Get("storage_account.0.resource_group_name").(string)),
				SubscriptionID:     utils.String(d.Get("storage_account.0.subscription_id").(string)),
				FilePath:           utils.String(filePath.(string)),
			},
		}
	} else if folderPath, ok := d.GetOk("folder_path"); ok {
		dataSet = datashare.BlobFolderDataSet{
			Kind: datashare.KindBlobFolder,
			BlobFolderProperties: &datashare.BlobFolderProperties{
				ContainerName:      utils.String(d.Get("container_name").(string)),
				StorageAccountName: utils.String(d.Get("storage_account.0.name").(string)),
				ResourceGroup:      utils.String(d.Get("storage_account.0.resource_group_name").(string)),
				SubscriptionID:     utils.String(d.Get("storage_account.0.subscription_id").(string)),
				Prefix:             utils.String(folderPath.(string)),
			},
		}
	} else {
		dataSet = datashare.BlobContainerDataSet{
			Kind: datashare.KindContainer,
			BlobContainerProperties: &datashare.BlobContainerProperties{
				ContainerName:      utils.String(d.Get("container_name").(string)),
				StorageAccountName: utils.String(d.Get("storage_account.0.name").(string)),
				ResourceGroup:      utils.String(d.Get("storage_account.0.resource_group_name").(string)),
				SubscriptionID:     utils.String(d.Get("storage_account.0.subscription_id").(string)),
			},
		}
	}

	if _, err := client.Create(ctx, shareId.ResourceGroup, shareId.AccountName, shareId.Name, name, dataSet); err != nil {
		return fmt.Errorf("creating DataShare Blob Storage DataSet %q (Resource Group %q / accountName %q / shareName %q): %+v", name, shareId.ResourceGroup, shareId.AccountName, shareId.Name, err)
	}

	resp, err := client.Get(ctx, shareId.ResourceGroup, shareId.AccountName, shareId.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving DataShare Blob Storage DataSet %q (Resource Group %q / accountName %q / shareName %q): %+v", name, shareId.ResourceGroup, shareId.AccountName, shareId.Name, err)
	}

	respId := helper.GetAzurermDataShareDataSetId(resp.Value)
	if respId == nil || *respId == "" {
		return fmt.Errorf("empty or nil ID returned for DataShare Blob Storage DataSet %q (Resource Group %q / accountName %q / shareName %q)", name, shareId.ResourceGroup, shareId.AccountName, shareId.Name)
	}

	d.SetId(*respId)
	return resourceArmDataShareDataSetBlobStorageRead(d, meta)
}

func resourceArmDataShareDataSetBlobStorageRead(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("retrieving DataShare Blob Storage DataSet %q (Resource Group %q / accountName %q / shareName %q): %+v", id.Name, id.ResourceGroup, id.AccountName, id.ShareName, err)
	}

	d.Set("name", id.Name)
	shareResp, err := shareClient.Get(ctx, id.ResourceGroup, id.AccountName, id.ShareName)
	if err != nil {
		return fmt.Errorf("retrieving DataShare %q (Resource Group %q / accountName %q): %+v", id.ShareName, id.ResourceGroup, id.AccountName, err)
	}
	if shareResp.ID == nil || *shareResp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for DataShare %q (Resource Group %q / accountName %q)", id.ShareName, id.ResourceGroup, id.AccountName)
	}

	d.Set("data_share_id", shareResp.ID)

	switch resp := resp.Value.(type) {
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
		return fmt.Errorf("data share dataset %q (Resource Group %q / accountName %q / shareName %q) is not a blob storage dataset", id.Name, id.ResourceGroup, id.AccountName, id.ShareName)
	}

	return nil
}

func resourceArmDataShareDataSetBlobStorageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.DataSetClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataShareDataSetID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.AccountName, id.ShareName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting DataShare Blob Storage DataSet %q (Resource Group %q / accountName %q / shareName %q): %+v", id.Name, id.ResourceGroup, id.AccountName, id.ShareName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of DataShare Blob Storage DataSet %q (Resource Group %q / accountName %q / shareName %q): %+v", id.Name, id.ResourceGroup, id.AccountName, id.ShareName, err)
	}
	return nil
}

func flattenAzureRmDataShareDataSetBlobStorageAccount(strName, strRG, strSubs *string) []interface{} {
	var name, rg, subs string
	if strName != nil {
		name = *strName
	}

	if strRG != nil {
		rg = *strRG
	}

	if strSubs != nil {
		subs = *strSubs
	}

	return []interface{}{
		map[string]interface{}{
			"name":                name,
			"resource_group_name": rg,
			"subscription_id":     subs,
		},
	}
}
