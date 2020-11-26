package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LoadBalancerFrontendIPConfigurationId struct {
	ResourceGroup    string
	LoadBalancerName string
	Name             string
}

func NewLoadBalancerFrontendIPConfigurationId(loadBalancer LoadBalancerId, name string) LoadBalancerFrontendIPConfigurationId {
	return LoadBalancerFrontendIPConfigurationId{
		ResourceGroup:    loadBalancer.ResourceGroup,
		LoadBalancerName: loadBalancer.Name,
		Name:             name,
	}
}

func (id LoadBalancerFrontendIPConfigurationId) ID(subscriptionId string) string {
	baseId := NewLoadBalancerID(id.ResourceGroup, id.LoadBalancerName).ID(subscriptionId)
	return fmt.Sprintf("%s/frontendIPConfigurations/%s", baseId, id.Name)
}

func LoadBalancerFrontendIPConfigurationID(input string) (*LoadBalancerFrontendIPConfigurationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Load Balancer Frontend IP Configuration ID %q: %+v", input, err)
	}

	frontendIPConfigurationId := LoadBalancerFrontendIPConfigurationId{
		ResourceGroup: id.ResourceGroup,
	}

	if frontendIPConfigurationId.LoadBalancerName, err = id.PopSegment("loadBalancers"); err != nil {
		return nil, err
	}

	if frontendIPConfigurationId.Name, err = id.PopSegment("frontendIPConfigurations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &frontendIPConfigurationId, nil
}
