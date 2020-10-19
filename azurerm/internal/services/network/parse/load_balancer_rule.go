package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LoadBalancerRuleId struct {
	ResourceGroup    string
	LoadBalancerName string
	Name             string
}

func LoadBalancerRuleID(input string) (*LoadBalancerRuleId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Load Balancer Rule ID %q: %+v", input, err)
	}

	ruleId := LoadBalancerRuleId{
		ResourceGroup: id.ResourceGroup,
	}

	if ruleId.LoadBalancerName, err = id.PopSegment("loadBalancers"); err != nil {
		return nil, err
	}

	if ruleId.Name, err = id.PopSegment("loadBalancingRules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &ruleId, nil
}
