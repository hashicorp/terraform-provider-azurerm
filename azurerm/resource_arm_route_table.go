package azurerm

import (
	"fmt"
	"log"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var routeTableResourceName = "azurerm_route_table"

func resourceArmRouteTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmRouteTableCreateUpdate,
		Read:   resourceArmRouteTableRead,
		Update: resourceArmRouteTableCreateUpdate,
		Delete: resourceArmRouteTableDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"route": {
				Type:       schema.TypeList,
				ConfigMode: schema.SchemaConfigModeAttr,
				Optional:   true,
				Computed:   true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},

						"address_prefix": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
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
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},
			},

			"disable_bgp_route_propagation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"subnets": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmRouteTableCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.RouteTablesClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM Route Table creation.")

	name := d.Get("name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	resGroup := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})

	if requireResourcesToBeImported && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Route Table %q (Resource Group %q): %+v", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_route_table", *existing.ID)
		}
	}

	routeSet := network.RouteTable{
		Name:     &name,
		Location: &location,
		RouteTablePropertiesFormat: &network.RouteTablePropertiesFormat{
			Routes:                     expandRouteTableRoutes(d),
			DisableBgpRoutePropagation: utils.Bool(d.Get("disable_bgp_route_propagation").(bool)),
		},
		Tags: expandTags(tags),
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, routeSet)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Route Table %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Route Table %q (Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Route Table %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmRouteTableRead(d, meta)
}

func resourceArmRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.RouteTablesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["routeTables"]

	resp, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Route Table %q: %+v", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.RouteTablePropertiesFormat; props != nil {
		d.Set("disable_bgp_route_propagation", props.DisableBgpRoutePropagation)
		if err := d.Set("route", flattenRouteTableRoutes(props.Routes)); err != nil {
			return err
		}

		if err := d.Set("subnets", flattenRouteTableSubnets(props.Subnets)); err != nil {
			return err
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmRouteTableDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.RouteTablesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["routeTables"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error deleting Route Table %q (Resource Group %q): %+v", name, resGroup, err)
		}
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Route Table %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}

func expandRouteTableRoutes(d *schema.ResourceData) *[]network.Route {
	configs := d.Get("route").([]interface{})
	routes := make([]network.Route, 0, len(configs))

	for _, configRaw := range configs {
		data := configRaw.(map[string]interface{})

		route := network.Route{
			Name: utils.String(data["name"].(string)),
			RoutePropertiesFormat: &network.RoutePropertiesFormat{
				AddressPrefix: utils.String(data["address_prefix"].(string)),
				NextHopType:   network.RouteNextHopType(data["next_hop_type"].(string)),
			},
		}

		if v := data["next_hop_in_ip_address"].(string); v != "" {
			route.RoutePropertiesFormat.NextHopIPAddress = &v
		}

		routes = append(routes, route)
	}

	return &routes
}

func flattenRouteTableRoutes(input *[]network.Route) []interface{} {
	results := make([]interface{}, 0)

	if routes := input; routes != nil {
		for _, route := range *routes {
			r := make(map[string]interface{})

			r["name"] = *route.Name

			if props := route.RoutePropertiesFormat; props != nil {
				r["address_prefix"] = *props.AddressPrefix
				r["next_hop_type"] = string(props.NextHopType)
				if ip := props.NextHopIPAddress; ip != nil {
					r["next_hop_in_ip_address"] = *ip
				}
			}

			results = append(results, r)
		}
	}

	return results
}

func flattenRouteTableSubnets(subnets *[]network.Subnet) []string {
	output := make([]string, 0)

	if subnets != nil {
		for _, subnet := range *subnets {
			output = append(output, *subnet.ID)
		}
	}

	return output
}
