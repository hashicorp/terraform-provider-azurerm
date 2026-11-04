// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package lb_probe

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
)

func FindLoadBalancerProbeByName(lb *loadbalancers.LoadBalancer, name string) (*loadbalancers.Probe, int, bool) {
	if lb == nil || lb.Properties == nil || lb.Properties.Probes == nil {
		return nil, -1, false
	}

	for i, p := range *lb.Properties.Probes {
		if p.Name != nil && *p.Name == name {
			return &p, i, true
		}
	}

	return nil, -1, false
}
