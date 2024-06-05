// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/publicipprefixes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePublicIpPrefix() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePublicIpPrefixCreate,
		Read:   resourcePublicIpPrefixRead,
		Update: resourcePublicIpPrefixUpdate,
		Delete: resourcePublicIpPrefixDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := publicipprefixes.ParsePublicIPPrefixID(id)
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
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"sku": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(publicipprefixes.PublicIPPrefixSkuNameStandard),
				ValidateFunc: validation.StringInSlice([]string{
					string(publicipprefixes.PublicIPPrefixSkuNameStandard),
				}, false),
			},

			"prefix_length": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      28,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 127),
			},

			"ip_version": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(publicipprefixes.IPVersionIPvFour),
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(publicipprefixes.IPVersionIPvFour),
					string(publicipprefixes.IPVersionIPvSix),
				}, false),
			},

			"zones": commonschema.ZonesMultipleOptionalForceNew(),

			"ip_prefix": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourcePublicIpPrefixCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PublicIPPrefixes
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Public IP Prefix creation.")

	id := publicipprefixes.NewPublicIPPrefixID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id, publicipprefixes.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_public_ip_prefix", id.ID())
	}

	publicIpPrefix := publicipprefixes.PublicIPPrefix{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Sku: &publicipprefixes.PublicIPPrefixSku{
			Name: pointer.To(publicipprefixes.PublicIPPrefixSkuName(d.Get("sku").(string))),
		},
		Properties: &publicipprefixes.PublicIPPrefixPropertiesFormat{
			PrefixLength:           pointer.To(int64(d.Get("prefix_length").(int))),
			PublicIPAddressVersion: pointer.To(publicipprefixes.IPVersion(d.Get("ip_version").(string))),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	zones := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
	if len(zones) > 0 {
		publicIpPrefix.Zones = &zones
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, publicIpPrefix); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourcePublicIpPrefixRead(d, meta)
}

func resourcePublicIpPrefixUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PublicIPPrefixes
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if d.HasChange("tags") {

		id, err := publicipprefixes.ParsePublicIPPrefixID(d.Id())
		if err != nil {
			return err
		}

		payload := publicipprefixes.TagsObject{
			Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
		}

		if _, err = client.UpdateTags(ctx, *id, payload); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourcePublicIpPrefixRead(d, meta)
}

func resourcePublicIpPrefixRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PublicIPPrefixes
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := publicipprefixes.ParsePublicIPPrefixID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, publicipprefixes.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.PublicIPPrefixName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		d.Set("zones", zones.FlattenUntyped(model.Zones))
		skuName := ""
		if sku := model.Sku; sku != nil {
			skuName = string(pointer.From(sku.Name))
		}
		d.Set("sku", skuName)
		if props := model.Properties; props != nil {
			d.Set("prefix_length", props.PrefixLength)
			d.Set("ip_prefix", props.IPPrefix)
			if version := props.PublicIPAddressVersion; version != nil {
				d.Set("ip_version", string(*version))
			}
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourcePublicIpPrefixDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PublicIPPrefixes
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := publicipprefixes.ParsePublicIPPrefixID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
