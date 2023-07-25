// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FrontDoorRuleId struct {
	SubscriptionId string
	ResourceGroup  string
	ProfileName    string
	RuleSetName    string
	RuleName       string
}

func NewFrontDoorRuleID(subscriptionId, resourceGroup, profileName, ruleSetName, ruleName string) FrontDoorRuleId {
	return FrontDoorRuleId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		ProfileName:    profileName,
		RuleSetName:    ruleSetName,
		RuleName:       ruleName,
	}
}

func (id FrontDoorRuleId) String() string {
	segments := []string{
		fmt.Sprintf("Rule Name %q", id.RuleName),
		fmt.Sprintf("Rule Set Name %q", id.RuleSetName),
		fmt.Sprintf("Profile Name %q", id.ProfileName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Front Door Rule", segmentsStr)
}

func (id FrontDoorRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s/ruleSets/%s/rules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.RuleSetName, id.RuleName)
}

// FrontDoorRuleID parses a FrontDoorRule ID into an FrontDoorRuleId struct
func FrontDoorRuleID(input string) (*FrontDoorRuleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("parsing %q as an FrontDoorRule ID: %+v", input, err)
	}

	resourceId := FrontDoorRuleId{
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

// FrontDoorRuleIDInsensitively parses an FrontDoorRule ID into an FrontDoorRuleId struct, insensitively
// This should only be used to parse an ID for rewriting, the FrontDoorRuleID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func FrontDoorRuleIDInsensitively(input string) (*FrontDoorRuleId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontDoorRuleId{
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

	// find the correct casing for the 'rules' segment
	rulesKey := "rules"
	for key := range id.Path {
		if strings.EqualFold(key, rulesKey) {
			rulesKey = key
			break
		}
	}
	if resourceId.RuleName, err = id.PopSegment(rulesKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
