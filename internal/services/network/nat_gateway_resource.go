// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

var natGatewayResourceName = "azurerm_nat_gateway"

func resourceNatGateway() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNatGatewayCreate,
		Read:   resourceNatGatewayRead,
		Update: resourceNatGatewayUpdate,
		Delete: resourceNatGatewayDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.NatGatewayID(id)
			return err
		}),

		Schema: resourceNatGatewaySchema(),
	}
}

func resourceNatGatewaySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.NatGatewayName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"idle_timeout_in_minutes": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      4,
			ValidateFunc: validation.IntBetween(4, 120),
		},

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(network.NatGatewaySkuNameStandard),
			ValidateFunc: validation.StringInSlice([]string{
				string(network.NatGatewaySkuNameStandard),
			}, false),
		},

		"zones": commonschema.ZonesMultipleOptionalForceNew(),

		"resource_guid": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": tags.Schema(),
	}
}

func resourceNatGatewayCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NatGatewayClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewNatGatewayID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByName(id.Name, natGatewayResourceName)
	defer locks.UnlockByName(id.Name, natGatewayResourceName)

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if !utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("checking for present of existing %s: %+v", id, err)
		}
	}
	if resp.ID != nil && *resp.ID != "" {
		return tf.ImportAsExistsError("azurerm_nat_gateway", *resp.ID)
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	idleTimeoutInMinutes := d.Get("idle_timeout_in_minutes").(int)
	skuName := d.Get("sku_name").(string)
	t := d.Get("tags").(map[string]interface{})

	parameters := network.NatGateway{
		Location: utils.String(location),
		NatGatewayPropertiesFormat: &network.NatGatewayPropertiesFormat{
			IdleTimeoutInMinutes: utils.Int32(int32(idleTimeoutInMinutes)),
		},
		Sku: &network.NatGatewaySku{
			Name: network.NatGatewaySkuName(skuName),
		},
		Tags: tags.Expand(t),
	}

	zones := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
	if len(zones) > 0 {
		parameters.Zones = &zones
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, parameters)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceNatGatewayRead(d, meta)
}

func resourceNatGatewayUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NatGatewayClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NatGatewayID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Name, natGatewayResourceName)
	defer locks.UnlockByName(id.Name, natGatewayResourceName)

	existing, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("%s was not found!", *id)
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if existing.NatGatewayPropertiesFormat == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}
	props := *existing.NatGatewayPropertiesFormat

	// intentionally building a new object rather than reusing due to the additional read-only fields
	parameters := network.NatGateway{
		Location: existing.Location,
		NatGatewayPropertiesFormat: &network.NatGatewayPropertiesFormat{
			IdleTimeoutInMinutes: props.IdleTimeoutInMinutes,
			PublicIPAddresses:    props.PublicIPAddresses, // note: these can be managed via the separate resource
			PublicIPPrefixes:     props.PublicIPPrefixes,
		},
		Sku:   existing.Sku,
		Tags:  existing.Tags,
		Zones: existing.Zones,
	}

	if d.HasChange("idle_timeout_in_minutes") {
		timeout := d.Get("idle_timeout_in_minutes").(int)
		parameters.NatGatewayPropertiesFormat.IdleTimeoutInMinutes = utils.Int32(int32(timeout))
	}

	if d.HasChange("sku_name") {
		skuName := d.Get("sku_name").(string)
		parameters.Sku = &network.NatGatewaySku{
			Name: network.NatGatewaySkuName(skuName),
		}
	}

	if d.HasChange("tags") {
		t := d.Get("tags").(map[string]interface{})
		parameters.Tags = tags.Expand(t)
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, parameters)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", *id, err)
	}

	return resourceNatGatewayRead(d, meta)
}

func resourceNatGatewayRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NatGatewayClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NatGatewayID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] NAT Gateway %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", sku.Name)
	}

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.NatGatewayPropertiesFormat; props != nil {
		d.Set("idle_timeout_in_minutes", props.IdleTimeoutInMinutes)
		d.Set("resource_guid", props.ResourceGUID)
	}

	d.Set("zones", zones.FlattenUntyped(resp.Zones))

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceNatGatewayDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NatGatewayClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NatGatewayID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Name, natGatewayResourceName)
	defer locks.UnlockByName(id.Name, natGatewayResourceName)

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return nil
}
