package network

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-02-01/network"
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

func resourceVirtualHubRoute() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualHubRouteCreateUpdate,
		Read:   resourceVirtualHubRouteRead,
		Update: resourceVirtualHubRouteCreateUpdate,
		Delete: resourceVirtualHubRouteDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, name, err := ParseHubRouteId(id)

			if err != nil {
				return err
			}

			if len(name) == 0 {
				return fmt.Errorf("route name is empty")
			}

			return nil
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

func resourceVirtualHubRouteCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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
			return fmt.Errorf("checking for present of existing HubRouteTable %q (Resource Group %q / Virtual Hub %q): %+v", routeTableId.Name, routeTableId.ResourceGroup, routeTableId.VirtualHubName, err)
		}
	}

	name := d.Get("name").(string)
	id := HubRouteID(routeTable, name)

	if d.IsNewResource() {
		for _, r := range *routeTable.Routes {
			if *r.Name == name {
				return tf.ImportAsExistsError("azurerm_virtual_hub_route", id)
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
		return fmt.Errorf("creating/updating Route %q in HubRouteTable %q (Resource Group %q / Virtual Hub %q): %+v", name, routeTableId.Name, routeTableId.ResourceGroup, routeTableId.Name, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creating/updating future for Route %q in HubRouteTable %q (Resource Group %q / Virtual Hub %q): %+v", name, routeTableId.Name, routeTableId.ResourceGroup, routeTableId.Name, err)
	}

	resp, err := client.Get(ctx, routeTableId.ResourceGroup, routeTableId.VirtualHubName, routeTableId.Name)
	if err != nil {
		return fmt.Errorf("retrieving HubRouteTable %q (Resource Group %q / Virtual Hub %q): %+v", name, routeTableId.ResourceGroup, routeTableId.Name, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for HubRouteTable %q (Resource Group %q / Virtual Hub %q) ID", name, routeTableId.ResourceGroup, routeTableId.Name)
	}

	d.SetId(id)

	return resourceVirtualHubRouteRead(d, meta)
}

func resourceVirtualHubRouteRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.HubRouteTableClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	routeTableId, name, err := ParseHubRouteId(d.Id())

	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, routeTableId.ResourceGroup, routeTableId.VirtualHubName, routeTableId.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Virtual Hub Route Table %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving HubRouteTable %q (Resource Group %q / Virtual Hub %q): %+v", routeTableId.Name, routeTableId.ResourceGroup, routeTableId.VirtualHubName, err)
	}

	if props := resp.HubRouteTableProperties; props != nil {
		found := false
		for _, r := range *props.Routes {
			if *r.Name == name {
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

	d.Set("name", name)
	d.Set("route_table_id", routeTableId.ID())

	return nil
}

func resourceVirtualHubRouteDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.HubRouteTableClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	routeTableId, name, err := ParseHubRouteId(d.Id())

	if err != nil {
		return err
	}

	locks.ByName(routeTableId.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(routeTableId.VirtualHubName, virtualHubResourceName)

	// get latest list of routes
	routeTable, err := client.Get(ctx, routeTableId.ResourceGroup, routeTableId.VirtualHubName, routeTableId.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(routeTable.Response) {
			return fmt.Errorf("checking for present of existing HubRouteTable %q (Resource Group %q / Virtual Hub %q): %+v", routeTableId.Name, routeTableId.ResourceGroup, routeTableId.VirtualHubName, err)
		}
	}

	if props := routeTable.HubRouteTableProperties; props != nil {
		if props.Routes != nil {
			routes := *props.Routes
			removeIndex := -1

			for i, r := range routes {
				if *r.Name == name {
					removeIndex = i
					break
				}
			}

			if removeIndex > -1 {
				if len(routes) == 1 && removeIndex == 0 {
					routes = nil
				} else {
					routes[removeIndex] = routes[len(routes)-1]
					routes = routes[:len(routes)-1]
				}
			}

			props.Routes = &routes
		}
	}

	future, err := client.CreateOrUpdate(ctx, routeTableId.ResourceGroup, routeTableId.VirtualHubName, routeTableId.Name, routeTable)
	if err != nil {
		return fmt.Errorf("removing Route %q from HubRouteTable %q (Resource Group %q / Virtual Hub %q): %+v", name, routeTableId.Name, routeTableId.ResourceGroup, routeTableId.Name, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on removing Route %q from future for HubRouteTable %q (Resource Group %q / Virtual Hub %q): %+v", name, routeTableId.Name, routeTableId.ResourceGroup, routeTableId.Name, err)
	}

	return nil
}

// As we are making a "virtual" sub-resource, the id stored in state needs to contain the name
// This ensures that terraform import works
// Note: This resource id DOES NOT exist in Azure.
// ID format: <route table id>/<route name>
// For example:
// Route Table ID: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/virtualHubs/virtualHub1/hubRouteTables/routeTable1
// Route Name: route1
// Resulting ID: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/virtualHubs/virtualHub1/hubRouteTables/routeTable1/route1

func HubRouteID(hub network.HubRouteTable, routeName string) string {
	return fmt.Sprintf("%s/%s", *hub.ID, routeName)
}

func ParseHubRouteId(input string) (*parse.HubRouteTableId, string, error) {
	i := strings.LastIndex(input, "/")
	routeTableID := input[:i]
	name := input[i+1:]

	routeTable, err := parse.HubRouteTableID(routeTableID)
	if err != nil {
		return nil, "", err
	}

	return routeTable, name, nil
}
