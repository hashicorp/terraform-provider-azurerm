// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/networkinterfaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

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

			// TODO 4.0: change this from enable_* to *_enabled
			"enable_accelerated_networking": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			// TODO 4.0: change this from enable_* to *_enabled
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

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceNetworkInterfaceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NetworkInterfaces
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewNetworkInterfaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("Error: %s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.NetworkInterfaceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	model := resp.Model
	if model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if location := model.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := model.Properties; props != nil {
		d.Set("mac_address", props.MacAddress)

		privateIpAddress := ""
		privateIpAddresses := make([]interface{}, 0)
		if configs := props.IPConfigurations; configs != nil {
			for _, config := range *configs {
				if config.Properties == nil {
					continue
				}
				if config.Properties.PrivateIPAddress == nil {
					continue
				}

				ipAddress := *config.Properties.PrivateIPAddress
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
		if props.VirtualMachine != nil && props.VirtualMachine.Id != nil {
			virtualMachineId = *props.VirtualMachine.Id
		}
		d.Set("virtual_machine_id", virtualMachineId)

		var appliedDNSServers []string
		var dnsServers []string
		if dnsSettings := props.DnsSettings; dnsSettings != nil {
			if s := dnsSettings.AppliedDnsServers; s != nil {
				appliedDNSServers = *s
			}

			if s := dnsSettings.DnsServers; s != nil {
				dnsServers = *s
			}

			d.Set("internal_dns_name_label", dnsSettings.InternalDnsNameLabel)
		}

		networkSecurityGroupId := ""
		if props.NetworkSecurityGroup != nil && props.NetworkSecurityGroup.Id != nil {
			networkSecurityGroupId = *props.NetworkSecurityGroup.Id
		}
		d.Set("network_security_group_id", networkSecurityGroupId)

		d.Set("applied_dns_servers", appliedDNSServers)
		d.Set("dns_servers", dnsServers)
		d.Set("enable_ip_forwarding", props.EnableIPForwarding)
		d.Set("enable_accelerated_networking", props.EnableAcceleratedNetworking)
	}

	return tags.FlattenAndSet(d, model.Tags)
}
