// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceVirtualHubRouteTableRoute() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualHubRouteTableRouteCreate,
		Read:   resourceVirtualHubRouteTableRouteRead,
		Update: resourceVirtualHubRouteTableRouteUpdate,
		Delete: resourceVirtualHubRouteTableRouteDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.HubRouteTableRouteID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"route_table_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: virtualwans.ValidateHubRouteTableID,
			},

			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"destinations": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"destinations_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"CIDR",
					"ResourceId",
					"Service",
				}, false),
			},

			"next_hop": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"next_hop_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "ResourceId",
				ValidateFunc: validation.StringInSlice([]string{
					"ResourceId",
				}, false),
			},
		},
	}
}

func resourceVirtualHubRouteTableRouteCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	routeTableId, err := virtualwans.ParseHubRouteTableID(d.Get("route_table_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(routeTableId.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(routeTableId.VirtualHubName, virtualHubResourceName)

	routeTable, err := client.HubRouteTablesGet(ctx, *routeTableId)
	if err != nil {
		if !response.WasNotFound(routeTable.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", routeTableId, err)
		}
		return fmt.Errorf("retrieving %s: %+v", routeTableId, err)
	}

	if routeTable.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", routeTableId)
	}
	if routeTable.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", routeTableId)
	}

	props := routeTable.Model.Properties

	name := d.Get("name").(string)
	id := parse.NewHubRouteTableRouteID(routeTableId.SubscriptionId, routeTableId.ResourceGroupName, routeTableId.VirtualHubName, routeTableId.HubRouteTableName, name)

	routes := make([]virtualwans.HubRoute, 0)
	if hubRoutes := props.Routes; hubRoutes != nil {
		for _, r := range *hubRoutes {
			if r.Name == name {
				return tf.ImportAsExistsError("azurerm_virtual_hub_route_table_route", id.ID())
			}
		}
		routes = *props.Routes

		result := virtualwans.HubRoute{
			Name:            d.Get("name").(string),
			DestinationType: d.Get("destinations_type").(string),
			Destinations:    pointer.From(utils.ExpandStringSlice(d.Get("destinations").(*pluginsdk.Set).List())),
			NextHopType:     d.Get("next_hop_type").(string),
			NextHop:         d.Get("next_hop").(string),
		}

		routes = append(routes, result)
	}

	routeTable.Model.Properties.Routes = pointer.To(routes)

	if err := client.HubRouteTablesCreateOrUpdateThenPoll(ctx, *routeTableId, *routeTable.Model); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualHubRouteTableRouteRead(d, meta)
}

func resourceVirtualHubRouteTableRouteUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	routeTableId, err := virtualwans.ParseHubRouteTableID(d.Get("route_table_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(routeTableId.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(routeTableId.VirtualHubName, virtualHubResourceName)

	routeTable, err := client.HubRouteTablesGet(ctx, *routeTableId)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", routeTableId, err)
	}

	if routeTable.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", routeTableId)
	}
	if routeTable.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", routeTableId)
	}

	props := routeTable.Model.Properties

	id, err := parse.HubRouteTableRouteID(d.Id())
	if err != nil {
		return err
	}

	routes := *props.Routes
	for i := range routes {
		if routes[i].Name == id.RouteName {
			if d.HasChange("destinations_type") {
				routes[i].DestinationType = d.Get("destinations_type").(string)
			}
			if d.HasChange("destinations") {
				routes[i].Destinations = pointer.From(utils.ExpandStringSlice(d.Get("destinations").(*pluginsdk.Set).List()))
			}
			if d.HasChange("next_hop_type") {
				routes[i].NextHopType = d.Get("next_hop_type").(string)
			}
			if d.HasChange("next_hop") {
				routes[i].NextHop = d.Get("next_hop").(string)
			}
			break
		}
	}

	routeTable.Model.Properties.Routes = pointer.To(routes)

	if err := client.HubRouteTablesCreateOrUpdateThenPoll(ctx, *routeTableId, *routeTable.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualHubRouteTableRouteRead(d, meta)
}

func resourceVirtualHubRouteTableRouteRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.HubRouteTableRouteID(d.Id())
	if err != nil {
		return err
	}

	routeTableId := virtualwans.NewHubRouteTableID(id.SubscriptionId, id.ResourceGroup, id.VirtualHubName, id.HubRouteTableName)

	resp, err := client.HubRouteTablesGet(ctx, routeTableId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", routeTableId.ID())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.RouteName)

	routeTableID := virtualwans.NewHubRouteTableID(id.SubscriptionId, id.ResourceGroup, id.VirtualHubName, id.HubRouteTableName)
	d.Set("route_table_id", routeTableID.ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			found := false
			for _, r := range *props.Routes {
				if r.Name == id.RouteName {
					found = true
					d.Set("destinations_type", r.DestinationType)
					d.Set("destinations", r.Destinations)
					d.Set("next_hop_type", r.NextHopType)
					d.Set("next_hop", r.NextHop)
					break
				}
			}

			if !found {
				// could not find existing id by name
				d.SetId("")
				return nil
			}
		}
	}

	return nil
}

func resourceVirtualHubRouteTableRouteDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.HubRouteTableRouteID(d.Id())
	if err != nil {
		return err
	}

	routeTableId := virtualwans.NewHubRouteTableID(id.SubscriptionId, id.ResourceGroup, id.VirtualHubName, id.HubRouteTableName)

	locks.ByName(id.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(id.VirtualHubName, virtualHubResourceName)

	// get latest list of routes
	routeTable, err := client.HubRouteTablesGet(ctx, routeTableId)
	if err != nil {
		if !response.WasNotFound(routeTable.HttpResponse) {
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if routeTable.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", routeTableId)
	}
	if routeTable.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", routeTableId)
	}

	props := routeTable.Model.Properties

	if props.Routes != nil {
		routes := *props.Routes

		newRoutes := make([]virtualwans.HubRoute, 0)
		for _, r := range routes {
			if r.Name != id.RouteName {
				newRoutes = append(newRoutes, r)
			}
		}
		props.Routes = &newRoutes

	}

	routeTable.Model.Properties = props

	if err := client.HubRouteTablesCreateOrUpdateThenPoll(ctx, routeTableId, *routeTable.Model); err != nil {
		return fmt.Errorf("removing %s: %+v", id, err)
	}

	return nil
}
