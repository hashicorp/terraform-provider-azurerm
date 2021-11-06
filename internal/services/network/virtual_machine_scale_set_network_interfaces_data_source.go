package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-02-01/network"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceVirtualMachineScaleSetNetworkInterfaces() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceVirtualMachineScaleSetNetworkInterfacesRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"virtual_machine_scale_set_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"network_interfaces": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"location": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"network_security_group_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"mac_address": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"virtual_machine_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"ip_configuration": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"subnet_id": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"private_ip_address": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"private_ip_address_version": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"private_ip_address_allocation": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"public_ip_address_id": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},

									"application_gateway_backend_address_pools_ids": {
										Type:     pluginsdk.TypeSet,
										Computed: true,
										Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
										Set:      pluginsdk.HashString,
									},

									"load_balancer_backend_address_pools_ids": {
										Type:     pluginsdk.TypeSet,
										Computed: true,
										Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
										Set:      pluginsdk.HashString,
									},

									"load_balancer_inbound_nat_rules_ids": {
										Type:     pluginsdk.TypeSet,
										Computed: true,
										Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
										Set:      pluginsdk.HashString,
									},

									"application_security_group_ids": {
										Type:     pluginsdk.TypeSet,
										Computed: true,
										Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
										Set:      pluginsdk.HashString,
									},

									"primary": {
										Type:     pluginsdk.TypeBool,
										Computed: true,
									},

									"gateway_load_balancer_frontend_ip_configuration_id": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
								},
							},
						},

						"dns_servers": {
							Type:     pluginsdk.TypeSet,
							Computed: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Set:      pluginsdk.HashString,
						},

						"internal_dns_name_label": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"applied_dns_servers": {
							Type:     pluginsdk.TypeSet,
							Computed: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
							Set:      pluginsdk.HashString,
						},

						"enable_accelerated_networking": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"enable_ip_forwarding": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"private_ip_address": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"private_ip_addresses": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},

						"tags": tags.SchemaDataSource(),
					},
				},
			},
		},
	}
}

func dataSourceVirtualMachineScaleSetNetworkInterfacesRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.InterfacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	virtualMachineScaleSetName := d.Get("virtual_machine_scale_set_name").(string)

	log.Printf("[DEBUG] Reading Network Interfaces of Virtual Machine Scale Set %q (Resource Group %q)", virtualMachineScaleSetName, resourceGroup)
	resp, err := client.ListVirtualMachineScaleSetNetworkInterfacesComplete(ctx, resourceGroup, virtualMachineScaleSetName)
	if err != nil {
		return fmt.Errorf("listing Network Interfaces of Virtual Machine Scale Set %q (Resource Group %q): %v", virtualMachineScaleSetName, resourceGroup, err)
	}

	networkInterfaces := make([]network.Interface, 0)
	for resp.NotDone() {
		networkInterfaces = append(networkInterfaces, resp.Value())
		if err := resp.NextWithContext(ctx); err != nil {
			return fmt.Errorf("listing next page of Network Interfaces of Virtual Machine Scale Set %q (Resource Group %q): %v", virtualMachineScaleSetName, resourceGroup, err)
		}
	}

	d.SetId(fmt.Sprintf("%s-%s", virtualMachineScaleSetName, resourceGroup))

	d.Set("resource_group_name", resourceGroup)
	d.Set("virtual_machine_scale_set_name", virtualMachineScaleSetName)
	results := flattenDataSourceNetworkInterfaces(networkInterfaces)
	if err := d.Set("network_interfaces", results); err != nil {
		return fmt.Errorf("setting `network_interfaces`: %+v", err)
	}

	return nil
}

func flattenDataSourceNetworkInterfaces(input []network.Interface) []interface{} {
	results := make([]interface{}, 0)

	for _, element := range input {
		var name, location string
		if element.Name != nil {
			name = *element.Name
		}
		if element.Location != nil {
			location = azure.NormalizeLocation(*element.Location)
		}

		var enableIPForwarding, enableAcceleratedNetworking bool
		var macAddress, privateIpAddress, virtualMachineId, internalDNSNameLabel, networkSecurityGroupId string
		var privateIpAddresses, ipConfigurations []interface{}
		var appliedDNSServers, dnsServers []string
		if props := element.InterfacePropertiesFormat; props != nil {
			if s := props.EnableIPForwarding; s != nil {
				enableIPForwarding = *s
			}

			if s := props.EnableAcceleratedNetworking; s != nil {
				enableAcceleratedNetworking = *s
			}

			if s := props.MacAddress; s != nil {
				macAddress = *s
			}

			privateIpAddress = ""
			privateIpAddresses = make([]interface{}, 0)
			if configs := props.IPConfigurations; configs != nil {
				for _, config := range *configs {
					if config.InterfaceIPConfigurationPropertiesFormat == nil {
						continue
					}
					if config.InterfaceIPConfigurationPropertiesFormat.PrivateIPAddress == nil {
						continue
					}

					ipAddress := *config.InterfaceIPConfigurationPropertiesFormat.PrivateIPAddress
					if privateIpAddress == "" {
						privateIpAddress = ipAddress
					}

					privateIpAddresses = append(privateIpAddresses, ipAddress)
				}
			}

			ipConfigurations = flattenNetworkInterfaceIPConfigurations(props.IPConfigurations)

			if props.VirtualMachine != nil && props.VirtualMachine.ID != nil {
				virtualMachineId = *props.VirtualMachine.ID
			}

			if dnsSettings := props.DNSSettings; dnsSettings != nil {
				if s := dnsSettings.AppliedDNSServers; s != nil {
					appliedDNSServers = *s
				}

				if s := dnsSettings.DNSServers; s != nil {
					dnsServers = *s
				}

				if s := dnsSettings.InternalDNSNameLabel; s != nil {
					internalDNSNameLabel = *s
				}
			}

			if props.NetworkSecurityGroup != nil && props.NetworkSecurityGroup.ID != nil {
				networkSecurityGroupId = *props.NetworkSecurityGroup.ID
			}
		}

		results = append(results, map[string]interface{}{
			"name":                          name,
			"location":                      location,
			"network_security_group_id":     networkSecurityGroupId,
			"mac_address":                   macAddress,
			"virtual_machine_id":            virtualMachineId,
			"ip_configuration":              ipConfigurations,
			"dns_servers":                   dnsServers,
			"internal_dns_name_label":       internalDNSNameLabel,
			"applied_dns_servers":           appliedDNSServers,
			"enable_accelerated_networking": enableAcceleratedNetworking,
			"enable_ip_forwarding":          enableIPForwarding,
			"private_ip_address":            privateIpAddress,
			"private_ip_addresses":          privateIpAddresses,
			"tags":                          tags.Flatten(element.Tags),
		})
	}

	return results
}
