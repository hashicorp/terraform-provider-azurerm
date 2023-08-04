// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mixedreality

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mixedreality/2021-01-01/resource"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSpatialAnchorsAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSpatialAnchorsAccountCreate,
		Read:   resourceSpatialAnchorsAccountRead,
		Update: resourceSpatialAnchorsAccountUpdate,
		Delete: resourceSpatialAnchorsAccountDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := resource.ParseSpatialAnchorsAccountID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[-\w._()]+$`),
					"Spatial Anchors Account name must be 1 - 90 characters long, contain only word characters and underscores.",
				),
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"account_domain": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"account_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceSpatialAnchorsAccountCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MixedReality.SpatialAnchorsAccountClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := resource.NewSpatialAnchorsAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.SpatialAnchorsAccountsGet(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_spatial_anchors_account", id.ID())
	}

	account := resource.SpatialAnchorsAccount{
		Location: location.Normalize(d.Get("location").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.SpatialAnchorsAccountsCreate(ctx, id, account); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceSpatialAnchorsAccountRead(d, meta)
}

func resourceSpatialAnchorsAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MixedReality.SpatialAnchorsAccountClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := resource.ParseSpatialAnchorsAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.SpatialAnchorsAccountsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.SpatialAnchorsAccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			d.Set("account_domain", props.AccountDomain)
			d.Set("account_id", props.AccountId)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceSpatialAnchorsAccountUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MixedReality.SpatialAnchorsAccountClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := resource.ParseSpatialAnchorsAccountID(d.Id())
	if err != nil {
		return err
	}

	account := resource.SpatialAnchorsAccount{
		Location: location.Normalize(d.Get("location").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.SpatialAnchorsAccountsUpdate(ctx, *id, account); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceSpatialAnchorsAccountRead(d, meta)
}

func resourceSpatialAnchorsAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MixedReality.SpatialAnchorsAccountClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := resource.ParseSpatialAnchorsAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.SpatialAnchorsAccountsDelete(ctx, *id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
