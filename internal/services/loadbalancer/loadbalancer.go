// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loadbalancer

import (
	"context"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/loadbalancers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
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

func FindLoadBalancerFrontEndIpConfigurationByName(lb *loadbalancers.LoadBalancer, name string) (*loadbalancers.FrontendIPConfiguration, bool) {
	if lb == nil || lb.Properties == nil || lb.Properties.FrontendIPConfigurations == nil {
		return nil, false
	}

	for _, feip := range *lb.Properties.FrontendIPConfigurations {
		if feip.Name != nil && *feip.Name == name {
			return &feip, true
		}
	}

	return nil, false
}

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

func loadBalancerSubResourceImporter(parser func(input string) (*loadbalancers.LoadBalancerId, error)) *schema.ResourceImporter {
	return pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
		_, err := parser(id)
		return err
	}, func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) ([]*pluginsdk.ResourceData, error) {
		lbId, err := parser(d.Id())
		if err != nil {
			return nil, err
		}

		d.Set("loadbalancer_id", lbId.ID())
		return []*pluginsdk.ResourceData{d}, nil
	})
}
