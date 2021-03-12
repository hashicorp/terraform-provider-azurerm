package loadbalancer

import (
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loadbalancer/parse"
)

// TODO: refactor this

func FindLoadBalancerBackEndAddressPoolByName(lb *network.LoadBalancer, name string) (*network.BackendAddressPool, int, bool) {
	if lb == nil || lb.LoadBalancerPropertiesFormat == nil || lb.LoadBalancerPropertiesFormat.BackendAddressPools == nil {
		return nil, -1, false
	}

	for i, apc := range *lb.LoadBalancerPropertiesFormat.BackendAddressPools {
		if apc.Name != nil && *apc.Name == name {
			return &apc, i, true
		}
	}

	return nil, -1, false
}

func FindLoadBalancerFrontEndIpConfigurationByName(lb *network.LoadBalancer, name string) (*network.FrontendIPConfiguration, bool) {
	if lb == nil || lb.LoadBalancerPropertiesFormat == nil || lb.LoadBalancerPropertiesFormat.FrontendIPConfigurations == nil {
		return nil, false
	}

	for _, feip := range *lb.LoadBalancerPropertiesFormat.FrontendIPConfigurations {
		if feip.Name != nil && *feip.Name == name {
			return &feip, true
		}
	}

	return nil, false
}

func FindLoadBalancerRuleByName(lb *network.LoadBalancer, name string) (*network.LoadBalancingRule, int, bool) {
	if lb == nil || lb.LoadBalancerPropertiesFormat == nil || lb.LoadBalancerPropertiesFormat.LoadBalancingRules == nil {
		return nil, -1, false
	}

	for i, lbr := range *lb.LoadBalancerPropertiesFormat.LoadBalancingRules {
		if lbr.Name != nil && *lbr.Name == name {
			return &lbr, i, true
		}
	}

	return nil, -1, false
}

func FindLoadBalancerOutboundRuleByName(lb *network.LoadBalancer, name string) (*network.OutboundRule, int, bool) {
	if lb == nil || lb.LoadBalancerPropertiesFormat == nil || lb.LoadBalancerPropertiesFormat.OutboundRules == nil {
		return nil, -1, false
	}

	for i, or := range *lb.LoadBalancerPropertiesFormat.OutboundRules {
		if or.Name != nil && *or.Name == name {
			return &or, i, true
		}
	}

	return nil, -1, false
}

func FindLoadBalancerNatRuleByName(lb *network.LoadBalancer, name string) (*network.InboundNatRule, int, bool) {
	if lb == nil || lb.LoadBalancerPropertiesFormat == nil || lb.LoadBalancerPropertiesFormat.InboundNatRules == nil {
		return nil, -1, false
	}

	for i, nr := range *lb.LoadBalancerPropertiesFormat.InboundNatRules {
		if nr.Name != nil && *nr.Name == name {
			return &nr, i, true
		}
	}

	return nil, -1, false
}

func FindLoadBalancerNatPoolByName(lb *network.LoadBalancer, name string) (*network.InboundNatPool, int, bool) {
	if lb == nil || lb.LoadBalancerPropertiesFormat == nil || lb.LoadBalancerPropertiesFormat.InboundNatPools == nil {
		return nil, -1, false
	}

	for i, np := range *lb.LoadBalancerPropertiesFormat.InboundNatPools {
		if np.Name != nil && *np.Name == name {
			return &np, i, true
		}
	}

	return nil, -1, false
}

func FindLoadBalancerProbeByName(lb *network.LoadBalancer, name string) (*network.Probe, int, bool) {
	if lb == nil || lb.LoadBalancerPropertiesFormat == nil || lb.LoadBalancerPropertiesFormat.Probes == nil {
		return nil, -1, false
	}

	for i, p := range *lb.LoadBalancerPropertiesFormat.Probes {
		if p.Name != nil && *p.Name == name {
			return &p, i, true
		}
	}

	return nil, -1, false
}

func loadBalancerSubResourceImporter(parser func(input string) (*parse.LoadBalancerId, error)) *schema.ResourceImporter {
	return &schema.ResourceImporter{
		State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
			lbId, err := parser(d.Id())
			if err != nil {
				return nil, err
			}

			d.Set("loadbalancer_id", lbId.ID())
			return []*schema.ResourceData{d}, nil
		},
	}
}
