// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/subnets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// subnetAndVnetPoller waits for both a subnet and its parent VNet to reach ProvisioningStateSucceeded.
type SubnetAndVnetPoller struct {
	SubnetClient *subnets.SubnetsClient
	VnetClient   *virtualnetworks.VirtualNetworksClient
	SubnetID     *commonids.SubnetId
	Deadline     time.Time
}

// NewSubnetAndVnetPoller creates a new poller for a subnet and its parent VNet.
func NewSubnetAndVnetPoller(subnetClient *subnets.SubnetsClient, vnetClient *virtualnetworks.VirtualNetworksClient, subnetID *commonids.SubnetId, deadline time.Time) *SubnetAndVnetPoller {
	return &SubnetAndVnetPoller{
		SubnetClient: subnetClient,
		VnetClient:   vnetClient,
		SubnetID:     subnetID,
		Deadline:     deadline,
	}
}

// Poll waits for both the subnet and its parent VNet to reach ProvisioningStateSucceeded.
func (p *SubnetAndVnetPoller) Poll(ctx context.Context) error {
	// Poll the subnet
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(subnets.ProvisioningStateUpdating)},
		Target:     []string{string(subnets.ProvisioningStateSucceeded)},
		Refresh:    SubnetProvisioningStateRefreshFunc(ctx, p.SubnetClient, *p.SubnetID),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(p.Deadline),
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("subnet %s provisioning: %w", p.SubnetID.ID(), err)
	}

	// Poll the parent VNet
	vnetID := commonids.NewVirtualNetworkID(
		p.SubnetID.SubscriptionId,
		p.SubnetID.ResourceGroupName,
		p.SubnetID.VirtualNetworkName,
	)
	vnetConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(subnets.ProvisioningStateUpdating)},
		Target:     []string{string(subnets.ProvisioningStateSucceeded)},
		Refresh:    VirtualNetworkProvisioningStateRefreshFunc(ctx, p.VnetClient, vnetID),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(p.Deadline),
	}
	if _, err := vnetConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("vnet %s provisioning: %w", vnetID, err)
	}

	return nil
}
