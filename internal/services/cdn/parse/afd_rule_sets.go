package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AfdRuleSetsId struct {
	SubscriptionId string
	ResourceGroup  string
	ProfileName    string
	RuleSetName    string
}

func NewAfdRuleSetsID(subscriptionId, resourceGroup, profileName, ruleSetName string) AfdRuleSetsId {
	return AfdRuleSetsId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ProfileName:    profileName,
		RuleSetName:    ruleSetName,
	}
}

func (id AfdRuleSetsId) String() string {
	segments := []string{
		fmt.Sprintf("Rule Set Name %q", id.RuleSetName),
		fmt.Sprintf("Profile Name %q", id.ProfileName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Afd Rule Sets", segmentsStr)
}

func (id AfdRuleSetsId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s/ruleSets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.RuleSetName)
}

// AfdRuleSetsID parses a AfdRuleSets ID into an AfdRuleSetsId struct
func AfdRuleSetsID(input string) (*AfdRuleSetsId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AfdRuleSetsId{
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

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
