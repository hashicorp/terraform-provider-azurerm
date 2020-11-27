package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LoadBalancerFrontendIpConfigurationId struct {
	ResourceGroup               string
	LoadBalancerName            string
	FrontendIPConfigurationName string
}

func NewLoadBalancerFrontendIPConfigurationId(loadBalancer LoadBalancerId, name string) LoadBalancerFrontendIpConfigurationId {
	return LoadBalancerFrontendIpConfigurationId{
		ResourceGroup:               loadBalancer.ResourceGroup,
		LoadBalancerName:            loadBalancer.Name,
		FrontendIPConfigurationName: name,
	}
}

func (id LoadBalancerFrontendIpConfigurationId) ID(subscriptionId string) string {
	baseId := NewLoadBalancerID(id.ResourceGroup, id.LoadBalancerName).ID(subscriptionId)
	return fmt.Sprintf("%s/frontendIPConfigurations/%s", baseId, id.FrontendIPConfigurationName)
}

func LoadBalancerFrontendIpConfigurationID(input string) (*LoadBalancerFrontendIpConfigurationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Load Balancer Frontend IP Configuration ID %q: %+v", input, err)
	}

	frontendIPConfigurationId := LoadBalancerFrontendIpConfigurationId{
		ResourceGroup: id.ResourceGroup,
	}

	if frontendIPConfigurationId.LoadBalancerName, err = id.PopSegment("loadBalancers"); err != nil {
		return nil, err
	}

	if frontendIPConfigurationId.FrontendIPConfigurationName, err = id.PopSegment("frontendIPConfigurations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &frontendIPConfigurationId, nil
}
