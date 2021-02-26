package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
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
