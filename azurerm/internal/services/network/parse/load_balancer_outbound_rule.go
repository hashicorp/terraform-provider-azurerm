package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LoadBalancerOutboundRuleId struct {
	SubscriptionId   string
	ResourceGroup    string
	LoadBalancerName string
	OutboundRuleName string
}

func NewLoadBalancerOutboundRuleID(subscriptionId, resourceGroup, loadBalancerName, outboundRuleName string) LoadBalancerOutboundRuleId {
	return LoadBalancerOutboundRuleId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		LoadBalancerName: loadBalancerName,
		OutboundRuleName: outboundRuleName,
	}
}

func (id LoadBalancerOutboundRuleId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s/outboundRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName, id.OutboundRuleName)
}

// LoadBalancerOutboundRuleID parses a LoadBalancerOutboundRule ID into an LoadBalancerOutboundRuleId struct
func LoadBalancerOutboundRuleID(input string) (*LoadBalancerOutboundRuleId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := LoadBalancerOutboundRuleId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.LoadBalancerName, err = id.PopSegment("loadBalancers"); err != nil {
		return nil, err
	}
	if resourceId.OutboundRuleName, err = id.PopSegment("outboundRules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
