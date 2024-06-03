// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/routes"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRoute() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceRouteCreate,
		Read:   resourceRouteRead,
		Update: resourceRouteUpdate,
		Delete: resourceRouteDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := routes.ParseRouteID(id)
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
				ValidateFunc: validate.RouteName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"route_table_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.RouteTableName,
			},

			"address_prefix": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"next_hop_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(routes.RouteNextHopTypeVirtualNetworkGateway),
					string(routes.RouteNextHopTypeVnetLocal),
					string(routes.RouteNextHopTypeInternet),
					string(routes.RouteNextHopTypeVirtualAppliance),
					string(routes.RouteNextHopTypeNone),
				}, false),
			},

			"next_hop_in_ip_address": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceRouteCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Routes
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	addressPrefix := d.Get("address_prefix").(string)
	nextHopType := d.Get("next_hop_type").(string)

	id := routes.NewRouteID(subscriptionId, d.Get("resource_group_name").(string), d.Get("route_table_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_route", id.ID())
	}

	locks.ByName(id.RouteTableName, routeTableResourceName)
	defer locks.UnlockByName(id.RouteTableName, routeTableResourceName)

	route := routes.Route{
		Name: pointer.To(id.RouteName),
		Properties: &routes.RoutePropertiesFormat{
			AddressPrefix: &addressPrefix,
			NextHopType:   routes.RouteNextHopType(nextHopType),
		},
	}

	if v, ok := d.GetOk("next_hop_in_ip_address"); ok {
		route.Properties.NextHopIPAddress = pointer.To(v.(string))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, route); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceRouteRead(d, meta)
}

func resourceRouteUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Routes
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := routes.ParseRouteID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := existing.Model

	locks.ByName(id.RouteTableName, routeTableResourceName)
	defer locks.UnlockByName(id.RouteTableName, routeTableResourceName)

	if d.HasChange("address_prefix") {
		payload.Properties.AddressPrefix = pointer.To(d.Get("address_prefix").(string))
	}

	if d.HasChange("next_hop_type") {
		payload.Properties.NextHopType = routes.RouteNextHopType(d.Get("next_hop_type").(string))
	}

	if d.HasChange("next_hop_in_ip_address") {
		payload.Properties.NextHopIPAddress = pointer.To(d.Get("next_hop_in_ip_address").(string))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceRouteRead(d, meta)
}

func resourceRouteRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Routes
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := routes.ParseRouteID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.RouteName)
	d.Set("route_table_name", id.RouteTableName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("address_prefix", props.AddressPrefix)
			d.Set("next_hop_type", string(props.NextHopType))
			d.Set("next_hop_in_ip_address", props.NextHopIPAddress)
		}
	}

	return nil
}

func resourceRouteDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Routes
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := routes.ParseRouteID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.RouteTableName, routeTableResourceName)
	defer locks.UnlockByName(id.RouteTableName, routeTableResourceName)

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
