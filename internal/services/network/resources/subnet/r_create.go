// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package subnet

import (
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routetables"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/subnets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSubnetCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.Subnets
	vnetClient := meta.(*clients.Client).Network.VirtualNetworks
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewSubnetID(meta.(*clients.Client).Account.SubscriptionId, d.Get("resource_group_name").(string), d.Get("virtual_network_name").(string), d.Get("name").(string))

	nsg, err := expandSubnetNetworkSecurityGroupID(d)
	if err != nil {
		return fmt.Errorf("expanding Network Security Group ID: %+v", err)
	}
	if nsg != nil && nsg.Id != nil {
		nsgId, err := networksecuritygroups.ParseNetworkSecurityGroupID(*nsg.Id)
		if err != nil {
			return err
		}
		locks.ByID(nsgId.ID())
		defer locks.UnlockByID(nsgId.ID())
	}

	rt, err := expandSubnetRouteTableID(d)
	if err != nil {
		return fmt.Errorf("expanding Route Table ID: %+v", err)
	}
	if rt != nil && rt.Id != nil {
		rtId, err := routetables.ParseRouteTableID(*rt.Id)
		if err != nil {
			return err
		}
		locks.ByID(rtId.ID())
		defer locks.UnlockByID(rtId.ID())
	}

	vnetId := commonids.NewVirtualNetworkID(id.SubscriptionId, id.ResourceGroupName, id.VirtualNetworkName)
	locks.ByID(vnetId.ID())
	defer locks.UnlockByID(vnetId.ID())

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	existing, err := client.Get(ctx, id, subnets.DefaultGetOperationOptions())
	if !response.WasNotFound(existing.HttpResponse) {
		if err != nil {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
		return tf.ImportAsExistsError("azurerm_subnet", id.ID())
	}

	subnet := subnets.Subnet{
		Name: pointer.To(id.SubnetName),
		Properties: &subnets.SubnetPropertiesFormat{
			AddressPrefixes:                   pluginsdk.ExpandStringSlice(d.Get("address_prefixes").([]any)),
			DefaultOutboundAccess:             pointer.To(d.Get("default_outbound_access_enabled").(bool)),
			Delegations:                       expandSubnetDelegation(d.Get("delegation").([]any)),
			IPamPoolPrefixAllocations:         expandSubnetIPAddressPool(d.Get("ip_address_pool").([]any)),
			PrivateEndpointNetworkPolicies:    pointer.ToEnum[subnets.VirtualNetworkPrivateEndpointNetworkPolicies](d.Get("private_endpoint_network_policies").(string)),
			PrivateLinkServiceNetworkPolicies: expandSubnetNetworkPolicy(d.Get("private_link_service_network_policies_enabled").(bool)),
			ServiceEndpoints:                  expandSubnetServiceEndpoint(d.Get("service_endpoint").([]any)),
			ServiceEndpointPolicies:           expandSubnetServiceEndpointPolicies(d.Get("service_endpoint_policy_ids").(*pluginsdk.Set).List()),
			SharingScope:                      pointer.ToEnum[subnets.SharingScope](d.Get("sharing_scope").(string)),

			// Write-only
			NetworkSecurityGroup: nsg,
			RouteTable:           rt,
		},
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, subnet, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return errors.New("internal-error: context had no deadline")
	}

	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(subnets.ProvisioningStateUpdating)},
		Target:     []string{string(subnets.ProvisioningStateSucceeded)},
		Refresh:    helpers.SubnetProvisioningStateRefreshFunc(ctx, client, id),
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
	if _, err := vnetStateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of virtual network for %s: %+v", id, err)
	}

	return resourceSubnetRead(d, meta)
}
