// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_network

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualNetworkUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.VirtualNetworks
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseVirtualNetworkID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id, virtualnetworks.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("address_space") {
		if v := d.Get("address_space").(*pluginsdk.Set).List(); len(v) > 0 {
			if payload.Properties.AddressSpace == nil {
				payload.Properties.AddressSpace = &virtualnetworks.AddressSpace{}
			}
			payload.Properties.AddressSpace.AddressPrefixes = pluginsdk.ExpandStringSlice(v)
		} else {
			payload.Properties.AddressSpace.AddressPrefixes = nil
		}
	}

	if d.HasChange("ip_address_pool") {
		if v := d.Get("ip_address_pool").([]any); len(v) > 0 {
			if payload.Properties.AddressSpace == nil {
				payload.Properties.AddressSpace = &virtualnetworks.AddressSpace{}
			}

			expandedIPAddressPool := expandVirtualNetworkIPAddressPool(v)

			if payload.Properties.AddressSpace.IPamPoolPrefixAllocations != nil {
				for _, existingAllocation := range *payload.Properties.AddressSpace.IPamPoolPrefixAllocations {
					for _, expandedAllocation := range *expandedIPAddressPool {
						if existingAllocation.Pool != nil && expandedAllocation.Pool != nil && strings.EqualFold(pointer.From(existingAllocation.Pool.Id), pointer.From(expandedAllocation.Pool.Id)) &&
							existingAllocation.NumberOfIPAddresses != nil && expandedAllocation.NumberOfIPAddresses != nil {
							existingNum, _ := new(big.Int).SetString(*existingAllocation.NumberOfIPAddresses, 10)
							newNum, _ := new(big.Int).SetString(*expandedAllocation.NumberOfIPAddresses, 10)
							if existingNum != nil && newNum != nil && existingNum.Cmp(newNum) == 1 {
								return fmt.Errorf("`number_of_ip_addresses` cannot be decreased from %v to %v on pool: %v", *existingAllocation.NumberOfIPAddresses, *expandedAllocation.NumberOfIPAddresses, *expandedAllocation.Pool.Id)
							}
						}
					}
				}
			}

			payload.Properties.AddressSpace.IPamPoolPrefixAllocations = expandedIPAddressPool
		} else {
			payload.Properties.AddressSpace.IPamPoolPrefixAllocations = nil
		}
	}

	if d.HasChange("bgp_community") {
		// nil out the current values in case `bgp_community` has been removed from the config file
		payload.Properties.BgpCommunities = nil
		if v := d.Get("bgp_community"); v.(string) != "" {
			payload.Properties.BgpCommunities = &virtualnetworks.VirtualNetworkBgpCommunities{VirtualNetworkCommunity: v.(string)}
		}
	}

	if d.HasChange("ddos_protection_plan") {
		ddosProtectionPlanId, enabled := expandVirtualNetworkDdosProtectionPlan(d.Get("ddos_protection_plan").([]any))
		payload.Properties.DdosProtectionPlan = ddosProtectionPlanId
		payload.Properties.EnableDdosProtection = enabled
	}

	if d.HasChange("encryption") {
		// nil out the current values in case `encryption` has been removed from the config file
		payload.Properties.Encryption = expandVirtualNetworkEncryption(d.Get("encryption").([]any))
	}

	if d.HasChange("dns_servers") {
		if payload.Properties.DhcpOptions == nil {
			payload.Properties.DhcpOptions = &virtualnetworks.DhcpOptions{}
		}

		payload.Properties.DhcpOptions.DnsServers = pluginsdk.ExpandStringSlice(d.Get("dns_servers").([]any))
	}

	if d.HasChange("flow_timeout_in_minutes") {
		payload.Properties.FlowTimeoutInMinutes = nil
		if v := d.Get("flow_timeout_in_minutes"); v.(int) != 0 {
			payload.Properties.FlowTimeoutInMinutes = pointer.To(int64(v.(int)))
		}
	}

	if d.HasChange("subnet") {
		subnetList, routeTables, err := expandVirtualNetworkSubnets(ctx, *client, d.Get("subnet").(*pluginsdk.Set).List(), *id)
		if err != nil {
			return fmt.Errorf("expanding `subnet`: %+v", err)
		}
		payload.Properties.Subnets = subnetList

		locks.ByID(id.ID())
		defer locks.UnlockByID(id.ID())

		locks.MultipleByID(routeTables)
		defer locks.UnlockMultipleByID(routeTables)
	}

	if d.HasChange("private_endpoint_vnet_policies") {
		payload.Properties.PrivateEndpointVNetPolicies = pointer.ToEnum[virtualnetworks.PrivateEndpointVNetPolicies](d.Get("private_endpoint_vnet_policies").(string))
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	networkSecurityGroupIds := make([]string, 0)
	if payload.Properties != nil && payload.Properties.Subnets != nil {
		for _, subnet := range *payload.Properties.Subnets {
			if subnet.Properties != nil {
				// remove readonly properties as they are not managed by TF - large networks can cause ARM API limit errors
				subnet.Properties.IPConfigurations = nil
				subnet.Properties.PrivateEndpoints = nil

				if subnet.Properties.NetworkSecurityGroup != nil && subnet.Properties.NetworkSecurityGroup.Id != nil {
					parsedNsgID, err := networksecuritygroups.ParseNetworkSecurityGroupID(*subnet.Properties.NetworkSecurityGroup.Id)
					if err != nil {
						return err
					}
					networkSecurityGroupIds = append(networkSecurityGroupIds, parsedNsgID.ID())
				}
			}
		}
	}

	locks.MultipleByID(&networkSecurityGroupIds)
	defer locks.UnlockMultipleByID(&networkSecurityGroupIds)

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	timeout, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(virtualnetworks.ProvisioningStateUpdating)},
		Target:     []string{string(virtualnetworks.ProvisioningStateSucceeded)},
		Refresh:    helpers.VirtualNetworkProvisioningStateRefreshFunc(ctx, meta.(*clients.Client).Network.VirtualNetworks, *id),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(timeout),
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVirtualNetworkRead(d, meta)
}
