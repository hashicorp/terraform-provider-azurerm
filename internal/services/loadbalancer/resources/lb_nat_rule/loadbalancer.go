// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package lb_nat_rule

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
)

func FindLoadBalancerNatRuleByName(lb *loadbalancers.LoadBalancer, name string) (*loadbalancers.InboundNatRule, int, bool) {
	if lb == nil || lb.Properties == nil || lb.Properties.InboundNatRules == nil {
		return nil, -1, false
	}

	for i, nr := range *lb.Properties.InboundNatRules {
		if nr.Name != nil && *nr.Name == name {
			return &nr, i, true
		}
	}

	return nil, -1, false
}
