package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmNetworkInterface() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmNetworkInterfaceRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

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

			"internal_fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			/**
			 * As of 2018-01-06: AN (aka. SR-IOV) on Azure is GA on Windows and Linux.
			 *
			 * Refer to: https://azure.microsoft.com/en-us/blog/maximize-your-vm-s-performance-with-accelerated-networking-now-generally-available-for-both-windows-and-linux/
			 *
			 * Refer to: https://docs.microsoft.com/en-us/azure/virtual-network/create-vm-accelerated-networking-cli
			 * For details, VM configuration and caveats.
			 */
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

			"tags": tagsForDataSourceSchema(),
		},
	}
}

func dataSourceArmNetworkInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).ifaceClient
	ctx := meta.(*ArmClient).StopContext

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

		if configs[0].InterfaceIPConfigurationPropertiesFormat != nil {
			privateIPAddress := configs[0].InterfaceIPConfigurationPropertiesFormat.PrivateIPAddress
			d.Set("private_ip_address", *privateIPAddress)
		}

		addresses := make([]interface{}, 0)
		for _, config := range configs {
			if config.InterfaceIPConfigurationPropertiesFormat != nil {
				addresses = append(addresses, *config.InterfaceIPConfigurationPropertiesFormat.PrivateIPAddress)
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
		d.Set("virtual_machine_id", *iface.VirtualMachine.ID)
	} else {
		d.Set("virtual_machine_id", "")
	}

	var appliedDNSServers []string
	var dnsServers []string
	if iface.DNSSettings != nil {
		if iface.DNSSettings.AppliedDNSServers != nil && len(*iface.DNSSettings.AppliedDNSServers) > 0 {
			for _, applied := range *iface.DNSSettings.AppliedDNSServers {
				appliedDNSServers = append(appliedDNSServers, applied)
			}
		}

		if iface.DNSSettings.DNSServers != nil && len(*iface.DNSSettings.DNSServers) > 0 {
			for _, dns := range *iface.DNSSettings.DNSServers {
				dnsServers = append(dnsServers, dns)
			}
		}

		if iface.DNSSettings.InternalFqdn != nil && *iface.DNSSettings.InternalFqdn != "" {
			d.Set("internal_fqdn", iface.DNSSettings.InternalFqdn)
		}

		d.Set("internal_dns_name_label", iface.DNSSettings.InternalDNSNameLabel)
	}

	if iface.NetworkSecurityGroup != nil {
		d.Set("network_security_group_id", resp.NetworkSecurityGroup.ID)
	} else {
		d.Set("network_security_group_id", "")
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	d.Set("applied_dns_servers", appliedDNSServers)
	d.Set("dns_servers", dnsServers)
	d.Set("enable_ip_forwarding", resp.EnableIPForwarding)
	d.Set("enable_accelerated_networking", resp.EnableAcceleratedNetworking)

	flattenAndSetTags(d, resp.Tags)

	return nil
}
