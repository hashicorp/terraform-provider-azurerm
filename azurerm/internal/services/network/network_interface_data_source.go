package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceNetworkInterface() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworkInterfaceRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"network_security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"mac_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"virtual_machine_id": {
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

						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"private_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"private_ip_address_version": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"private_ip_address_allocation": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"public_ip_address_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"application_gateway_backend_address_pools_ids": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},

						"load_balancer_backend_address_pools_ids": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},

						"load_balancer_inbound_nat_rules_ids": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},

						"application_security_group_ids": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},

						"primary": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"dns_servers": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"internal_dns_name_label": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"applied_dns_servers": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"enable_accelerated_networking": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"enable_ip_forwarding": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"private_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"private_ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceNetworkInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.InterfacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	resp, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Network Interface %q (Resource Group %q) was not found", name, resGroup)
		}
		return fmt.Errorf("Error making Read request on Azure Network Interface %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.SetId(*resp.ID)

	iface := *resp.InterfacePropertiesFormat

	d.Set("mac_address", iface.MacAddress)

	if iface.IPConfigurations != nil && len(*iface.IPConfigurations) > 0 {
		configs := *iface.IPConfigurations

		d.Set("private_ip_address", configs[0].InterfaceIPConfigurationPropertiesFormat.PrivateIPAddress)

		addresses := make([]interface{}, 0)
		for _, config := range configs {
			if config.InterfaceIPConfigurationPropertiesFormat != nil {
				addresses = append(addresses, config.InterfaceIPConfigurationPropertiesFormat.PrivateIPAddress)
			}
		}

		if err := d.Set("private_ip_addresses", addresses); err != nil {
			return err
		}
	}

	if iface.IPConfigurations != nil {
		d.Set("ip_configuration", flattenNetworkInterfaceIPConfigurations(iface.IPConfigurations))
	}

	if iface.VirtualMachine != nil {
		d.Set("virtual_machine_id", iface.VirtualMachine.ID)
	} else {
		d.Set("virtual_machine_id", "")
	}

	var appliedDNSServers []string
	var dnsServers []string
	if dnsSettings := iface.DNSSettings; dnsSettings != nil {
		if s := dnsSettings.AppliedDNSServers; s != nil {
			appliedDNSServers = *s
		}

		if s := dnsSettings.DNSServers; s != nil {
			dnsServers = *s
		}

		d.Set("internal_dns_name_label", dnsSettings.InternalDNSNameLabel)
	}

	if iface.NetworkSecurityGroup != nil {
		d.Set("network_security_group_id", resp.NetworkSecurityGroup.ID)
	} else {
		d.Set("network_security_group_id", "")
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	d.Set("applied_dns_servers", appliedDNSServers)
	d.Set("dns_servers", dnsServers)
	d.Set("enable_ip_forwarding", resp.EnableIPForwarding)
	d.Set("enable_accelerated_networking", resp.EnableAcceleratedNetworking)

	return tags.FlattenAndSet(d, resp.Tags)
}
