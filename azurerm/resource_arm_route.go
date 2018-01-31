package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmRouteCreateUpdate,
		Read:   resourceArmRouteRead,
		Update: resourceArmRouteCreateUpdate,
		Delete: resourceArmRouteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"route_table_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"address_prefix": {
				Type:     schema.TypeString,
				Required: true,
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
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"next_hop_in_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceArmRouteCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).routesClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	rtName := d.Get("route_table_name").(string)
	resGroup := d.Get("resource_group_name").(string)

	addressPrefix := d.Get("address_prefix").(string)
	nextHopType := d.Get("next_hop_type").(string)

	azureRMLockByName(rtName, routeTableResourceName)
	defer azureRMUnlockByName(rtName, routeTableResourceName)

	properties := network.RoutePropertiesFormat{
		AddressPrefix: &addressPrefix,
		NextHopType:   network.RouteNextHopType(nextHopType),
	}

	if v, ok := d.GetOk("next_hop_in_ip_address"); ok {
		nextHopInIpAddress := v.(string)
		properties.NextHopIPAddress = &nextHopInIpAddress
	}

	route := network.Route{
		Name: &name,
		RoutePropertiesFormat: &properties,
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, rtName, name, route)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Route %q (Route Table %q / Resource Group %q): %+v", name, rtName, resGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
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

	return resourceArmRouteRead(d, meta)
}

func resourceArmRouteRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).routesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
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

		if ip := props.NextHopIPAddress; ip != nil {
			d.Set("next_hop_in_ip_address", props.NextHopIPAddress)
		}
	}

	return nil
}

func resourceArmRouteDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).routesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	rtName := id.Path["routeTables"]
	routeName := id.Path["routes"]

	azureRMLockByName(rtName, routeTableResourceName)
	defer azureRMUnlockByName(rtName, routeTableResourceName)

	future, err := client.Delete(ctx, resGroup, rtName, routeName)
	if err != nil {
		return fmt.Errorf("Error deleting Route %q (Route Table %q / Resource Group %q): %+v", routeName, rtName, resGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for deletion of Route %q (Route Table %q / Resource Group %q): %+v", routeName, rtName, resGroup, err)
	}

	return nil
}
