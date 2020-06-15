package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-03-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var virtualHubRouteTableResourceName = "azurerm_virtual_hub_route_table"

// https://docs.microsoft.com/en-us/rest/api/virtualwan/virtualhubroutetablev2s/createorupdate

func resourceArmVirtualHubRouteTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVirtualHubRouteTableCreateUpdate,
		Read:   resourceArmVirtualHubRouteTableRead,
		Update: resourceArmVirtualHubRouteTableCreateUpdate,
		Delete: resourceArmVirtualHubRouteTableDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"virtual_hub_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"attached_connections": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"route": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"destinations": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"next_hop_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"next_hops": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func resourceArmVirtualHubRouteTableCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualHubRouteTableV2sClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM VirtualHub Route Table creation")

	name := d.Get("name").(string)
	virtualHubName := d.Get("virtual_hub_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, virtualHubName, name)

		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing VirtualHub Route Table %q (Resource Group %q, Hub %q): %s", name, resourceGroup, virtualHubName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_virtual_hub_route_table", *existing.ID)
		}
	}

	locks.ByName(virtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(virtualHubName, virtualHubResourceName)

	locks.ByName(name, virtualHubRouteTableResourceName)
	defer locks.UnlockByName(name, virtualHubRouteTableResourceName)

	parameters := network.VirtualHubRouteTableV2{
		VirtualHubRouteTableV2Properties: &network.VirtualHubRouteTableV2Properties{},
	}

	attachedConnections := make([]string, 0)
	rawAttachedConnections := d.Get("attached_connections").([]interface{})
	for _, rawAttachedConnection := range rawAttachedConnections {
		data := rawAttachedConnection.(string)
		attachedConnections = append(attachedConnections, data)
	}
	parameters.AttachedConnections = &attachedConnections

	routes := make([]network.VirtualHubRouteV2, 0)
	rawRoutes := d.Get("route").([]interface{})
	for _, rawRoute := range rawRoutes {
		raw := rawRoute.(map[string]interface{})

		destinationType := raw["destination_type"].(string)
		destinations := make([]string, 0)
		for _, rawDestination := range raw["destinations"].([]interface{}) {
			destination := rawDestination.(string)
			destinations = append(destinations, destination)
		}

		nextHopType := raw["next_hop_type"].(string)
		nextHops := make([]string, 0)
		for _, rawNextHop := range raw["next_hops"].([]interface{}) {
			nextHop := rawNextHop.(string)
			nextHops = append(nextHops, nextHop)
		}

		route := network.VirtualHubRouteV2{
			DestinationType: &destinationType,
			Destinations:    &destinations,
			NextHopType:     &nextHopType,
			NextHops:        &nextHops,
		}

		routes = append(routes, route)
	}
	parameters.Routes = &routes

	future, err := client.CreateOrUpdate(ctx, resourceGroup, virtualHubName, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating/updating VirtualHub Route Table %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation/update of VirtualHub Route Table %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, virtualHubName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving VirtualHub Route Table %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read VirtualHub Route Table %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmVirtualHubRouteTableRead(d, meta)
}

func resourceArmVirtualHubRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualHubRouteTableV2sClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["routeTables"]
	virtualHubName := id.Path["virtualHubs"]

	read, err := client.Get(ctx, resourceGroup, virtualHubName, name)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			log.Printf("[DEBUG] VirtualHub Route Table %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on VirtualHub Route Table %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", read.Name)
	d.Set("virtual_hub_name", virtualHubName)
	d.Set("resource_group_name", resourceGroup)

	if attachedConnections := read.AttachedConnections; attachedConnections != nil {
		d.Set("attached_connections", attachedConnections)
	}

	if routes := read.Routes; routes != nil {
		rawRoutes := make([]interface{}, 0)

		for _, route := range *routes {
			rawRoute := make(map[string]interface{})

			rawRoute["destination_type"] = *route.DestinationType
			rawRoute["destinations"] = *route.Destinations
			rawRoute["next_hop_type"] = *route.NextHopType
			rawRoute["next_hops"] = *route.NextHops

			rawRoutes = append(rawRoutes, rawRoute)
		}

		if len(rawRoutes) > 0 {
			d.Set("route", rawRoutes)
		}
	}

	return nil
}

func resourceArmVirtualHubRouteTableDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualHubRouteTableV2sClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["routeTables"]
	virtualHubName := id.Path["virtualHubs"]

	read, err := client.Get(ctx, resourceGroup, virtualHubName, name)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			// deleted outside of TF
			log.Printf("[DEBUG] VirtualHub Route Table %q was not found in Resource Group %q - assuming removed!", name, resourceGroup)
			return nil
		}

		return fmt.Errorf("Error retrieving VirtualHub Route Table %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	locks.ByName(name, virtualHubRouteTableResourceName)
	defer locks.UnlockByName(name, virtualHubRouteTableResourceName)

	locks.ByName(virtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(virtualHubName, virtualHubResourceName)

	future, err := client.Delete(ctx, resourceGroup, virtualHubName, name)
	if err != nil {
		return fmt.Errorf("Error deleting VirtualHub Route Table %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the deletion of Virtual Hub Route Table %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return err
}
