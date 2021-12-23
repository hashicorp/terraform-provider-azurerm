package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AfdRuleSetRulesId struct {
	SubscriptionId string
	ResourceGroup  string
	ProfileName    string
	RuleSetName    string
	RuleName       string
}

func NewAfdRuleSetRulesID(subscriptionId, resourceGroup, profileName, ruleSetName, ruleName string) AfdRuleSetRulesId {
	return AfdRuleSetRulesId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ProfileName:    profileName,
		RuleSetName:    ruleSetName,
		RuleName:       ruleName,
	}
}

func (id AfdRuleSetRulesId) String() string {
	segments := []string{
		fmt.Sprintf("Rule Name %q", id.RuleName),
		fmt.Sprintf("Rule Set Name %q", id.RuleSetName),
		fmt.Sprintf("Profile Name %q", id.ProfileName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Afd Rule Set Rules", segmentsStr)
}

func (id AfdRuleSetRulesId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s/ruleSets/%s/rules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.RuleSetName, id.RuleName)
}

// AfdRuleSetRulesID parses a AfdRuleSetRules ID into an AfdRuleSetRulesId struct
func AfdRuleSetRulesID(input string) (*AfdRuleSetRulesId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AfdRuleSetRulesId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ProfileName, err = id.PopSegment("profiles"); err != nil {
		return nil, err
	}
	if resourceId.RuleSetName, err = id.PopSegment("ruleSets"); err != nil {
		return nil, err
	}
	if resourceId.RuleName, err = id.PopSegment("rules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
