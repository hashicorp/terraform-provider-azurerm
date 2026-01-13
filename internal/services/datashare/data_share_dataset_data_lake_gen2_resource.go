// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package datashare

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/dataset"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/share"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name data_share_dataset_data_lake_gen2 -service-package-name datashare -properties "name" -compare-values "resource_group_name:share_id,account_name:share_id,share_name:share_id" -known-values "subscription_id:data.Subscriptions.Primary" -test-name basicFile

func resourceDataShareDataSetDataLakeGen2() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataShareDataSetDataLakeGen2Create,
		Read:   resourceDataShareDataSetDataLakeGen2Read,
		Delete: resourceDataShareDataSetDataLakeGen2Delete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingIdentity(&dataset.DataSetId{}),
		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&dataset.DataSetId{}),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataSetName(),
			},

			"share_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: share.ValidateShareID,
			},

			"storage_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateStorageAccountID,
			},

			"file_system_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"file_path": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringIsNotEmpty,
				ConflictsWith: []string{"folder_path"},
			},

			"folder_path": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validation.StringIsNotEmpty,
				ConflictsWith: []string{"file_path"},
			},

			"display_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDataShareDataSetDataLakeGen2Create(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.DataSetClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	shareId, err := share.ParseShareID(d.Get("share_id").(string))
	if err != nil {
		return err
	}
	id := dataset.NewDataSetID(shareId.SubscriptionId, shareId.ResourceGroupName, shareId.AccountName, shareId.ShareName, d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_data_share_dataset_data_lake_gen2", id.ID())
	}

	strId, err := commonids.ParseStorageAccountID(d.Get("storage_account_id").(string))
	if err != nil {
		return err
	}

	var dataSet dataset.DataSet

	if filePath, ok := d.GetOk("file_path"); ok {
		dataSet = dataset.ADLSGen2FileDataSet{
			Properties: dataset.ADLSGen2FileProperties{
				StorageAccountName: strId.StorageAccountName,
				ResourceGroup:      strId.ResourceGroupName,
				SubscriptionId:     strId.SubscriptionId,
				FileSystem:         d.Get("file_system_name").(string),
				FilePath:           filePath.(string),
			},
		}
	} else if folderPath, ok := d.GetOk("folder_path"); ok {
		dataSet = dataset.ADLSGen2FolderDataSet{
			Properties: dataset.ADLSGen2FolderProperties{
				StorageAccountName: strId.StorageAccountName,
				ResourceGroup:      strId.ResourceGroupName,
				SubscriptionId:     strId.SubscriptionId,
				FileSystem:         d.Get("file_system_name").(string),
				FolderPath:         folderPath.(string),
			},
		}
	} else {
		dataSet = dataset.ADLSGen2FileSystemDataSet{
			Properties: dataset.ADLSGen2FileSystemProperties{
				StorageAccountName: strId.StorageAccountName,
				ResourceGroup:      strId.ResourceGroupName,
				SubscriptionId:     strId.SubscriptionId,
				FileSystem:         d.Get("file_system_name").(string),
			},
		}
	}

	if _, err := client.Create(ctx, id, dataSet); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}
	return resourceDataShareDataSetDataLakeGen2Read(d, meta)
}

func resourceDataShareDataSetDataLakeGen2Read(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.DataSetClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dataset.ParseDataSetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] DataShare %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	d.Set("name", id.DataSetName)

	shareId := share.NewShareID(id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.ShareName)
	d.Set("share_id", shareId.ID())

	if model := resp.Model; model != nil {
		if ds, ok := model.(dataset.ADLSGen2FileDataSet); ok {
			props := ds.Properties
			d.Set("storage_account_id", commonids.NewStorageAccountID(props.SubscriptionId, props.ResourceGroup, props.StorageAccountName).ID())
			d.Set("file_system_name", props.FileSystem)
			d.Set("file_path", props.FilePath)
			d.Set("display_name", props.DataSetId)
		} else if ds, ok := model.(dataset.ADLSGen2FolderDataSet); ok {
			props := ds.Properties
			d.Set("storage_account_id", commonids.NewStorageAccountID(props.SubscriptionId, props.ResourceGroup, props.StorageAccountName).ID())
			d.Set("file_system_name", props.FileSystem)
			d.Set("folder_path", props.FolderPath)
			d.Set("display_name", props.DataSetId)
		} else if ds, ok := model.(dataset.ADLSGen2FileSystemDataSet); ok {
			props := ds.Properties
			d.Set("storage_account_id", commonids.NewStorageAccountID(props.SubscriptionId, props.ResourceGroup, props.StorageAccountName).ID())
			d.Set("file_system_name", props.FileSystem)
			d.Set("display_name", props.DataSetId)
		} else {
			return fmt.Errorf("%s is not a datalake store gen2 dataset", *id)
		}
	}
	return pluginsdk.SetResourceIdentityData(d, id)
}

func resourceDataShareDataSetDataLakeGen2Delete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataShare.DataSetClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dataset.ParseDataSetID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
