package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FirewallPolicyRuleCollectionId struct {
	SubscriptionId          string
	ResourceGroup           string
	FirewallPolicyName      string
	RuleCollectionGroupName string
	RuleCollectionName      string
}

func NewFirewallPolicyRuleCollectionID(subscriptionId, resourceGroup, firewallPolicyName, ruleCollectionGroupName, ruleCollectionName string) FirewallPolicyRuleCollectionId {
	return FirewallPolicyRuleCollectionId{
		SubscriptionId:          subscriptionId,
		ResourceGroup:           resourceGroup,
		FirewallPolicyName:      firewallPolicyName,
		RuleCollectionGroupName: ruleCollectionGroupName,
		RuleCollectionName:      ruleCollectionName,
	}
}

func (id FirewallPolicyRuleCollectionId) String() string {
	segments := []string{
		fmt.Sprintf("Rule Collection Group %q", id.RuleCollectionName),
		fmt.Sprintf("Rule Collection Group Name %q", id.RuleCollectionGroupName),
		fmt.Sprintf("Firewall Policy Name %q", id.FirewallPolicyName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Firewall Policy Rule Collection", segmentsStr)
}

func (id FirewallPolicyRuleCollectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/firewallPolicies/%s/ruleCollectionGroups/%s/ruleCollections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.FirewallPolicyName, id.RuleCollectionGroupName, id.RuleCollectionName)
}

// FirewallPolicyRuleCollectionID parses a FirewallPolicyRuleCollection ID into an FirewallPolicyRuleCollectionId struct
func FirewallPolicyRuleCollectionID(input string) (*FirewallPolicyRuleCollectionId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FirewallPolicyRuleCollectionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.FirewallPolicyName, err = id.PopSegment("firewallPolicies"); err != nil {
		return nil, err
	}

	if resourceId.RuleCollectionGroupName, err = id.PopSegment("ruleCollectionGroups"); err != nil {
		return nil, err
	}

	if resourceId.RuleCollectionName, err = id.PopSegment("ruleCollections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
