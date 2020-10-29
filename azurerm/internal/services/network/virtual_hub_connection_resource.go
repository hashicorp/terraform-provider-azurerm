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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmVirtualHubConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVirtualHubConnectionCreateOrUpdate,
		Read:   resourceArmVirtualHubConnectionRead,
		Update: resourceArmVirtualHubConnectionCreateOrUpdate,
		Delete: resourceArmVirtualHubConnectionDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateVirtualHubConnectionName,
			},

			"virtual_hub_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"remote_virtual_network_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			// TODO 3.0: remove this property
			"hub_to_vitual_network_traffic_allowed": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Due to a breaking behavioural change in the Azure API this property is no longer functional and will be removed in version 3.0 of the provider",
			},

			// TODO 3.0: remove this property
			"vitual_network_to_hub_gateways_traffic_allowed": {
				Type:       schema.TypeBool,
				Optional:   true,
				Deprecated: "Due to a breaking behavioural change in the Azure API this property is no longer functional and will be removed in version 3.0 of the provider",
			},

			"internet_security_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
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
											Type: schema.TypeString,
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
		},
	}
}

func resourceArmVirtualHubConnectionCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.HubVirtualNetworkConnectionClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualHubID(d.Get("virtual_hub_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(id.Name, virtualHubResourceName)
	defer locks.UnlockByName(id.Name, virtualHubResourceName)

	name := d.Get("name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Connection %q (Virtual Hub %q / Resource Group %q): %+v", name, id.Name, id.ResourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_virtual_hub_connection", *existing.ID)
		}
	}

	connection := network.HubVirtualNetworkConnection{
		Name: utils.String(name),
		HubVirtualNetworkConnectionProperties: &network.HubVirtualNetworkConnectionProperties{
			RemoteVirtualNetwork: &network.SubResource{
				ID: utils.String(d.Get("remote_virtual_network_id").(string)),
			},
			EnableInternetSecurity: utils.Bool(d.Get("internet_security_enabled").(bool)),
		},
	}

	if v, ok := d.GetOk("routing"); ok {
		connection.HubVirtualNetworkConnectionProperties.RoutingConfiguration = expandArmVirtualHubConnectionRouting(v.([]interface{}))
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, name, connection)
	if err != nil {
		return fmt.Errorf("creating Connection %q (Virtual Hub %q / Resource Group %q): %+v", name, id.Name, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Connection %q (Virtual Hub %q / Resource Group %q): %+v", name, id.Name, id.ResourceGroup, err)
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving Connection %q (Virtual Hub %q / Resource Group %q): %+v", name, id.Name, id.ResourceGroup, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("cannot read Connection %q (Virtual Hub %q / Resource Group %q) ID", name, id.Name, id.ResourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmVirtualHubConnectionRead(d, meta)
}

func resourceArmVirtualHubConnectionRead(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Network.HubVirtualNetworkConnectionClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualHubConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.VirtualHubName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Connection %q (Virtual Hub %q / Resource Group %q) does not exist - removing from state", id.Name, id.VirtualHubName, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Connection %q (Virtual Hub %q / Resource Group %q): %+v", id.Name, id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("virtual_hub_id", parse.NewVirtualHubID(id.ResourceGroup, id.VirtualHubName).ID(subscriptionId))

	if props := resp.HubVirtualNetworkConnectionProperties; props != nil {
		// The following two attributes are deprecated by API (which will always return `true`).
		// Hence, we explicitly set them to `false` (as false is the default value when users omit that property).
		// TODO: 3.0: Remove below lines.
		d.Set("hub_to_vitual_network_traffic_allowed", false)
		d.Set("vitual_network_to_hub_gateways_traffic_allowed", false)

		d.Set("internet_security_enabled", props.EnableInternetSecurity)
		remoteVirtualNetworkId := ""
		if props.RemoteVirtualNetwork != nil && props.RemoteVirtualNetwork.ID != nil {
			remoteVirtualNetworkId = *props.RemoteVirtualNetwork.ID
		}
		d.Set("remote_virtual_network_id", remoteVirtualNetworkId)

		if err := d.Set("routing", flattenArmVirtualHubConnectionRouting(props.RoutingConfiguration)); err != nil {
			return fmt.Errorf("setting `routing`: %+v", err)
		}
	}

	return nil
}

func resourceArmVirtualHubConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.HubVirtualNetworkConnectionClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualHubConnectionID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(id.VirtualHubName, virtualHubResourceName)

	future, err := client.Delete(ctx, id.ResourceGroup, id.VirtualHubName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Connection %q (Virtual Hub %q / Resource Group %q): %+v", id.Name, id.VirtualHubName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deleting Connection %q (Virtual Hub %q / Resource Group %q): %+v", id.Name, id.VirtualHubName, id.ResourceGroup, err)
	}

	return nil
}

func expandArmVirtualHubConnectionRouting(input []interface{}) *network.RoutingConfiguration {
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
		result.VnetRoutes = expandArmVirtualHubConnectionVnetStaticRoute(vnetStaticRoute)
	}

	if propagatedRouteTable := v["propagated_route_table"].([]interface{}); len(propagatedRouteTable) != 0 {
		result.PropagatedRouteTables = expandArmVirtualHubConnectionPropagatedRouteTable(propagatedRouteTable)
	}

	return &result
}

func expandArmVirtualHubConnectionPropagatedRouteTable(input []interface{}) *network.PropagatedRouteTable {
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

func expandArmVirtualHubConnectionVnetStaticRoute(input []interface{}) *network.VnetRoute {
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

func flattenArmVirtualHubConnectionRouting(input *network.RoutingConfiguration) []interface{} {
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
			"propagated_route_table":    flattenArmVirtualHubConnectionPropagatedRouteTable(input.PropagatedRouteTables),
			"static_vnet_route":         flattenArmVirtualHubConnectionVnetStaticRoute(input.VnetRoutes),
		},
	}
}

func flattenArmVirtualHubConnectionPropagatedRouteTable(input *network.PropagatedRouteTable) []interface{} {
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

func flattenArmVirtualHubConnectionVnetStaticRoute(input *network.VnetRoute) []interface{} {
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
