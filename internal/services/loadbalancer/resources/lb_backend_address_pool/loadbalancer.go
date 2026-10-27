// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package lb_backend_address_pool

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
)

// TODO: refactor this

func FindLoadBalancerBackEndAddressPoolByName(lb *loadbalancers.LoadBalancer, name string) (*loadbalancers.BackendAddressPool, int, bool) {
	if lb == nil || lb.Properties == nil || lb.Properties.BackendAddressPools == nil {
		return nil, -1, false
	}

	for i, apc := range *lb.Properties.BackendAddressPools {
		if apc.Name != nil && *apc.Name == name {
			return &apc, i, true
		}
	}

	return nil, -1, false
}
