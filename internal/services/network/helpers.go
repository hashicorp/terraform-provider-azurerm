// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/subnets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/virtualnetworks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// subnetAndVnetPoller waits for both a subnet and its parent VNet to reach ProvisioningStateSucceeded.
type subnetAndVnetPoller struct {
	subnetClient *subnets.SubnetsClient
	vnetClient   *virtualnetworks.VirtualNetworksClient
	subnetID     *commonids.SubnetId
	deadline     time.Time
}

// NewSubnetAndVnetPoller creates a new poller for a subnet and its parent VNet.
func NewSubnetAndVnetPoller(subnetClient *subnets.SubnetsClient, vnetClient *virtualnetworks.VirtualNetworksClient, subnetID *commonids.SubnetId, deadline time.Time) *subnetAndVnetPoller {
	return &subnetAndVnetPoller{
		subnetClient: subnetClient,
		vnetClient:   vnetClient,
		subnetID:     subnetID,
		deadline:     deadline,
	}
}

// Poll waits for both the subnet and its parent VNet to reach ProvisioningStateSucceeded.
func (p *subnetAndVnetPoller) Poll(ctx context.Context) error {
	// Poll the subnet
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(subnets.ProvisioningStateUpdating)},
		Target:     []string{string(subnets.ProvisioningStateSucceeded)},
		Refresh:    SubnetProvisioningStateRefreshFunc(ctx, p.subnetClient, *p.subnetID),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(p.deadline),
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("subnet %s provisioning: %w", p.subnetID.ID(), err)
	}

	// Poll the parent VNet
	vnetID := commonids.NewVirtualNetworkID(
		p.subnetID.SubscriptionId,
		p.subnetID.ResourceGroupName,
		p.subnetID.VirtualNetworkName,
	)
	vnetConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(subnets.ProvisioningStateUpdating)},
		Target:     []string{string(subnets.ProvisioningStateSucceeded)},
		Refresh:    VirtualNetworkProvisioningStateRefreshFunc(ctx, p.vnetClient, vnetID),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(p.deadline),
	}
	if _, err := vnetConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("vnet %s provisioning: %w", vnetID, err)
	}

	return nil
}
