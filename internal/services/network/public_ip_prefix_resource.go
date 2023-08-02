// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

func resourcePublicIpPrefix() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePublicIpPrefixCreateUpdate,
		Read:   resourcePublicIpPrefixRead,
		Update: resourcePublicIpPrefixCreateUpdate,
		Delete: resourcePublicIpPrefixDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.PublicIpPrefixID(id)
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
				Default:  string(network.PublicIPPrefixSkuNameStandard),
				ValidateFunc: validation.StringInSlice([]string{
					string(network.PublicIPPrefixSkuNameStandard),
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
				Default:  string(network.IPVersionIPv4),
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.IPVersionIPv4),
					string(network.IPVersionIPv6),
				}, false),
			},

			"zones": commonschema.ZonesMultipleOptionalForceNew(),

			"ip_prefix": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourcePublicIpPrefixCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PublicIPPrefixesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Public IP Prefix creation.")

	id := parse.NewPublicIpPrefixID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_public_ip_prefix", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	sku := d.Get("sku").(string)
	prefixLength := d.Get("prefix_length").(int)
	ipVersion := d.Get("ip_version").(string)
	t := d.Get("tags").(map[string]interface{})

	publicIpPrefix := network.PublicIPPrefix{
		Location: &location,
		Sku: &network.PublicIPPrefixSku{
			Name: network.PublicIPPrefixSkuName(sku),
		},
		PublicIPPrefixPropertiesFormat: &network.PublicIPPrefixPropertiesFormat{
			PrefixLength:           utils.Int32(int32(prefixLength)),
			PublicIPAddressVersion: network.IPVersion(ipVersion),
		},
		Tags: tags.Expand(t),
	}

	zones := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
	if len(zones) > 0 {
		publicIpPrefix.Zones = &zones
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, publicIpPrefix)
	if err != nil {
		return fmt.Errorf("creating/Updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for completion of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourcePublicIpPrefixRead(d, meta)
}

func resourcePublicIpPrefixRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PublicIPPrefixesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PublicIpPrefixID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	d.Set("zones", zones.FlattenUntyped(resp.Zones))

	skuName := ""
	if sku := resp.Sku; sku != nil {
		skuName = string(sku.Name)
	}
	d.Set("sku", skuName)

	if props := resp.PublicIPPrefixPropertiesFormat; props != nil {
		d.Set("prefix_length", props.PrefixLength)
		d.Set("ip_prefix", props.IPPrefix)

		if version := props.PublicIPAddressVersion; version != "" {
			d.Set("ip_version", string(props.PublicIPAddressVersion))
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourcePublicIpPrefixDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PublicIPPrefixesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PublicIpPrefixID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}
