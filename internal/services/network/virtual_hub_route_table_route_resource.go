package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-05-01/network"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceVirtualHubRouteTableRoute() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualHubRouteTableRouteCreateUpdate,
		Read:   resourceVirtualHubRouteTableRouteRead,
		Update: resourceVirtualHubRouteTableRouteCreateUpdate,
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
				ValidateFunc: networkValidate.HubRouteTableID,
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

func resourceVirtualHubRouteTableRouteCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.HubRouteTableClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	routeTableId, err := parse.HubRouteTableID(d.Get("route_table_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(routeTableId.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(routeTableId.VirtualHubName, virtualHubResourceName)

	routeTable, err := client.Get(ctx, routeTableId.ResourceGroup, routeTableId.VirtualHubName, routeTableId.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(routeTable.Response) {
			return fmt.Errorf("checking for existing %s: %+v", routeTableId, err)
		}

		return fmt.Errorf("retrieving %s: %+v", routeTableId, err)
	}

	name := d.Get("name").(string)
	id := parse.NewHubRouteTableRouteID(routeTableId.SubscriptionId, routeTableId.ResourceGroup, routeTableId.VirtualHubName, routeTableId.Name, name)

	if d.IsNewResource() {
		for _, r := range *routeTable.Routes {
			if *r.Name == name {
				return tf.ImportAsExistsError("azurerm_virtual_hub_route_table_route", id.ID())
			}
		}

		routes := *routeTable.Routes
		result := network.HubRoute{
			Name:            utils.String(d.Get("name").(string)),
			DestinationType: utils.String(d.Get("destinations_type").(string)),
			Destinations:    utils.ExpandStringSlice(d.Get("destinations").(*pluginsdk.Set).List()),
			NextHopType:     utils.String(d.Get("next_hop_type").(string)),
			NextHop:         utils.String(d.Get("next_hop").(string)),
		}

		routes = append(routes, result)
		routeTable.Routes = &routes
	} else {
		for _, r := range *routeTable.Routes {
			if *r.Name == name {
				r.DestinationType = utils.String(d.Get("destinations_type").(string))
				r.Destinations = utils.ExpandStringSlice(d.Get("destinations").(*pluginsdk.Set).List())
				r.NextHopType = utils.String(d.Get("next_hop_type").(string))
				r.NextHop = utils.String(d.Get("next_hop").(string))
				break
			}
		}
	}

	future, err := client.CreateOrUpdate(ctx, routeTableId.ResourceGroup, routeTableId.VirtualHubName, routeTableId.Name, routeTable)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting to create/update %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualHubRouteTableRouteRead(d, meta)
}

func resourceVirtualHubRouteTableRouteRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.HubRouteTableClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	route, err := parse.HubRouteTableRouteID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, route.ResourceGroup, route.VirtualHubName, route.HubRouteTableName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Virtual Hub Route Table %q does not exist - removing route %s from state", route.HubRouteTableName, route.RouteName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", route, err)
	}

	if props := resp.HubRouteTableProperties; props != nil {
		found := false
		for _, r := range *props.Routes {
			if *r.Name == route.RouteName {
				found = true

				d.Set("destinations_type", r.DestinationType)
				d.Set("destinations", utils.FlattenStringSlice(r.Destinations))
				d.Set("next_hop_type", r.NextHopType)
				d.Set("next_hop", r.NextHop)

				break
			}
		}

		if !found {
			// could not find existing route by name
			d.SetId("")
			return nil
		}
	}

	d.Set("name", route.RouteName)
	routeTableID := parse.NewHubRouteTableID(route.SubscriptionId, route.ResourceGroup, route.VirtualHubName, route.HubRouteTableName)
	d.Set("route_table_id", routeTableID.ID())

	return nil
}

func resourceVirtualHubRouteTableRouteDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.HubRouteTableClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	route, err := parse.HubRouteTableRouteID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(route.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(route.VirtualHubName, virtualHubResourceName)

	// get latest list of routes
	routeTable, err := client.Get(ctx, route.ResourceGroup, route.VirtualHubName, route.HubRouteTableName)
	if err != nil {
		if !utils.ResponseWasNotFound(routeTable.Response) {
			// route table does not exist, therefore route does not exist
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", route, err)
	}

	if props := routeTable.HubRouteTableProperties; props != nil {
		if props.Routes != nil {
			routes := *props.Routes

			newRoutes := make([]network.HubRoute, 0)

			for _, r := range routes {
				if *r.Name != route.RouteName {
					newRoutes = append(newRoutes, r)
				}
			}

			props.Routes = &newRoutes
		}
	}

	future, err := client.CreateOrUpdate(ctx, route.ResourceGroup, route.VirtualHubName, route.HubRouteTableName, routeTable)
	if err != nil {
		return fmt.Errorf("removing %s: %+v", route, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting to remove %s: %+v", route, err)
	}

	return nil
}
