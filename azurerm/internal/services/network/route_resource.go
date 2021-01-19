package network

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceRouteCreateUpdate,
		Read:   resourceRouteRead,
		Update: resourceRouteCreateUpdate,
		Delete: resourceRouteDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateRouteName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"route_table_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateRouteTableName,
			},

			"address_prefix": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"next_hop_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.RouteNextHopTypeVirtualNetworkGateway),
					string(network.RouteNextHopTypeVnetLocal),
					string(network.RouteNextHopTypeInternet),
					string(network.RouteNextHopTypeVirtualAppliance),
					string(network.RouteNextHopTypeNone),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"next_hop_in_ip_address": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceRouteCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.RoutesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	rtName := d.Get("route_table_name").(string)
	resGroup := d.Get("resource_group_name").(string)

	addressPrefix := d.Get("address_prefix").(string)
	nextHopType := d.Get("next_hop_type").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, rtName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Route %q (Route Table %q / Resource Group %q): %+v", name, rtName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_route", *existing.ID)
		}
	}

	locks.ByName(rtName, routeTableResourceName)
	defer locks.UnlockByName(rtName, routeTableResourceName)

	route := network.Route{
		Name: &name,
		RoutePropertiesFormat: &network.RoutePropertiesFormat{
			AddressPrefix: &addressPrefix,
			NextHopType:   network.RouteNextHopType(nextHopType),
		},
	}

	if v, ok := d.GetOk("next_hop_in_ip_address"); ok {
		route.RoutePropertiesFormat.NextHopIPAddress = utils.String(v.(string))
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, rtName, name, route)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Route %q (Route Table %q / Resource Group %q): %+v", name, rtName, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion for Route %q (Route Table %q / Resource Group %q): %+v", name, rtName, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, rtName, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Route %q/%q (resource group %q) ID", rtName, name, resGroup)
	}
	d.SetId(*read.ID)

	return resourceRouteRead(d, meta)
}

func resourceRouteRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.RoutesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	rtName := id.Path["routeTables"]
	routeName := id.Path["routes"]

	resp, err := client.Get(ctx, resGroup, rtName, routeName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Route %q: %+v", routeName, err)
	}

	d.Set("name", routeName)
	d.Set("resource_group_name", resGroup)
	d.Set("route_table_name", rtName)

	if props := resp.RoutePropertiesFormat; props != nil {
		d.Set("address_prefix", props.AddressPrefix)
		d.Set("next_hop_type", string(props.NextHopType))
		d.Set("next_hop_in_ip_address", props.NextHopIPAddress)
	}

	return nil
}

func resourceRouteDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.RoutesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	rtName := id.Path["routeTables"]
	routeName := id.Path["routes"]

	locks.ByName(rtName, routeTableResourceName)
	defer locks.UnlockByName(rtName, routeTableResourceName)

	future, err := client.Delete(ctx, resGroup, rtName, routeName)
	if err != nil {
		return fmt.Errorf("Error deleting Route %q (Route Table %q / Resource Group %q): %+v", routeName, rtName, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Route %q (Route Table %q / Resource Group %q): %+v", routeName, rtName, resGroup, err)
	}

	return nil
}
