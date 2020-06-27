package network

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-03-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},

									"load_balancer_backend_address_pools_ids": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},

									"load_balancer_inbound_nat_rules_ids": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},

									"application_security_group_ids": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},

									"primary": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},

						"dns_servers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

						"internal_dns_name_label": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"applied_dns_servers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
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
							Elem:     &schema.Schema{Type: schema.TypeString},
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

	resp, err := client.ListVirtualMachineScaleSetNetworkInterfacesComplete(ctx, resGroup, vmssName)
	if err != nil {
		return fmt.Errorf("Error making Read request on Azure VMSS %q (Resource Group %q): %+v", vmssName, resGroup, err)
	}

	results := make([]interface{}, 0)

	for resp.NotDone() {
		iface := resp.Value()

		name := ""
		if iface.Name != nil {
			name = *iface.Name
		}

		networkSecurityGroupID := ""
		if iface.NetworkSecurityGroup != nil {
			networkSecurityGroupID = *iface.NetworkSecurityGroup.ID
		}

		macAddress := ""
		if iface.MacAddress != nil {
			macAddress = *iface.MacAddress
		}

		virtualMachineID := ""
		if iface.VirtualMachine != nil {
			virtualMachineID = *iface.VirtualMachine.ID
		}

		id := ""
		if iface.ID != nil {
			id = *iface.ID
		}

		ipConfiguration := flattenNetworkInterfaceIPConfigurations(iface.IPConfigurations)

		privateIPAddress := ""
		privateIPAddresses := flattenNetworkInterfacePrivateIPAddresses(iface.IPConfigurations)
		if len(privateIPAddresses) > 0 {
			privateIPAddress = privateIPAddresses[0].(string)
		}

		dnsServers := make([]interface{}, 0)
		appliedDNSServers := make([]interface{}, 0)
		internalDNSNameLabel := ""
		if dnsSettings := iface.DNSSettings; dnsSettings != nil {
			dnsServers = utils.FlattenStringSlice(dnsSettings.DNSServers)
			appliedDNSServers = utils.FlattenStringSlice(dnsSettings.AppliedDNSServers)
			if dnsSettings.InternalDNSNameLabel != nil {
				internalDNSNameLabel = *dnsSettings.InternalDNSNameLabel
			}
		}

		enableIPForwarding := false
		if iface.EnableIPForwarding != nil {
			enableIPForwarding = *iface.EnableIPForwarding
		}
		enableAcceleratedNetworking := false
		if iface.EnableAcceleratedNetworking != nil {
			enableAcceleratedNetworking = *iface.EnableAcceleratedNetworking
		}
		results = append(results, map[string]interface{}{
			"name":                          name,
			"network_security_group_id":     networkSecurityGroupID,
			"mac_address":                   macAddress,
			"virtual_machine_id":            virtualMachineID,
			"id":                            id,
			"ip_configuration":              ipConfiguration,
			"dns_servers":                   dnsServers,
			"internal_dns_name_label":       internalDNSNameLabel,
			"applied_dns_servers":           appliedDNSServers,
			"enable_ip_forwarding":          enableIPForwarding,
			"enable_accelerated_networking": enableAcceleratedNetworking,
			"private_ip_address":            privateIPAddress,
			"private_ip_addresses":          privateIPAddresses,
		})
	}

	d.SetId(time.Now().UTC().String())
	if err := d.Set("network_interfaces", results); err != nil {
		return fmt.Errorf("Error setting `network_interfaces`: %+v", err)
	}
	return nil
}

func flattenNetworkInterfacePrivateIPAddresses(input *[]network.InterfaceIPConfiguration) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	result := make([]interface{}, 0)
	for _, ipConfig := range *input {
		if ipConfig.InterfaceIPConfigurationPropertiesFormat != nil {
			result = append(result, ipConfig.InterfaceIPConfigurationPropertiesFormat.PrivateIPAddress)
		}
	}
	return result
}
