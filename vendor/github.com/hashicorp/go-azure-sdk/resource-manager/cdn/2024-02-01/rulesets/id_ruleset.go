package rulesets

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&RuleSetId{})
}

var _ resourceids.ResourceId = &RuleSetId{}

// RuleSetId is a struct representing the Resource ID for a Rule Set
type RuleSetId struct {
	SubscriptionId    string
	ResourceGroupName string
	ProfileName       string
	RuleSetName       string
}

// NewRuleSetID returns a new RuleSetId struct
func NewRuleSetID(subscriptionId string, resourceGroupName string, profileName string, ruleSetName string) RuleSetId {
	return RuleSetId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ProfileName:       profileName,
		RuleSetName:       ruleSetName,
	}
}

// ParseRuleSetID parses 'input' into a RuleSetId
func ParseRuleSetID(input string) (*RuleSetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RuleSetId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RuleSetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseRuleSetIDInsensitively parses 'input' case-insensitively into a RuleSetId
// note: this method should only be used for API response data and not user input
func ParseRuleSetIDInsensitively(input string) (*RuleSetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&RuleSetId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := RuleSetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *RuleSetId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ProfileName, ok = input.Parsed["profileName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "profileName", input)
	}

	if id.RuleSetName, ok = input.Parsed["ruleSetName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "ruleSetName", input)
	}

	return nil
}

// ValidateRuleSetID checks that 'input' can be parsed as a Rule Set ID
func ValidateRuleSetID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRuleSetID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Rule Set ID
func (id RuleSetId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s/ruleSets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProfileName, id.RuleSetName)
}

// Segments returns a slice of Resource ID Segments which comprise this Rule Set ID
func (id RuleSetId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCdn", "Microsoft.Cdn", "Microsoft.Cdn"),
		resourceids.StaticSegment("staticProfiles", "profiles", "profiles"),
		resourceids.UserSpecifiedSegment("profileName", "profileName"),
		resourceids.StaticSegment("staticRuleSets", "ruleSets", "ruleSets"),
		resourceids.UserSpecifiedSegment("ruleSetName", "ruleSetName"),
	}
}

// String returns a human-readable description of this Rule Set ID
func (id RuleSetId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Profile Name: %q", id.ProfileName),
		fmt.Sprintf("Rule Set Name: %q", id.RuleSetName),
	}
	return fmt.Sprintf("Rule Set (%s)", strings.Join(components, "\n"))
}
