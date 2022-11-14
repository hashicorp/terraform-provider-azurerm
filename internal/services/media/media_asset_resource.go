package media

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/mediaservices/mgmt/2021-05-01/media"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2020-05-01/assetsandassetfilters"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/media/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMediaAsset() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMediaAssetCreateUpdate,
		Read:   resourceMediaAssetRead,
		Update: resourceMediaAssetCreateUpdate,
		Delete: resourceMediaAssetDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := assetsandassetfilters.ParseAssetID(id)
			return err
		}),

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.AssetV0ToV1{},
		}),
		SchemaVersion: 1,

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-zA-Z0-9]{1,128}$"),
					"Asset name must be 1 - 128 characters long, contain only letters, hyphen and numbers.",
				),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"media_services_account_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^[-a-z0-9]{3,24}$"),
					"Media Services Account name must be 3 - 24 characters long, contain only lowercase letters and numbers.",
				),
			},

			"alternate_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"container": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validate.StorageContainerName,
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"storage_account_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile("^([a-z0-9]{3,24})$"),
					"Storage Account Name can only consist of lowercase letters and numbers, and must be between 3 and 24 characters long.",
				),
			},
		},
	}
}

func resourceMediaAssetCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.AssetsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := assetsandassetfilters.NewAssetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("media_services_account_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroupName, id.AccountName, id.AssetName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_media_asset", id.ID())
		}
	}

	parameters := media.Asset{
		AssetProperties: &media.AssetProperties{
			Description: utils.String(d.Get("description").(string)),
		},
	}

	if v, ok := d.GetOk("container"); ok {
		parameters.Container = utils.String(v.(string))
	}

	if v, ok := d.GetOk("alternate_id"); ok {
		parameters.AlternateID = utils.String(v.(string))
	}

	if v, ok := d.GetOk("storage_account_name"); ok {
		parameters.StorageAccountName = utils.String(v.(string))
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroupName, id.AccountName, id.AssetName, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceMediaAssetRead(d, meta)
}

func resourceMediaAssetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.AssetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := assetsandassetfilters.ParseAssetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroupName, id.AccountName, id.AssetName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.AssetName)
	d.Set("media_services_account_name", id.AccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if props := resp.AssetProperties; props != nil {
		d.Set("description", props.Description)
		d.Set("alternate_id", props.AlternateID)
		d.Set("container", props.Container)
		d.Set("storage_account_name", props.StorageAccountName)
	}

	return nil
}

func resourceMediaAssetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.AssetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := assetsandassetfilters.ParseAssetID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroupName, id.AccountName, id.AssetName); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
