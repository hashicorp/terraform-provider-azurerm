// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func FlattenVirtualNetworkDNSServers(input *virtualnetworks.DhcpOptions) []string {
	results := make([]string, 0)

	if input != nil {
		if servers := input.DnsServers; servers != nil {
			results = *servers
		}
	}

	return results
}

func VirtualNetworkProvisioningStateRefreshFunc(ctx context.Context, client *virtualnetworks.VirtualNetworksClient, id commonids.VirtualNetworkId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		res, err := client.Get(ctx, id, virtualnetworks.DefaultGetOperationOptions())
		if err != nil {
			return nil, "", fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if res.Model != nil && res.Model.Properties != nil {
			return res, pointer.FromEnum(res.Model.Properties.ProvisioningState), nil
		}
		return res, "", fmt.Errorf("polling for %s: %+v", id, err)
	}
}
