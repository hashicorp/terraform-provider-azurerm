package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmExpressRouteConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmExpressRouteConnectionCreateUpdate,
		Read:   resourceArmExpressRouteConnectionRead,
		Update: resourceArmExpressRouteConnectionCreateUpdate,
		Delete: resourceArmExpressRouteConnectionDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ExpressRouteConnectionID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"express_route_gateway_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ExpressRouteGatewayID,
			},

			"express_route_circuit_peering_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ExpressRouteCircuitPeeringID,
			},

			"authorization_key": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"enable_internet_security": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"routing": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"associated_route_table_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: azure.ValidateResourceID,
						},

						"propagated_route_table": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"labels": {
										Type:     schema.TypeSet,
										Optional: true,
										Computed: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringIsNotEmpty,
										},
									},

									"route_table_ids": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: azure.ValidateResourceID,
										},
									},
								},
							},
						},

						"static_vnet_route": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"address_prefixes": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.IsCIDR,
										},
									},

									"next_hop_ip_address": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.IsIPv4Address,
									},
								},
							},
						},
					},
				},
			},

			"routing_weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 32000),
			},
		},
	}
}
func resourceArmExpressRouteConnectionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteConnectionsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)

	id, err := parse.ExpressRouteGatewayID(d.Get("express_route_gateway_id").(string))
	if err != nil {
		return err
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for present of existing ExpressRouteConnection %q (Resource Group %q / expressRouteGatewayName %q): %+v", name, id.ResourceGroup, id.Name, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_network_express_route_connection", *existing.ID)
		}
	}

	putExpressRouteConnectionParameters := network.ExpressRouteConnection{
		Name: utils.String(d.Get("name").(string)),
		ExpressRouteConnectionProperties: &network.ExpressRouteConnectionProperties{
			ExpressRouteCircuitPeering: &network.ExpressRouteCircuitPeeringID{
				ID: utils.String(d.Get("express_route_circuit_peering_id").(string)),
			},
		},
	}

	if v, ok := d.GetOk("authorization_key"); ok {
		putExpressRouteConnectionParameters.ExpressRouteConnectionProperties.AuthorizationKey = utils.String(v.(string))
	}

	if v, ok := d.GetOk("enable_internet_security"); ok {
		putExpressRouteConnectionParameters.ExpressRouteConnectionProperties.EnableInternetSecurity = utils.Bool(v.(bool))
	}

	if v, ok := d.GetOk("routing_weight"); ok {
		putExpressRouteConnectionParameters.ExpressRouteConnectionProperties.RoutingWeight = utils.Int32(int32(v.(int)))
	}

	if v, ok := d.GetOk("routing"); ok {
		putExpressRouteConnectionParameters.ExpressRouteConnectionProperties.RoutingConfiguration = expandArmExpressRouteConnectionRouting(v.([]interface{}))
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, name, putExpressRouteConnectionParameters)
	if err != nil {
		return fmt.Errorf("creating/updating ExpressRouteConnection %q (Resource Group %q / expressRouteGatewayName %q): %+v", name, id.ResourceGroup, id.Name, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creating/updating future for ExpressRouteConnection %q (Resource Group %q / expressRouteGatewayName %q): %+v", name, id.ResourceGroup, id.Name, err)
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving ExpressRouteConnection %q (Resource Group %q / expressRouteGatewayName %q): %+v", name, id.ResourceGroup, id.Name, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for ExpressRouteConnection %q (Resource Group %q / expressRouteGatewayName %q) ID", name, id.ResourceGroup, id.Name)
	}

	d.SetId(*resp.ID)

	return resourceArmExpressRouteConnectionRead(d, meta)
}

func resourceArmExpressRouteConnectionRead(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Network.ExpressRouteConnectionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ExpressRouteConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ExpressRouteGatewayName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] network %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Network ExpressRouteConnection %q (Resource Group %q / expressRouteGatewayName %q): %+v", id.Name, id.ResourceGroup, id.ExpressRouteGatewayName, err)
	}

	d.Set("name", id.Name)
	d.Set("expresss_route_gateway_id", parse.NewExpressRouteGatewayID(id.ResourceGroup, id.Name).ID(subscriptionId))

	if props := resp.ExpressRouteConnectionProperties; props != nil {
		if v := props.ExpressRouteCircuitPeering; v != nil {
			d.Set("express_route_circuit_peering_id", v.ID)
		}

		if v := props.AuthorizationKey; v != nil {
			d.Set("authorization_key", v)
		}

		if v := props.EnableInternetSecurity; v != nil {
			d.Set("enable_internet_security", v)
		}

		if v := props.RoutingWeight; v != nil {
			d.Set("routing_weight", v)
		}

		if err := d.Set("routing", flattenArmExpressRouteConnectionRouting(props.RoutingConfiguration)); err != nil {
			return fmt.Errorf("setting `routing`: %+v", err)
		}
	}

	return nil
}

func resourceArmExpressRouteConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ExpressRouteConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ExpressRouteConnectionID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ExpressRouteGatewayName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting ExpressRouteConnection %q (Resource Group %q / expressRouteGatewayName %q): %+v", id.Name, id.ResourceGroup, id.ExpressRouteGatewayName, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on deleting future for ExpressRouteConnection %q (Resource Group %q / expressRouteGatewayName %q): %+v", id.Name, id.ResourceGroup, id.ExpressRouteGatewayName, err)
	}

	return nil
}

func expandArmExpressRouteConnectionRouting(input []interface{}) *network.RoutingConfiguration {
	if len(input) == 0 {
		return &network.RoutingConfiguration{}
	}

	v := input[0].(map[string]interface{})
	result := network.RoutingConfiguration{}

	if associatedRouteTableId := v["associated_route_table_id"].(string); associatedRouteTableId != "" {
		result.AssociatedRouteTable = &network.SubResource{
			ID: utils.String(associatedRouteTableId),
		}
	}

	if vnetStaticRoute := v["static_vnet_route"].([]interface{}); len(vnetStaticRoute) != 0 {
		result.VnetRoutes = expandArmExpressRouteConnectionStaticVnetRoute(vnetStaticRoute)
	}

	if propagatedRouteTable := v["propagated_route_table"].([]interface{}); len(propagatedRouteTable) != 0 {
		result.PropagatedRouteTables = expandArmExpressRouteConnectionPropagatedRouteTable(propagatedRouteTable)
	}

	return &result
}

func expandArmExpressRouteConnectionPropagatedRouteTable(input []interface{}) *network.PropagatedRouteTable {
	if len(input) == 0 {
		return &network.PropagatedRouteTable{}
	}

	v := input[0].(map[string]interface{})

	result := network.PropagatedRouteTable{}

	if labels := v["labels"].(*schema.Set).List(); len(labels) != 0 {
		result.Labels = utils.ExpandStringSlice(labels)
	}

	if routeTableIds := v["route_table_ids"].([]interface{}); len(routeTableIds) != 0 {
		result.Ids = expandIDsToSubResources(routeTableIds)
	}

	return &result
}

func expandArmExpressRouteConnectionStaticVnetRoute(input []interface{}) *network.VnetRoute {
	if len(input) == 0 {
		return &network.VnetRoute{}
	}

	results := make([]network.StaticRoute, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		result := network.StaticRoute{}

		if name := v["name"].(string); name != "" {
			result.Name = utils.String(name)
		}

		if addressPrefixes := v["address_prefixes"].(*schema.Set).List(); len(addressPrefixes) != 0 {
			result.AddressPrefixes = utils.ExpandStringSlice(addressPrefixes)
		}

		if nextHopIPAddress := v["next_hop_ip_address"].(string); nextHopIPAddress != "" {
			result.NextHopIPAddress = utils.String(nextHopIPAddress)
		}

		results = append(results, result)
	}

	return &network.VnetRoute{
		StaticRoutes: &results,
	}
}

func expandIDsToSubResources(input []interface{}) *[]network.SubResource {
	ids := make([]network.SubResource, 0)

	for _, v := range input {
		ids = append(ids, network.SubResource{
			ID: utils.String(v.(string)),
		})
	}

	return &ids
}

func flattenArmExpressRouteConnectionRouting(input *network.RoutingConfiguration) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	associatedRouteTableId := ""
	if input.AssociatedRouteTable != nil && input.AssociatedRouteTable.ID != nil {
		associatedRouteTableId = *input.AssociatedRouteTable.ID
	}

	return []interface{}{
		map[string]interface{}{
			"associated_route_table_id": associatedRouteTableId,
			"propagated_route_table":    flattenArmExpressRouteConnectionPropagatedRouteTable(input.PropagatedRouteTables),
			"static_vnet_route":         flattenArmExpressRouteConnectionStaticVnetRoute(input.VnetRoutes),
		},
	}
}

func flattenArmExpressRouteConnectionPropagatedRouteTable(input *network.PropagatedRouteTable) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	labels := make([]interface{}, 0)
	if input.Labels != nil {
		labels = utils.FlattenStringSlice(input.Labels)
	}

	routeTableIds := make([]interface{}, 0)
	if input.Ids != nil {
		routeTableIds = flattenSubResourcesToIDs(input.Ids)
	}

	return []interface{}{
		map[string]interface{}{
			"labels":          labels,
			"route_table_ids": routeTableIds,
		},
	}
}

func flattenArmExpressRouteConnectionStaticVnetRoute(input *network.VnetRoute) []interface{} {
	results := make([]interface{}, 0)
	if input == nil || input.StaticRoutes == nil {
		return results
	}

	for _, item := range *input.StaticRoutes {
		var name string
		if item.Name != nil {
			name = *item.Name
		}

		var nextHopIpAddress string
		if item.NextHopIPAddress != nil {
			nextHopIpAddress = *item.NextHopIPAddress
		}

		addressPrefixes := make([]interface{}, 0)
		if item.AddressPrefixes != nil {
			addressPrefixes = utils.FlattenStringSlice(item.AddressPrefixes)
		}

		v := map[string]interface{}{
			"name":                name,
			"address_prefixes":    addressPrefixes,
			"next_hop_ip_address": nextHopIpAddress,
		}

		results = append(results, v)
	}

	return results
}

func flattenSubResourcesToIDs(input *[]network.SubResource) []interface{} {
	ids := make([]interface{}, 0)
	if input == nil {
		return ids
	}

	for _, v := range *input {
		if v.ID == nil {
			continue
		}

		ids = append(ids, *v.ID)
	}

	return ids
}
