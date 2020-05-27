package network

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-03-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// TODO: refactor this

func resourceGroupAndLBNameFromId(loadBalancerId string) (string, string, error) {
	id, err := azure.ParseAzureResourceID(loadBalancerId)
	if err != nil {
		return "", "", err
	}
	name := id.Path["loadBalancers"]
	resGroup := id.ResourceGroup

	return resGroup, name, nil
}

func retrieveLoadBalancerById(d *schema.ResourceData, loadBalancerId string, meta interface{}) (*network.LoadBalancer, bool, error) {
	client := meta.(*clients.Client).Network.LoadBalancersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resGroup, name, err := resourceGroupAndLBNameFromId(loadBalancerId)
	if err != nil {
		return nil, false, fmt.Errorf("Error Getting Load Balancer Name and Group:: %+v", err)
	}

	resp, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("Error making Read request on Azure Load Balancer %s: %s", name, err)
	}

	return &resp, true, nil
}

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

// sets the loadbalancer_id in the ResourceData from the sub resources full id
func loadBalancerSubResourceStateImporter(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	r := regexp.MustCompile(`.+/loadBalancers/.+?/`)

	lbID := strings.TrimSuffix(r.FindString(d.Id()), "/")
	parsed, err := azure.ParseAzureResourceID(lbID)
	if err != nil {
		return nil, fmt.Errorf("unable to parse loadbalancer id from %s", d.Id())
	}

	if parsed.Path["loadBalancers"] == "" {
		return nil, fmt.Errorf("parsed ID is invalid")
	}

	d.Set("loadbalancer_id", lbID)
	return []*schema.ResourceData{d}, nil
}
