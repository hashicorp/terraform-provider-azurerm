// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package lb_nat_pool

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
)

func FindLoadBalancerNatPoolByName(lb *loadbalancers.LoadBalancer, name string) (*loadbalancers.InboundNatPool, int, bool) {
	if lb == nil || lb.Properties == nil || lb.Properties.InboundNatPools == nil {
		return nil, -1, false
	}

	for i, np := range *lb.Properties.InboundNatPools {
		if np.Name != nil && *np.Name == name {
			return &np, i, true
		}
	}

	return nil, -1, false
}
