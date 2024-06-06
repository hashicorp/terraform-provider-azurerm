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
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/natgateways"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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
			_, err := natgateways.ParseNatGatewayID(id)
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
			Default:  string(natgateways.NatGatewaySkuNameStandard),
			ValidateFunc: validation.StringInSlice([]string{
				string(natgateways.NatGatewaySkuNameStandard),
			}, false),
		},

		"zones": commonschema.ZonesMultipleOptionalForceNew(),

		"resource_guid": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.Tags(),
	}
}

func resourceNatGatewayCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.NatGateways
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := natgateways.NewNatGatewayID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByName(id.NatGatewayName, natGatewayResourceName)
	defer locks.UnlockByName(id.NatGatewayName, natGatewayResourceName)

	resp, err := client.Get(ctx, id, natgateways.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("checking for present of existing %s: %+v", id, err)
		}
	}
	if resp.Model != nil && resp.Model.Id != nil && *resp.Model.Id != "" {
		return tf.ImportAsExistsError("azurerm_nat_gateway", id.ID())
	}

	parameters := natgateways.NatGateway{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &natgateways.NatGatewayPropertiesFormat{
			IdleTimeoutInMinutes: pointer.To(int64(d.Get("idle_timeout_in_minutes").(int))),
		},
		Sku: &natgateways.NatGatewaySku{
			Name: pointer.To(natgateways.NatGatewaySkuName(d.Get("sku_name").(string))),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	zones := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
	if len(zones) > 0 {
		parameters.Zones = &zones
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceNatGatewayRead(d, meta)
}

func resourceNatGatewayUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.NatGateways
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := natgateways.ParseNatGatewayID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.NatGatewayName, natGatewayResourceName)
	defer locks.UnlockByName(id.NatGatewayName, natGatewayResourceName)

	existing, err := client.Get(ctx, *id, natgateways.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("%s was not found", *id)
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}
	props := *existing.Model.Properties

	// intentionally building a new object rather than reusing due to the additional read-only fields
	payload := natgateways.NatGateway{
		Location: existing.Model.Location,
		Properties: &natgateways.NatGatewayPropertiesFormat{
			IdleTimeoutInMinutes: props.IdleTimeoutInMinutes,
			PublicIPAddresses:    props.PublicIPAddresses, // note: these can be managed via the separate resource
			PublicIPPrefixes:     props.PublicIPPrefixes,
		},
		Sku:   existing.Model.Sku,
		Tags:  existing.Model.Tags,
		Zones: existing.Model.Zones,
	}

	if d.HasChange("idle_timeout_in_minutes") {
		timeout := d.Get("idle_timeout_in_minutes").(int)
		payload.Properties.IdleTimeoutInMinutes = pointer.To(int64(timeout))
	}

	if d.HasChange("sku_name") {
		payload.Sku = &natgateways.NatGatewaySku{
			Name: pointer.To(natgateways.NatGatewaySkuName(d.Get("sku_name").(string))),
		}
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceNatGatewayRead(d, meta)
}

func resourceNatGatewayRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.NatGateways
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := natgateways.ParseNatGatewayID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, natgateways.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s does not exist - removing from state", id.ID())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.NatGatewayName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		sku := ""
		if model.Sku != nil {
			sku = string(pointer.From(model.Sku.Name))
		}
		d.Set("sku_name", sku)
		d.Set("zones", zones.FlattenUntyped(model.Zones))
		if props := model.Properties; props != nil {
			d.Set("idle_timeout_in_minutes", props.IdleTimeoutInMinutes)
			d.Set("resource_guid", props.ResourceGuid)
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceNatGatewayDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.NatGateways
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := natgateways.ParseNatGatewayID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.NatGatewayName, natGatewayResourceName)
	defer locks.UnlockByName(id.NatGatewayName, natGatewayResourceName)

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
