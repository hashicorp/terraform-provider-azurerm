// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package maps

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maps/2023-06-01/accounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maps/2023-06-01/creators"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMapsCreator() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMapsCreatorCreateUpdate,
		Read:   resourceMapsCreatorRead,
		Update: resourceMapsCreatorUpdate,
		Delete: resourceMapsCreatorDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := creators.ParseCreatorID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"maps_account_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: accounts.ValidateAccountID,
			},

			"location": commonschema.Location(),

			"storage_units": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 100),
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceMapsCreatorCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maps.CreatorsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	accountId, err := accounts.ParseAccountID(d.Get("maps_account_id").(string))
	if err != nil {
		return err
	}

	id := creators.NewCreatorID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.AccountName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_maps_creator", id.ID())
		}
	}

	props := creators.Creator{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: creators.CreatorProperties{
			StorageUnits: int64(d.Get("storage_units").(int)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}
	if _, err := client.CreateOrUpdate(ctx, id, props); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceMapsCreatorRead(d, meta)
}

func resourceMapsCreatorRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maps.CreatorsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := creators.ParseCreatorID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.CreatorName)
	d.Set("maps_account_id", accounts.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName).ID())

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		props := model.Properties
		d.Set("storage_units", props.StorageUnits)
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceMapsCreatorUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maps.CreatorsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := creators.ParseCreatorID(d.Id())
	if err != nil {
		return err
	}

	props := creators.CreatorUpdateParameters{
		Properties: &creators.CreatorProperties{
			StorageUnits: int64(d.Get("storage_units").(int)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.Update(ctx, *id, props); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceMapsCreatorRead(d, meta)
}

func resourceMapsCreatorDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Maps.CreatorsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := creators.ParseCreatorID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
