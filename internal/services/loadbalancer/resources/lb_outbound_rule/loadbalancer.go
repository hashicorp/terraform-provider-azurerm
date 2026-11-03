// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package lb_outbound_rule

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
)

func FindLoadBalancerOutboundRuleByName(lb *loadbalancers.LoadBalancer, name string) (*loadbalancers.OutboundRule, int, bool) {
	if lb == nil || lb.Properties == nil || lb.Properties.OutboundRules == nil {
		return nil, -1, false
	}

	for i, or := range *lb.Properties.OutboundRules {
		if or.Name != nil && *or.Name == name {
			return &or, i, true
		}
	}

	return nil, -1, false
}
