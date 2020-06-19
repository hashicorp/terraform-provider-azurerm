package parse

import (
	"fmt"
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type FirewallPolicyId struct {
	ResourceGroup string
	Name          string
}

func FirewallPolicyID(input string) (*FirewallPolicyId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing Firewall Policy ID %q: %+v", input, err)
	}

	policy := FirewallPolicyId{
		ResourceGroup: id.ResourceGroup,
	}

	if policy.Name, err = id.PopSegment("firewallPolicies"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &policy, nil
}

type FirewallPolicyRuleCollectionGroupId struct {
	ResourceGroup string
	PolicyName    string
	Name          string
}

func FirewallPolicyRuleCollectionGroupID(input string) (*FirewallPolicyRuleCollectionGroupId, error) {
	groups := regexp.MustCompile(`^(.+)/ruleCollectionGroups/(.+)$`).FindStringSubmatch(input)
	if len(groups) != 3 {
		return nil, fmt.Errorf("failed to parse Firewall Policy Rule Collection Group ID: %q", input)
	}
	policy, name := groups[1], groups[2]

	policyId, err := FirewallPolicyID(policy)
	if err != nil {
		return nil, err
	}

	return &FirewallPolicyRuleCollectionGroupId{
		ResourceGroup: policyId.ResourceGroup,
		PolicyName:    policyId.Name,
		Name:          name,
	}, nil
}
