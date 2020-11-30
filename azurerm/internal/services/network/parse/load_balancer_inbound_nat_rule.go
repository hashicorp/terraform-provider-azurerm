package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LoadBalancerInboundNatRuleId struct {
	SubscriptionId     string
	ResourceGroup      string
	LoadBalancerName   string
	InboundNatRuleName string
}

func NewLoadBalancerInboundNatRuleID(subscriptionId, resourceGroup, loadBalancerName, inboundNatRuleName string) LoadBalancerInboundNatRuleId {
	return LoadBalancerInboundNatRuleId{
		SubscriptionId:     subscriptionId,
		ResourceGroup:      resourceGroup,
		LoadBalancerName:   loadBalancerName,
		InboundNatRuleName: inboundNatRuleName,
	}
}

func (id LoadBalancerInboundNatRuleId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s/inboundNatRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName, id.InboundNatRuleName)
}

// LoadBalancerInboundNatRuleID parses a LoadBalancerInboundNatRule ID into an LoadBalancerInboundNatRuleId struct
func LoadBalancerInboundNatRuleID(input string) (*LoadBalancerInboundNatRuleId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := LoadBalancerInboundNatRuleId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.LoadBalancerName, err = id.PopSegment("loadBalancers"); err != nil {
		return nil, err
	}
	if resourceId.InboundNatRuleName, err = id.PopSegment("inboundNatRules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
