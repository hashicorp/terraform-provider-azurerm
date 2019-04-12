package azurerm

import (
	"fmt"

	"bytes"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmVirtualNetworkGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmVirtualNetworkGatewayRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
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
							Set: hashVirtualNetworkGatewayDataSourceRootCert,
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
							Set: hashVirtualNetworkGatewayDataSourceRevokedCert,
						},

						"radius_server_address": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"radius_server_secret": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"vpn_client_protocols": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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

		if err := d.Set("ip_configuration", flattenArmVirtualNetworkGatewayDataSourceIPConfigurations(gw.IPConfigurations)); err != nil {
			return fmt.Errorf("Error setting `ip_configuration`: %+v", err)
		}

		vpnConfigFlat := flattenArmVirtualNetworkGatewayDataSourceVpnClientConfig(gw.VpnClientConfiguration)
		if err := d.Set("vpn_client_configuration", vpnConfigFlat); err != nil {
			return fmt.Errorf("Error setting `vpn_client_configuration`: %+v", err)
		}

		bgpSettingsFlat := flattenArmVirtualNetworkGatewayDataSourceBgpSettings(gw.BgpSettings)
		if err := d.Set("bgp_settings", bgpSettingsFlat); err != nil {
			return fmt.Errorf("Error setting `bgp_settings`: %+v", err)
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func flattenArmVirtualNetworkGatewayDataSourceIPConfigurations(ipConfigs *[]network.VirtualNetworkGatewayIPConfiguration) []interface{} {
	flat := make([]interface{}, 0)

	if ipConfigs != nil {
		for _, cfg := range *ipConfigs {
			props := cfg.VirtualNetworkGatewayIPConfigurationPropertiesFormat
			v := make(map[string]interface{})

			if name := cfg.Name; name != nil {
				v["name"] = *name
			}
			v["private_ip_address_allocation"] = string(props.PrivateIPAllocationMethod)

			if subnet := props.Subnet; subnet != nil {
				if id := subnet.ID; id != nil {
					v["subnet_id"] = *id
				}
			}

			if pip := props.PublicIPAddress; pip != nil {
				if id := pip.ID; id != nil {
					v["public_ip_address_id"] = *id
				}
			}

			flat = append(flat, v)
		}
	}

	return flat
}

func flattenArmVirtualNetworkGatewayDataSourceVpnClientConfig(cfg *network.VpnClientConfiguration) []interface{} {
	if cfg == nil {
		return []interface{}{}
	}

	flat := make(map[string]interface{})

	addressSpace := make([]interface{}, 0)
	if pool := cfg.VpnClientAddressPool; pool != nil {
		if prefixes := pool.AddressPrefixes; prefixes != nil {
			for _, addr := range *prefixes {
				addressSpace = append(addressSpace, addr)
			}
		}
	}
	flat["address_space"] = addressSpace

	rootCerts := make([]interface{}, 0)
	if certs := cfg.VpnClientRootCertificates; certs != nil {
		for _, cert := range *certs {
			v := map[string]interface{}{
				"name":             *cert.Name,
				"public_cert_data": *cert.VpnClientRootCertificatePropertiesFormat.PublicCertData,
			}
			rootCerts = append(rootCerts, v)
		}
	}
	flat["root_certificate"] = schema.NewSet(hashVirtualNetworkGatewayDataSourceRootCert, rootCerts)

	revokedCerts := make([]interface{}, 0)
	if certs := cfg.VpnClientRevokedCertificates; certs != nil {
		for _, cert := range *certs {
			v := map[string]interface{}{
				"name":       *cert.Name,
				"thumbprint": *cert.VpnClientRevokedCertificatePropertiesFormat.Thumbprint,
			}
			revokedCerts = append(revokedCerts, v)
		}
	}
	flat["revoked_certificate"] = schema.NewSet(hashVirtualNetworkGatewayDataSourceRevokedCert, revokedCerts)

	vpnClientProtocols := &schema.Set{F: schema.HashString}
	if vpnProtocols := cfg.VpnClientProtocols; vpnProtocols != nil {
		for _, protocol := range *vpnProtocols {
			vpnClientProtocols.Add(string(protocol))
		}
	}
	flat["vpn_client_protocols"] = vpnClientProtocols

	if v := cfg.RadiusServerAddress; v != nil {
		flat["radius_server_address"] = *v
	}

	if v := cfg.RadiusServerSecret; v != nil {
		flat["radius_server_secret"] = *v
	}

	return []interface{}{flat}
}

func flattenArmVirtualNetworkGatewayDataSourceBgpSettings(settings *network.BgpSettings) []interface{} {
	output := make([]interface{}, 0)

	if settings != nil {
		flat := make(map[string]interface{})

		if asn := settings.Asn; asn != nil {
			flat["asn"] = int(*asn)
		}
		if address := settings.BgpPeeringAddress; address != nil {
			flat["peering_address"] = *address
		}
		if weight := settings.PeerWeight; weight != nil {
			flat["peer_weight"] = int(*weight)
		}

		output = append(output, flat)
	}

	return output
}

func hashVirtualNetworkGatewayDataSourceRootCert(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["public_cert_data"].(string)))

	return hashcode.String(buf.String())
}

func hashVirtualNetworkGatewayDataSourceRevokedCert(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["thumbprint"].(string)))

	return hashcode.String(buf.String())
}
