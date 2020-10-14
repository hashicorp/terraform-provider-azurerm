package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

			"routing_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"propagated_route_table": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"labels": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},

						"vnet_route": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"static_route": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Optional: true,
												},

												"address_prefixes": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},

												"next_hop_ip_address": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
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
			RoutingConfiguration:   expandArmVirtualHubConnectionRoutingConfiguration(d.Get("routing_configuration").([]interface{})),
		},
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

		if err := d.Set("routing_configuration", flattenArmVirtualHubConnectionRoutingConfiguration(props.RoutingConfiguration)); err != nil {
			return fmt.Errorf("setting `routing_configuration`: %+v", err)
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

func expandArmVirtualHubConnectionRoutingConfiguration(input []interface{}) *network.RoutingConfiguration {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &network.RoutingConfiguration{
		PropagatedRouteTables: expandArmVirtualHubConnectionPropagatedRouteTable(v["propagated_route_table"].([]interface{})),
		VnetRoutes:            expandArmVirtualHubConnectionVnetRoute(v["vnet_route"].([]interface{})),
	}
}

func expandArmVirtualHubConnectionPropagatedRouteTable(input []interface{}) *network.PropagatedRouteTable {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &network.PropagatedRouteTable{
		Labels: utils.ExpandStringSlice(v["labels"].(*schema.Set).List()),
	}
}

func expandArmVirtualHubConnectionVnetRoute(input []interface{}) *network.VnetRoute {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})

	return &network.VnetRoute{
		StaticRoutes: expandArmVirtualHubConnectionStaticRoutes(v["static_route"].(*schema.Set).List()),
	}
}

func expandArmVirtualHubConnectionStaticRoutes(input []interface{}) *[]network.StaticRoute {
	results := make([]network.StaticRoute, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		result := network.StaticRoute{
			Name:             utils.String(v["name"].(string)),
			AddressPrefixes:  utils.ExpandStringSlice(v["address_prefixes"].(*schema.Set).List()),
			NextHopIPAddress: utils.String(v["next_hop_ip_address"].(string)),
		}

		results = append(results, result)
	}

	return &results
}

func flattenArmVirtualHubConnectionRoutingConfiguration(input *network.RoutingConfiguration) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"propagated_route_table": flattenArmVirtualHubConnectionPropagatedRouteTable(input.PropagatedRouteTables),
			"vnet_route":             flattenArmVirtualHubConnectionVnetRoute(input.VnetRoutes),
		},
	}
}

func flattenArmVirtualHubConnectionPropagatedRouteTable(input *network.PropagatedRouteTable) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"labels": utils.FlattenStringSlice(input.Labels),
		},
	}
}

func flattenArmVirtualHubConnectionVnetRoute(input *network.VnetRoute) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"static_route": flattenArmVirtualHubConnectionStaticRoutes(input.StaticRoutes),
		},
	}
}

func flattenArmVirtualHubConnectionStaticRoutes(input *[]network.StaticRoute) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var name string
		if item.Name != nil {
			name = *item.Name
		}

		var nextHopIpAddress string
		if item.NextHopIPAddress != nil {
			nextHopIpAddress = *item.NextHopIPAddress
		}

		v := map[string]interface{}{
			"name":                name,
			"address_prefixes":    utils.FlattenStringSlice(item.AddressPrefixes),
			"next_hop_ip_address": nextHopIpAddress,
		}

		results = append(results, v)
	}

	return results
}
