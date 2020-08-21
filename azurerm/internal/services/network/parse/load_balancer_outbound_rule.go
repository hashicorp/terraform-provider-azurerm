package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LoadBalancerOutboundRuleId struct {
	ResourceGroup    string
	LoadBalancerName string
	Name             string
}

func LoadBalancerOutboundRuleID(input string) (*LoadBalancerOutboundRuleId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Load Balancer Outbound Rule ID %q: %+v", input, err)
	}

	outboundRuleId := LoadBalancerOutboundRuleId{
		ResourceGroup: id.ResourceGroup,
	}

	if outboundRuleId.LoadBalancerName, err = id.PopSegment("loadBalancers"); err != nil {
		return nil, err
	}

	if outboundRuleId.Name, err = id.PopSegment("outboundRules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &outboundRuleId, nil
}
