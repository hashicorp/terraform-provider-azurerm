package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type FirewallPolicyRuleCollectionGroupId struct {
	SubscriptionId          string
	ResourceGroup           string
	FirewallPolicyName      string
	RuleCollectionGroupName string
}

func NewFirewallPolicyRuleCollectionGroupID(subscriptionId, resourceGroup, firewallPolicyName, ruleCollectionGroupName string) FirewallPolicyRuleCollectionGroupId {
	return FirewallPolicyRuleCollectionGroupId{
		SubscriptionId:          subscriptionId,
		ResourceGroup:           resourceGroup,
		FirewallPolicyName:      firewallPolicyName,
		RuleCollectionGroupName: ruleCollectionGroupName,
	}
}

func (id FirewallPolicyRuleCollectionGroupId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/firewallPolicies/%s/ruleCollectionGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.FirewallPolicyName, id.RuleCollectionGroupName)
}

// FirewallPolicyRuleCollectionGroupID parses a FirewallPolicyRuleCollectionGroup ID into an FirewallPolicyRuleCollectionGroupId struct
func FirewallPolicyRuleCollectionGroupID(input string) (*FirewallPolicyRuleCollectionGroupId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FirewallPolicyRuleCollectionGroupId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.FirewallPolicyName, err = id.PopSegment("firewallPolicies"); err != nil {
		return nil, err
	}
	if resourceId.RuleCollectionGroupName, err = id.PopSegment("ruleCollectionGroups"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
