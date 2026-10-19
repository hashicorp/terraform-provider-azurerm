// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package stack_hci_logical_network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/logicalnetworks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.LogicalNetworks

			var config StackHCILogicalNetworkResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			id := logicalnetworks.NewLogicalNetworkID(subscriptionId, config.ResourceGroupName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			payload := logicalnetworks.LogicalNetworks{
				Name:     pointer.To(config.Name),
				Location: location.Normalize(config.Location),
				Tags:     tags.Expand(config.Tags),
				ExtendedLocation: &logicalnetworks.ExtendedLocation{
					Name: pointer.To(config.CustomLocationId),
					Type: pointer.To(logicalnetworks.ExtendedLocationTypesCustomLocation),
				},
				Properties: &logicalnetworks.LogicalNetworkProperties{
					VMSwitchName: pointer.To(config.VirtualSwitchName),
					Subnets:      expandStackHCILogicalNetworkSubnet(config.Subnet),
					DhcpOptions: &logicalnetworks.LogicalNetworkPropertiesDhcpOptions{
						DnsServers: pointer.To(config.DNSServers),
					},
				},
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("performing create %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func expandStackHCILogicalNetworkSubnet(input []StackHCISubnetModel) *[]logicalnetworks.Subnet {
	if len(input) == 0 {
		return nil
	}

	results := make([]logicalnetworks.Subnet, 0)
	for _, v := range input {
		results = append(results, logicalnetworks.Subnet{
			Properties: &logicalnetworks.SubnetPropertiesFormat{
				AddressPrefix:      pointer.To(v.AddressPrefix),
				IPAllocationMethod: pointer.ToEnum[logicalnetworks.IPAllocationMethodEnum](v.IpAllocationMethod),
				IPPools:            expandStackHCILogicalNetworkIPPool(v.IpPool),
				RouteTable:         expandStackHCILogicalNetworkRouteTable(v.Route),
				Vlan:               pointer.To(v.VlanId),
			},
		})
	}

	return &results
}

func expandStackHCILogicalNetworkIPPool(input []StackHCIIPPoolModel) *[]logicalnetworks.IPPool {
	if len(input) == 0 {
		return nil
	}

	results := make([]logicalnetworks.IPPool, 0)
	for _, v := range input {
		results = append(results, logicalnetworks.IPPool{
			Start: pointer.To(v.Start),
			End:   pointer.To(v.End),
		})
	}

	return &results
}

func expandStackHCILogicalNetworkRouteTable(input []StackHCIRouteModel) *logicalnetworks.RouteTable {
	if len(input) == 0 {
		return nil
	}

	routes := make([]logicalnetworks.Route, 0)
	for _, v := range input {
		route := logicalnetworks.Route{
			Properties: &logicalnetworks.RoutePropertiesFormat{
				AddressPrefix:    pointer.To(v.AddressPrefix),
				NextHopIPAddress: pointer.To(v.NextHopIpAddress),
			},
		}

		if v.Name != "" {
			route.Name = pointer.To(v.Name)
		}

		routes = append(routes, route)
	}

	return &logicalnetworks.RouteTable{
		Properties: &logicalnetworks.RouteTablePropertiesFormat{
			Routes: pointer.To(routes),
		},
	}
}
