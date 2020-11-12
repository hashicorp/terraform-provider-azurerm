package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LoadBalancerId struct {
	ResourceGroup string
	Name          string
}

func NewLoadBalancerID(resourceGroup, name string) LoadBalancerId {
	return LoadBalancerId{
		Name:          name,
		ResourceGroup: resourceGroup,
	}
}

func (id LoadBalancerId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s", subscriptionId, id.ResourceGroup, id.Name)
}

func LoadBalancerID(input string) (*LoadBalancerId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Load Balancer ID %q: %+v", input, err)
	}

	loadBalancer := LoadBalancerId{
		ResourceGroup: id.ResourceGroup,
	}

	if loadBalancer.Name, err = id.PopSegment("loadBalancers"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &loadBalancer, nil
}
