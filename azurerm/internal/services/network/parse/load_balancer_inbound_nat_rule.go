package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LoadBalancerInboundNATRuleId struct {
	ResourceGroup    string
	LoadBalancerName string
	Name             string
}

func LoadBalancerInboundNATRuleID(input string) (*LoadBalancerInboundNATRuleId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Load Balancer Inbound NAT Rule ID %q: %+v", input, err)
	}

	natRuleId := LoadBalancerInboundNATRuleId{
		ResourceGroup: id.ResourceGroup,
	}

	if natRuleId.LoadBalancerName, err = id.PopSegment("loadBalancers"); err != nil {
		return nil, err
	}

	if natRuleId.Name, err = id.PopSegment("inboundNatRules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &natRuleId, nil
}
