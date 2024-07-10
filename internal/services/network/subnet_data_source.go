// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/subnets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceSubnet() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Read: dataSourceSubnetRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"virtual_network_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"address_prefix": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"address_prefixes": {
				Type:     pluginsdk.TypeList,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				Computed: true,
			},

			"network_security_group_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"route_table_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"service_endpoints": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"default_outbound_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"private_endpoint_network_policies": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"private_link_service_network_policies_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},
		},
	}

	if !features.FourPointOhBeta() {
		resource.Schema["private_endpoint_network_policies_enabled"] = &pluginsdk.Schema{
			Type:       pluginsdk.TypeBool,
			Computed:   true,
			Deprecated: "This property has been superseded by `private_endpoint_network_policies` and will be removed in v4.0 of the AzureRM Provider.",
		}

		resource.Schema["enforce_private_link_endpoint_network_policies"] = &pluginsdk.Schema{
			Type:       pluginsdk.TypeBool,
			Computed:   true,
			Deprecated: "This property has been superseded by `private_endpoint_network_policies` and will be removed in v4.0 of the AzureRM Provider.",
		}

		resource.Schema["enforce_private_link_service_network_policies"] = &pluginsdk.Schema{
			Type:       pluginsdk.TypeBool,
			Computed:   true,
			Deprecated: "This property has been superseded by `private_link_service_network_policies_enabled` and will be removed in v4.0 of the AzureRM Provider.",
		}
	}

	return resource
}

func dataSourceSubnetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Client.Subnets
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewSubnetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("virtual_network_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id, subnets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.SubnetName)
	d.Set("virtual_network_name", id.VirtualNetworkName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("address_prefix", props.AddressPrefix)
			if props.AddressPrefixes == nil {
				if props.AddressPrefix != nil && len(*props.AddressPrefix) > 0 {
					d.Set("address_prefixes", []string{*props.AddressPrefix})
				} else {
					d.Set("address_prefixes", []string{})
				}
			} else {
				d.Set("address_prefixes", utils.FlattenStringSlice(props.AddressPrefixes))
			}

			if !features.FourPointOhBeta() {
				d.Set("enforce_private_link_endpoint_network_policies", flattenEnforceSubnetNetworkPolicy(string(*props.PrivateEndpointNetworkPolicies)))
				d.Set("private_endpoint_network_policies_enabled", flattenSubnetNetworkPolicy(string(*props.PrivateEndpointNetworkPolicies)))
				d.Set("enforce_private_link_service_network_policies", flattenEnforceSubnetNetworkPolicy(string(*props.PrivateLinkServiceNetworkPolicies)))
			}

			defaultOutboundAccessEnabled := true
			if props.DefaultOutboundAccess != nil {
				defaultOutboundAccessEnabled = *props.DefaultOutboundAccess
			}
			d.Set("default_outbound_access_enabled", defaultOutboundAccessEnabled)

			d.Set("private_endpoint_network_policies", string(pointer.From(props.PrivateEndpointNetworkPolicies)))
			d.Set("private_link_service_network_policies_enabled", flattenSubnetNetworkPolicy(string(*props.PrivateLinkServiceNetworkPolicies)))

			networkSecurityGroupId := ""
			if props.NetworkSecurityGroup != nil && props.NetworkSecurityGroup.Id != nil {
				networkSecurityGroupId = *props.NetworkSecurityGroup.Id
			}
			d.Set("network_security_group_id", networkSecurityGroupId)

			routeTableId := ""
			if props.RouteTable != nil && props.RouteTable.Id != nil {
				routeTableId = *props.RouteTable.Id
			}
			d.Set("route_table_id", routeTableId)

			serviceEndpoints := flattenSubnetServiceEndpoints(props.ServiceEndpoints)
			if err := d.Set("service_endpoints", serviceEndpoints); err != nil {
				return fmt.Errorf("setting `service_endpoints`: %+v", err)
			}
		}
	}

	return nil
}
