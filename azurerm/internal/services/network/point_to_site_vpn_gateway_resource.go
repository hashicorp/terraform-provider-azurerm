package network

import (
	"fmt"
	"log"
	"time"

	validate2 "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPointToSiteVPNGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPointToSiteVPNGatewayCreateUpdate,
		Read:   resourceArmPointToSiteVPNGatewayRead,
		Update: resourceArmPointToSiteVPNGatewayCreateUpdate,
		Delete: resourceArmPointToSiteVPNGatewayDelete,
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"virtual_hub_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate2.ValidateVirtualHubID,
			},

			"vpn_server_configuration_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateVpnServerConfigurationID,
			},

			"connection_configuration": {
				Type:     schema.TypeList,
				Required: true,
				// Code="P2SVpnGatewayCanHaveOnlyOneP2SConnectionConfiguration"
				// Message="Currently, P2SVpnGateway [ID] can have only one P2SConnectionConfiguration. Specified number of P2SConnectionConfiguration are 2.
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"vpn_client_address_pool": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address_prefixes": {
										Type:     schema.TypeSet,
										Required: true,
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validate.CIDR,
										},
									},
								},
							},
						},

						"route": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"associated_route_table_id": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validate2.VirtualHubRouteTableID,
									},

									"propagated_route_table": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ids": {
													Type:     schema.TypeList,
													Required: true,
													Elem: &schema.Schema{
														Type:         schema.TypeString,
														ValidateFunc: validate2.VirtualHubRouteTableID,
													},
												},

												"labels": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Schema{
														Type:         schema.TypeString,
														ValidateFunc: validation.StringIsNotEmpty,
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

			"scale_unit": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},

			"dns_servers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsIPv4Address,
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmPointToSiteVPNGatewayCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PointToSiteVpnGatewaysClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Point-to-Site VPN Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_point_to_site_vpn_gateway", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	scaleUnit := d.Get("scale_unit").(int)
	virtualHubId := d.Get("virtual_hub_id").(string)
	vpnServerConfigurationId := d.Get("vpn_server_configuration_id").(string)
	t := d.Get("tags").(map[string]interface{})

	connectionConfigurationsRaw := d.Get("connection_configuration").([]interface{})
	connectionConfigurations := expandPointToSiteVPNGatewayConnectionConfiguration(connectionConfigurationsRaw)

	parameters := network.P2SVpnGateway{
		Location: utils.String(location),
		P2SVpnGatewayProperties: &network.P2SVpnGatewayProperties{
			P2SConnectionConfigurations: connectionConfigurations,
			VpnServerConfiguration: &network.SubResource{
				ID: utils.String(vpnServerConfigurationId),
			},
			VirtualHub: &network.SubResource{
				ID: utils.String(virtualHubId),
			},
			VpnGatewayScaleUnit: utils.Int32(int32(scaleUnit)),
		},
		Tags: tags.Expand(t),
	}
	customDNSServers := utils.ExpandStringSlice(d.Get("dns_servers").([]interface{}))
	if len(*customDNSServers) != 0 {
		parameters.P2SVpnGatewayProperties.CustomDNSServers = customDNSServers
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating/updating Point-to-Site VPN Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation/update of Point-to-Site VPN Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Point-to-Site VPN Gateway %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceArmPointToSiteVPNGatewayRead(d, meta)
}

func resourceArmPointToSiteVPNGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PointToSiteVpnGatewaysClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := ParsePointToSiteVPNGatewayID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Point-to-Site VPN Gateway %q was not found in Resource Group %q - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Point-to-Site VPN Gateway %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.P2SVpnGatewayProperties; props != nil {
		d.Set("dns_servers", utils.FlattenStringSlice(props.CustomDNSServers))
		flattenedConfigurations := flattenPointToSiteVPNGatewayConnectionConfiguration(props.P2SConnectionConfigurations)
		if err := d.Set("connection_configuration", flattenedConfigurations); err != nil {
			return fmt.Errorf("Error setting `connection_configuration`: %+v", err)
		}

		scaleUnit := 0
		if props.VpnGatewayScaleUnit != nil {
			scaleUnit = int(*props.VpnGatewayScaleUnit)
		}
		d.Set("scale_unit", scaleUnit)

		virtualHubId := ""
		if props.VirtualHub != nil && props.VirtualHub.ID != nil {
			virtualHubId = *props.VirtualHub.ID
		}
		d.Set("virtual_hub_id", virtualHubId)

		vpnServerConfigurationId := ""
		if props.VpnServerConfiguration != nil && props.VpnServerConfiguration.ID != nil {
			vpnServerConfigurationId = *props.VpnServerConfiguration.ID
		}
		d.Set("vpn_server_configuration_id", vpnServerConfigurationId)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmPointToSiteVPNGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PointToSiteVpnGatewaysClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := ParsePointToSiteVPNGatewayID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error deleting Point-to-Site VPN Gateway %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Point-to-Site VPN Gateway %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func expandPointToSiteVPNGatewayConnectionConfiguration(input []interface{}) *[]network.P2SConnectionConfiguration {
	configurations := make([]network.P2SConnectionConfiguration, 0)

	for _, v := range input {
		raw := v.(map[string]interface{})

		addressPrefixes := make([]string, 0)
		name := raw["name"].(string)

		clientAddressPoolsRaw := raw["vpn_client_address_pool"].([]interface{})
		for _, clientV := range clientAddressPoolsRaw {
			clientRaw := clientV.(map[string]interface{})

			addressPrefixesRaw := clientRaw["address_prefixes"].(*schema.Set).List()
			for _, prefix := range addressPrefixesRaw {
				addressPrefixes = append(addressPrefixes, prefix.(string))
			}
		}

		configurations = append(configurations, network.P2SConnectionConfiguration{
			Name: utils.String(name),
			P2SConnectionConfigurationProperties: &network.P2SConnectionConfigurationProperties{
				VpnClientAddressPool: &network.AddressSpace{
					AddressPrefixes: &addressPrefixes,
				},
				RoutingConfiguration: expandPointToSiteVPNGatewayConnectionRouteConfiguration(raw["route"].([]interface{})),
			},
		})
	}

	return &configurations
}

func expandPointToSiteVPNGatewayConnectionRouteConfiguration(input []interface{}) *network.RoutingConfiguration {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &network.RoutingConfiguration{
		AssociatedRouteTable: &network.SubResource{
			ID: utils.String(v["associated_route_table_id"].(string)),
		},
		PropagatedRouteTables: expandPointToSiteVPNGatewayConnectionRouteConfigurationPropagatedRouteTable(v["propagated_route_table"].([]interface{})),
	}
}

func expandPointToSiteVPNGatewayConnectionRouteConfigurationPropagatedRouteTable(input []interface{}) *network.PropagatedRouteTable {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	idRaws := utils.ExpandStringSlice(v["ids"].([]interface{}))
	ids := make([]network.SubResource, len(*idRaws))
	for i, item := range *idRaws {
		ids[i] = network.SubResource{
			ID: utils.String(item),
		}
	}
	return &network.PropagatedRouteTable{
		Labels: utils.ExpandStringSlice(v["labels"].(*schema.Set).List()),
		Ids:    &ids,
	}
}

func flattenPointToSiteVPNGatewayConnectionConfiguration(input *[]network.P2SConnectionConfiguration) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, v := range *input {
		name := ""
		if v.Name != nil {
			name = *v.Name
		}

		addressPrefixes := make([]interface{}, 0)
		if props := v.P2SConnectionConfigurationProperties; props != nil {
			if props.VpnClientAddressPool == nil {
				continue
			}

			if props.VpnClientAddressPool.AddressPrefixes != nil {
				for _, prefix := range *props.VpnClientAddressPool.AddressPrefixes {
					addressPrefixes = append(addressPrefixes, prefix)
				}
			}
		}

		output = append(output, map[string]interface{}{
			"name": name,
			"vpn_client_address_pool": []interface{}{
				map[string]interface{}{
					"address_prefixes": addressPrefixes,
				},
			},
			"route": flattenPointToSiteVPNGatewayConnectionRouteConfiguration(v.RoutingConfiguration),
		})
	}

	return output
}

func flattenPointToSiteVPNGatewayConnectionRouteConfiguration(input *network.RoutingConfiguration) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	var associatedRouteTableId string
	if input.AssociatedRouteTable != nil && input.AssociatedRouteTable.ID != nil {
		associatedRouteTableId = *input.AssociatedRouteTable.ID
	}
	return []interface{}{
		map[string]interface{}{
			"associated_route_table_id": associatedRouteTableId,
			"propagated_route_table":    flattenPointToSiteVPNGatewayConnectionRouteConfigurationPropagatedRouteTable(input.PropagatedRouteTables),
		},
	}
}

func flattenPointToSiteVPNGatewayConnectionRouteConfigurationPropagatedRouteTable(input *network.PropagatedRouteTable) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	ids := make([]string, 0)
	if input.Ids != nil {
		for _, item := range *input.Ids {
			if item.ID != nil {
				ids = append(ids, *item.ID)
			}
		}
	}
	return []interface{}{
		map[string]interface{}{
			"ids":    ids,
			"labels": utils.FlattenStringSlice(input.Labels),
		},
	}
}
