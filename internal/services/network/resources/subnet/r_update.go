// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package subnet

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routetables"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/subnets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSubnetUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.Subnets
	vnetClient := meta.(*clients.Client).Network.VirtualNetworks

	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseSubnetID(d.Id())
	if err != nil {
		return err
	}

	existingUnlocked, err := client.Get(ctx, *id, subnets.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if existingUnlocked.Model == nil || existingUnlocked.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `model` or `properties` was nil", *id)
	}

	propsUnlocked := existingUnlocked.Model.Properties

	var nsgIds []string
	if propsUnlocked.NetworkSecurityGroup != nil && propsUnlocked.NetworkSecurityGroup.Id != nil {
		nsgId, err := networksecuritygroups.ParseNetworkSecurityGroupID(*propsUnlocked.NetworkSecurityGroup.Id)
		if err != nil {
			return fmt.Errorf("parsing existing Network Security Group ID: %+v", err)
		}
		nsgIds = append(nsgIds, nsgId.ID())
	}

	if d.HasChange("network_security_group_id_wo_version") {
		nsg, err := expandSubnetNetworkSecurityGroupID(d)
		if err != nil {
			return fmt.Errorf("expanding Network Security Group ID: %+v", err)
		}
		if nsg != nil && nsg.Id != nil {
			nsgId, err := networksecuritygroups.ParseNetworkSecurityGroupID(*nsg.Id)
			if err != nil {
				return err
			}
			nsgIds = append(nsgIds, nsgId.ID())
		}
	}

	var rtIds []string
	if propsUnlocked.RouteTable != nil && propsUnlocked.RouteTable.Id != nil {
		rtId, err := routetables.ParseRouteTableID(*propsUnlocked.RouteTable.Id)
		if err != nil {
			return fmt.Errorf("parsing existing Route Table ID: %+v", err)
		}
		rtIds = append(rtIds, rtId.ID())
	}

	if d.HasChange("route_table_id_wo_version") {
		rt, err := expandSubnetRouteTableID(d)
		if err != nil {
			return fmt.Errorf("expanding Route Table ID: %+v", err)
		}
		if rt != nil && rt.Id != nil {
			rtId, err := routetables.ParseRouteTableID(*rt.Id)
			if err != nil {
				return err
			}
			rtIds = append(rtIds, rtId.ID())
		}
	}

	locks.MultipleByID(&nsgIds)
	defer locks.UnlockMultipleByID(&nsgIds)

	locks.MultipleByID(&rtIds)
	defer locks.UnlockMultipleByID(&rtIds)

	vnetId := commonids.NewVirtualNetworkID(id.SubscriptionId, id.ResourceGroupName, id.VirtualNetworkName)
	locks.ByID(vnetId.ID())
	defer locks.UnlockByID(vnetId.ID())

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	existing, err := client.Get(ctx, *id, subnets.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if existing.Model == nil || existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `model` or `properties` was nil", *id)
	}

	props := existing.Model.Properties

	if d.HasChange("address_prefixes") {
		addressPrefixesRaw := d.Get("address_prefixes").([]any)
		switch len(addressPrefixesRaw) {
		case 0:
			// this is the case IPAddressPool is used, so we shall clear the `AddressPrefix` and `AddressPrefixes`.
			props.AddressPrefix = nil
			props.AddressPrefixes = nil
		default:
			props.AddressPrefixes = pluginsdk.ExpandStringSlice(addressPrefixesRaw)
			props.AddressPrefix = nil
		}
	}

	if d.HasChange("ip_address_pool") {
		if v := d.Get("ip_address_pool").([]any); len(v) > 0 {
			expandedIPAddressPool := expandSubnetIPAddressPool(d.Get("ip_address_pool").([]any))

			if props.IPamPoolPrefixAllocations != nil {
				for _, existingAllocation := range *props.IPamPoolPrefixAllocations {
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

			props.IPamPoolPrefixAllocations = expandedIPAddressPool

			// Set nil for AddressPrefixes when changing from `AddressPrefixes` to `IPAddressPool` but the change for `AddressPrefix` is not detected due to the diffSuppressFunc.
			props.AddressPrefix = nil
			props.AddressPrefixes = nil
		} else {
			props.IPamPoolPrefixAllocations = nil
		}
	}

	if d.HasChange("default_outbound_access_enabled") {
		props.DefaultOutboundAccess = pointer.To(d.Get("default_outbound_access_enabled").(bool))
	}

	if d.HasChange("delegation") {
		delegationsRaw := d.Get("delegation").([]any)
		props.Delegations = expandSubnetDelegation(delegationsRaw)
	}

	if d.HasChange("private_endpoint_network_policies") {
		v := d.Get("private_endpoint_network_policies").(string)
		props.PrivateEndpointNetworkPolicies = pointer.ToEnum[subnets.VirtualNetworkPrivateEndpointNetworkPolicies](v)
	}

	if d.HasChange("private_link_service_network_policies_enabled") {
		props.PrivateLinkServiceNetworkPolicies = expandSubnetNetworkPolicy(d.Get("private_link_service_network_policies_enabled").(bool))
	}

	if d.HasChange("sharing_scope") {
		props.SharingScope = pointer.ToEnum[subnets.SharingScope](d.Get("sharing_scope").(string))
	}

	if d.HasChange("service_endpoint") {
		props.ServiceEndpoints = expandSubnetServiceEndpoint(d.Get("service_endpoint").([]any))
	}

	if d.HasChange("service_endpoint_policy_ids") {
		props.ServiceEndpointPolicies = expandSubnetServiceEndpointPolicies(d.Get("service_endpoint_policy_ids").(*pluginsdk.Set).List())
	}

	if d.HasChange("network_security_group_id_wo_version") {
		nsg, err := expandSubnetNetworkSecurityGroupID(d)
		if err != nil {
			return fmt.Errorf("expanding Network Security Group ID: %+v", err)
		}

		props.NetworkSecurityGroup = nsg
	}

	if d.HasChange("route_table_id_wo_version") {
		rt, err := expandSubnetRouteTableID(d)
		if err != nil {
			return fmt.Errorf("expanding Route Table ID: %+v", err)
		}

		props.RouteTable = rt
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return errors.New("internal-error: context had no deadline")
	}

	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(subnets.ProvisioningStateUpdating)},
		Target:     []string{string(subnets.ProvisioningStateSucceeded)},
		Refresh:    helpers.SubnetProvisioningStateRefreshFunc(ctx, client, *id),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(deadline),
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of %s: %+v", id, err)
	}

	vnetStateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(subnets.ProvisioningStateUpdating)},
		Target:     []string{string(subnets.ProvisioningStateSucceeded)},
		Refresh:    helpers.VirtualNetworkProvisioningStateRefreshFunc(ctx, vnetClient, vnetId),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(deadline),
	}

	if _, err = vnetStateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of virtual network for %s: %+v", id, err)
	}

	return resourceSubnetRead(d, meta)
}
