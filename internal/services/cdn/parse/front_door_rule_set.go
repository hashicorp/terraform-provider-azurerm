package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FrontDoorRuleSetId struct {
	SubscriptionId string
	ResourceGroup  string
	ProfileName    string
	RuleSetName    string
}

func NewFrontDoorRuleSetID(subscriptionId, resourceGroup, profileName, ruleSetName string) FrontDoorRuleSetId {
	return FrontDoorRuleSetId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ProfileName:    profileName,
		RuleSetName:    ruleSetName,
	}
}

func (id FrontDoorRuleSetId) String() string {
	segments := []string{
		fmt.Sprintf("Rule Set Name %q", id.RuleSetName),
		fmt.Sprintf("Profile Name %q", id.ProfileName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Front Door Rule Set", segmentsStr)
}

func (id FrontDoorRuleSetId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s/ruleSets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.RuleSetName)
}

// FrontDoorRuleSetID parses a FrontDoorRuleSet ID into an FrontDoorRuleSetId struct
func FrontDoorRuleSetID(input string) (*FrontDoorRuleSetId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontDoorRuleSetId{
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

// FrontDoorRuleSetIDInsensitively parses an FrontDoorRuleSet ID into an FrontDoorRuleSetId struct, insensitively
// This should only be used to parse an ID for rewriting, the FrontDoorRuleSetID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func FrontDoorRuleSetIDInsensitively(input string) (*FrontDoorRuleSetId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontDoorRuleSetId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'profiles' segment
	profilesKey := "profiles"
	for key := range id.Path {
		if strings.EqualFold(key, profilesKey) {
			profilesKey = key
			break
		}
	}
	if resourceId.ProfileName, err = id.PopSegment(profilesKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'ruleSets' segment
	ruleSetsKey := "ruleSets"
	for key := range id.Path {
		if strings.EqualFold(key, ruleSetsKey) {
			ruleSetsKey = key
			break
		}
	}
	if resourceId.RuleSetName, err = id.PopSegment(ruleSetsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
