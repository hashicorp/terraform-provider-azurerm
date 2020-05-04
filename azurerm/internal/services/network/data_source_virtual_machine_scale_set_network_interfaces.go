package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func dataSourceArmVirtualMachineScaleSetNetworkInterfaces() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmVirtualMachineScaleSetNetworkInterfacesRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"virtual_machine_scale_set_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"network_interfaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"name": {
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

						"id": {
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
					},
				},
			},
		},
	}
}

func dataSourceArmVirtualMachineScaleSetNetworkInterfacesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.InterfacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resGroup := d.Get("resource_group_name").(string)
	vmssName := d.Get("virtual_machine_scale_set_name").(string)

	resp, err := client.ListVirtualMachineScaleSetNetworkInterfaces(ctx, resGroup, vmssName)
	if err != nil {
		return fmt.Errorf("Error making Read request on Azure VMSS %q (Resource Group %q): %+v", vmssName, resGroup, err)
	}

	results := make([]interface{}, 0)

	for _, iface := range resp.Values() {
		result := make(map[string]interface{})

		if iface.Name != nil {
			result["name"] = *iface.Name
		} else {
			result["name"] = ""
		}

		if iface.NetworkSecurityGroup != nil {
			result["network_security_group_id"] = *iface.NetworkSecurityGroup.ID
		} else {
			result["network_security_group_id"] = ""
		}

		if iface.MacAddress != nil {
			result["mac_address"] = *iface.MacAddress
		} else {
			result["mac_address"] = ""
		}

		if iface.VirtualMachine != nil {
			result["virtual_machine_id"] = *iface.VirtualMachine.ID
		} else {
			result["virtual_machine_id"] = ""
		}

		if iface.ID != nil {
			result["id"] = *iface.ID
		} else {
			result["id"] = ""
		}

		if iface.IPConfigurations != nil && len(*iface.IPConfigurations) > 0 {
			configs := *iface.IPConfigurations

			result["private_ip_address"] = *configs[0].InterfaceIPConfigurationPropertiesFormat.PrivateIPAddress

			addresses := make([]interface{}, 0)
			for _, config := range configs {
				if config.InterfaceIPConfigurationPropertiesFormat != nil {
					addresses = append(addresses, config.InterfaceIPConfigurationPropertiesFormat.PrivateIPAddress)
				}
			}
			result["private_ip_addresses"] = addresses
		}

		if iface.IPConfigurations != nil {
			result["ip_configuration"] = flattenNetworkInterfaceIPConfigurations(iface.IPConfigurations)
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

			if dnsSettings.InternalDNSNameLabel != nil {
				result["internal_dns_name_label"] = *dnsSettings.InternalDNSNameLabel
			} else {
				result["internal_dns_name_label"] = ""
			}
		}

		result["applied_dns_servers"] = appliedDNSServers
		result["dns_servers"] = dnsServers
		result["enable_ip_forwarding"] = *iface.EnableIPForwarding
		result["enable_accelerated_networking"] = *iface.EnableAcceleratedNetworking
		results = append(results, result)
	}

	d.SetId(time.Now().UTC().String())
	if err := d.Set("network_interfaces", results); err != nil {
		return fmt.Errorf("Error setting `network_interfaces`: %+v", err)
	}
	return nil
}
