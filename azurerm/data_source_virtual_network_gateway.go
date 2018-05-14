package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmVirtualNetworkGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmVirtualNetworkGatewayRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"location": locationForDataSourceSchema(),

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vpn_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"enable_bgp": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"active_active": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"sku": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ip_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip_address_allocation": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip_address_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"vpn_client_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address_space": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"root_certificate": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"public_cert_data": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
							Set: hashVirtualNetworkGatewayRootCert,
						},
						"revoked_certificate": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"thumbprint": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
							Set: hashVirtualNetworkGatewayRevokedCert,
						},
					},
				},
			},

			"bgp_settings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"asn": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"peering_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"peer_weight": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"default_local_network_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsForDataSourceSchema(),
		},
	}
}

func dataSourceArmVirtualNetworkGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vnetGatewayClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Virtual Network Gateway %q (Resource Group %q) was not found", name, resGroup)
		}

		return fmt.Errorf("Error making Read request on AzureRM Virtual Network Gateway %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.SetId(*resp.ID)

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if resp.VirtualNetworkGatewayPropertiesFormat != nil {
		gw := *resp.VirtualNetworkGatewayPropertiesFormat

		d.Set("type", string(gw.GatewayType))
		d.Set("enable_bgp", gw.EnableBgp)
		d.Set("active_active", gw.ActiveActive)

		if string(gw.VpnType) != "" {
			d.Set("vpn_type", string(gw.VpnType))
		}

		if gw.GatewayDefaultSite != nil {
			d.Set("default_local_network_gateway_id", gw.GatewayDefaultSite.ID)
		}

		if gw.Sku != nil {
			d.Set("sku", string(gw.Sku.Name))
		}

		d.Set("ip_configuration", flattenArmVirtualNetworkGatewayIPConfigurations(gw.IPConfigurations))

		if gw.VpnClientConfiguration != nil {
			vpnConfigFlat := flattenArmVirtualNetworkGatewayVpnClientConfig(gw.VpnClientConfiguration)
			if err := d.Set("vpn_client_configuration", vpnConfigFlat); err != nil {
				return fmt.Errorf("Error setting `vpn_client_configuration`: %+v", err)
			}
		}

		if gw.BgpSettings != nil {
			bgpSettingsFlat := flattenArmVirtualNetworkGatewayBgpSettings(gw.BgpSettings)
			if err := d.Set("bgp_settings", bgpSettingsFlat); err != nil {
				return fmt.Errorf("Error setting `bgp_settings`: %+v", err)
			}
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}
