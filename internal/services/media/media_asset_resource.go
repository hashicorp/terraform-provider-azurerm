// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package media

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/assetsandassetfilters"
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

		DeprecationMessage: azureMediaRetirementMessage,

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
	client := meta.(*clients.Client).Media.V20220801Client.AssetsAndAssetFilters
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := assetsandassetfilters.NewAssetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("media_services_account_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.AssetsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_media_asset", id.ID())
		}
	}

	payload := assetsandassetfilters.Asset{
		Properties: &assetsandassetfilters.AssetProperties{
			Description: utils.String(d.Get("description").(string)),
		},
	}

	if v, ok := d.GetOk("container"); ok {
		payload.Properties.Container = utils.String(v.(string))
	}

	if v, ok := d.GetOk("alternate_id"); ok {
		payload.Properties.AlternateId = utils.String(v.(string))
	}

	if v, ok := d.GetOk("storage_account_name"); ok {
		payload.Properties.StorageAccountName = utils.String(v.(string))
	}

	if _, err := client.AssetsCreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceMediaAssetRead(d, meta)
}

func resourceMediaAssetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20220801Client.AssetsAndAssetFilters
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := assetsandassetfilters.ParseAssetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.AssetsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.AssetName)
	d.Set("media_services_account_name", id.MediaServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("description", props.Description)
			d.Set("alternate_id", props.AlternateId)
			d.Set("container", props.Container)
			d.Set("storage_account_name", props.StorageAccountName)
		}
	}

	return nil
}

func resourceMediaAssetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Media.V20220801Client.AssetsAndAssetFilters
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := assetsandassetfilters.ParseAssetID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.AssetsDelete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
