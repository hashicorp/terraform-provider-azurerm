package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LoadBalancingRuleId struct {
	SubscriptionId   string
	ResourceGroup    string
	LoadBalancerName string
	Name             string
}

func NewLoadBalancingRuleID(subscriptionId, resourceGroup, loadBalancerName, name string) LoadBalancingRuleId {
	return LoadBalancingRuleId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		LoadBalancerName: loadBalancerName,
		Name:             name,
	}
}

func (id LoadBalancingRuleId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s/loadBalancingRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName, id.Name)
}

// LoadBalancingRuleID parses a LoadBalancingRule ID into an LoadBalancingRuleId struct
func LoadBalancingRuleID(input string) (*LoadBalancingRuleId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := LoadBalancingRuleId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.LoadBalancerName, err = id.PopSegment("loadBalancers"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("loadBalancingRules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
