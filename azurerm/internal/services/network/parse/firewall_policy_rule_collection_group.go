package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type FirewallPolicyRuleCollectionGroupId struct {
	ResourceGroup           string
	FirewallPolicyName      string
	RuleCollectionGroupName string
}

func (id FirewallPolicyRuleCollectionGroupId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/firewallPolicies/%s/ruleCollectionGroups/%s",
		subscriptionId, id.ResourceGroup, id.FirewallPolicyName, id.RuleCollectionGroupName)
}

func NewFirewallPolicyRuleCollectionGroupID(policyId FirewallPolicyId, name string) FirewallPolicyRuleCollectionGroupId {
	return FirewallPolicyRuleCollectionGroupId{
		ResourceGroup:           policyId.ResourceGroup,
		FirewallPolicyName:      policyId.Name,
		RuleCollectionGroupName: name,
	}
}

func FirewallPolicyRuleCollectionGroupID(input string) (*FirewallPolicyRuleCollectionGroupId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Firewall Policy ID %q: %+v", input, err)
	}

	group := FirewallPolicyRuleCollectionGroupId{
		ResourceGroup: id.ResourceGroup,
	}

	if group.FirewallPolicyName, err = id.PopSegment("firewallPolicies"); err != nil {
		return nil, err
	}

	if group.RuleCollectionGroupName, err = id.PopSegment("ruleCollectionGroups"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &group, nil
}
