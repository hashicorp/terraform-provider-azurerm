package network

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceNetworkInterface() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceNetworkInterfaceRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

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
	}
}

func dataSourceNetworkInterfaceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.InterfacesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewNetworkInterfaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: %s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.InterfacePropertiesFormat; props != nil {
		d.Set("mac_address", props.MacAddress)

		privateIpAddress := ""
		privateIpAddresses := make([]interface{}, 0)
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
		d.Set("private_ip_address", privateIpAddress)
		if err := d.Set("private_ip_addresses", privateIpAddresses); err != nil {
			return fmt.Errorf("setting `private_ip_addresses`: %+v", err)
		}

		if err := d.Set("ip_configuration", flattenNetworkInterfaceIPConfigurations(props.IPConfigurations)); err != nil {
			return fmt.Errorf("setting `ip_configuration`: %+v", err)
		}

		virtualMachineId := ""
		if props.VirtualMachine != nil && props.VirtualMachine.ID != nil {
			virtualMachineId = *props.VirtualMachine.ID
		}
		d.Set("virtual_machine_id", virtualMachineId)

		var appliedDNSServers []string
		var dnsServers []string
		if dnsSettings := props.DNSSettings; dnsSettings != nil {
			if s := dnsSettings.AppliedDNSServers; s != nil {
				appliedDNSServers = *s
			}

			if s := dnsSettings.DNSServers; s != nil {
				dnsServers = *s
			}

			d.Set("internal_dns_name_label", dnsSettings.InternalDNSNameLabel)
		}

		networkSecurityGroupId := ""
		if props.NetworkSecurityGroup != nil && props.NetworkSecurityGroup.ID != nil {
			networkSecurityGroupId = *props.NetworkSecurityGroup.ID
		}
		d.Set("network_security_group_id", networkSecurityGroupId)

		d.Set("applied_dns_servers", appliedDNSServers)
		d.Set("dns_servers", dnsServers)
		d.Set("enable_ip_forwarding", props.EnableIPForwarding)
		d.Set("enable_accelerated_networking", props.EnableAcceleratedNetworking)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
