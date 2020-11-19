package media

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/mediaservices/mgmt/2018-07-01/media"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/media/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMediaAssets() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMediaAssetsCreateUpdate,
		Read:   resourceArmMediaAssetsRead,
		Update: resourceArmMediaAssetsCreateUpdate,
		Delete: resourceArmMediaAssetsDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.MediaAssetsID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-z0-9]{3,24}$"),
					"Asset name must be 3 - 24 characters long, contain only lowercase letters and numbers.",
				),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"media_services_account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-z0-9]{3,24}$"),
					"Media Services Account name must be 3 - 24 characters long, contain only lowercase letters and numbers.",
				),
			},

			"alternate_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"container": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"storage_account_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceArmMediaAssetsCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.AssetsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	assetName := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	accountName := d.Get("media_services_account_name").(string)
	description := d.Get("description").(string)
	alternateId := d.Get("alternate_id").(string)
	container := d.Get("container").(string)
	storageAccountName := d.Get("storage_account_name").(string)

	parameters := media.Asset{
		AssetProperties: &media.AssetProperties{
			Description:        &description,
			AlternateID:        &alternateId,
			StorageAccountName: &storageAccountName,
		},
	}

	if container != "" {
		parameters.Container = &container
	}

	if _, e := client.CreateOrUpdate(ctx, resourceGroup, accountName, assetName, parameters); e != nil {
		return fmt.Errorf("Error creating Asset %q in Media Services Account %q (Resource Group %q): %+v", assetName, accountName, resourceGroup, e)
	}

	asset, err := client.Get(ctx, resourceGroup, accountName, assetName)
	if err != nil {
		return fmt.Errorf("Error retrieving Asset %q in Media Services Account %q (Resource Group %q): %+v", assetName, accountName, resourceGroup, err)
	}
	log.Printf(*asset.ID)
	d.SetId(*asset.ID)

	return resourceArmMediaAssetsRead(d, meta)
}

func resourceArmMediaAssetsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.AssetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MediaAssetsID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.AccountName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Asset %q was not found in Media Services Account %q and Resource Group %q - removing from state", id.Name, id.AccountName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Asset %q from Media Services Account %q (Resource Group %q): %+v", id.Name, id.AccountName, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("media_services_account_name", id.AccountName)
	d.Set("description", resp.AssetProperties.Description)
	d.Set("alternate_id", resp.AssetProperties.AlternateID)
	d.Set("container", resp.AssetProperties.Container)
	d.Set("storage_account_name", resp.AssetProperties.StorageAccountName)

	return nil
}

func resourceArmMediaAssetsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.AssetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MediaAssetsID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.AccountName, id.Name)
	if err != nil {
		if response.WasNotFound(resp.Response) {
			return nil
		}
		return fmt.Errorf("Error issuing AzureRM delete request for Asset '%s': %+v", id.Name, err)
	}

	return nil
}
