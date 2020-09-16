package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LoadBalancerBackendAddressPoolId struct {
	ResourceGroup    string
	LoadBalancerName string
	Name             string
}

func (id LoadBalancerBackendAddressPoolId) ID(subscriptionId string) string {
	baseId := NewLoadBalancerID(id.ResourceGroup, id.LoadBalancerName).ID(subscriptionId)
	return fmt.Sprintf("%s/backendAddressPools/%s", baseId, id.Name)
}

func NewLoadBalancerBackendAddressPoolId(loadBalancerId LoadBalancerId, name string) LoadBalancerBackendAddressPoolId {
	return LoadBalancerBackendAddressPoolId{
		ResourceGroup:    loadBalancerId.ResourceGroup,
		LoadBalancerName: loadBalancerId.Name,
		Name:             name,
	}
}

func LoadBalancerBackendAddressPoolID(input string) (*LoadBalancerBackendAddressPoolId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Load Balancer Backend Address Pool ID %q: %+v", input, err)
	}

	backendAddressPoolId := LoadBalancerBackendAddressPoolId{
		ResourceGroup: id.ResourceGroup,
	}

	if backendAddressPoolId.LoadBalancerName, err = id.PopSegment("loadBalancers"); err != nil {
		return nil, err
	}

	if backendAddressPoolId.Name, err = id.PopSegment("backendAddressPools"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &backendAddressPoolId, nil
}
