package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LoadBalancerInboundNATPoolId struct {
	ResourceGroup      string
	LoadBalancerName   string
	InboundNatPoolName string
}

func LoadBalancerInboundNATPoolID(input string) (*LoadBalancerInboundNATPoolId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Load Balancer Inbound NAT Pool ID %q: %+v", input, err)
	}

	natPoolId := LoadBalancerInboundNATPoolId{
		ResourceGroup: id.ResourceGroup,
	}

	if natPoolId.LoadBalancerName, err = id.PopSegment("loadBalancers"); err != nil {
		return nil, err
	}

	if natPoolId.InboundNatPoolName, err = id.PopSegment("inboundNatPools"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &natPoolId, nil
}
