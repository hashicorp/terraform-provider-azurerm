package custompollers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/subnets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/virtualnetworks"
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	network "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/client"
)

var _ pollers.PollerType = &virtualNetworkAndSubnetProvisioningSucceededPoller{}

type virtualNetworkAndSubnetProvisioningSucceededPoller struct {
	client *network.Client
	id     *commonids.SubnetId
}

func NewVirtualNetworkAndSubnetProvisioningSucceededPoller(client *network.Client, id *commonids.SubnetId) pollers.PollerType {
	return &virtualNetworkAndSubnetProvisioningSucceededPoller{
		client: client,
		id:     id,
	}
}

func (v virtualNetworkAndSubnetProvisioningSucceededPoller) Poll(ctx context.Context) (*pollers.PollResult, error) {
	pr := &pollers.PollResult{
		PollInterval: 10 * time.Second,
		Status:       pollers.PollingStatusInProgress,
	}

	subnetResp, err := v.client.Subnets.Get(ctx, *v.id, subnets.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %w", v.id, err)
	}

	if subnetResp.Model == nil || subnetResp.Model.Properties == nil || subnetResp.Model.Properties.ProvisioningState == nil {
		return nil, fmt.Errorf("retrieving %s: unable to determine provisioningState", v.id)
	}
	subnetDone := *subnetResp.Model.Properties.ProvisioningState == subnets.ProvisioningStateSucceeded

	vnetID := commonids.NewVirtualNetworkID(v.id.SubscriptionId, v.id.ResourceGroupName, v.id.VirtualNetworkName)
	vnetResp, err := v.client.VirtualNetworks.Get(ctx, vnetID, virtualnetworks.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %w", v.id, err)
	}

	if vnetResp.Model == nil || vnetResp.Model.Properties == nil || vnetResp.Model.Properties.ProvisioningState == nil {
		return nil, fmt.Errorf("retrieving %s: unable to determine provisioningState", vnetID)
	}
	vnetDone := *vnetResp.Model.Properties.ProvisioningState == virtualnetworks.ProvisioningStateSucceeded

	if subnetDone && vnetDone {
		pr.Status = pollers.PollingStatusSucceeded
		return pr, nil
	}

	return pr, nil
}
