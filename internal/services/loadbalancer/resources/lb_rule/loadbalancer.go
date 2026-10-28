// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package lb_rule

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
)

func FindLoadBalancerRuleByName(lb *loadbalancers.LoadBalancer, name string) (*loadbalancers.LoadBalancingRule, int, bool) {
	if lb == nil || lb.Properties == nil || lb.Properties.LoadBalancingRules == nil {
		return nil, -1, false
	}

	for i, lbr := range *lb.Properties.LoadBalancingRules {
		if lbr.Name != nil && *lbr.Name == name {
			return &lbr, i, true
		}
	}

	return nil, -1, false
}
