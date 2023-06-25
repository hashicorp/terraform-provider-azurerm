package sentinelonboardingstates

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = OnboardingStateId{}

// OnboardingStateId is a struct representing the Resource ID for a Onboarding State
type OnboardingStateId struct {
	SubscriptionId      string
	ResourceGroupName   string
	WorkspaceName       string
	OnboardingStateName string
}

// NewOnboardingStateID returns a new OnboardingStateId struct
func NewOnboardingStateID(subscriptionId string, resourceGroupName string, workspaceName string, onboardingStateName string) OnboardingStateId {
	return OnboardingStateId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		WorkspaceName:       workspaceName,
		OnboardingStateName: onboardingStateName,
	}
}

// ParseOnboardingStateID parses 'input' into a OnboardingStateId
func ParseOnboardingStateID(input string) (*OnboardingStateId, error) {
	parser := resourceids.NewParserFromResourceIdType(OnboardingStateId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := OnboardingStateId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.OnboardingStateName, ok = parsed.Parsed["onboardingStateName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "onboardingStateName", *parsed)
	}

	return &id, nil
}

// ParseOnboardingStateIDInsensitively parses 'input' case-insensitively into a OnboardingStateId
// note: this method should only be used for API response data and not user input
func ParseOnboardingStateIDInsensitively(input string) (*OnboardingStateId, error) {
	parser := resourceids.NewParserFromResourceIdType(OnboardingStateId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := OnboardingStateId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.OnboardingStateName, ok = parsed.Parsed["onboardingStateName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "onboardingStateName", *parsed)
	}

	return &id, nil
}

// ValidateOnboardingStateID checks that 'input' can be parsed as a Onboarding State ID
func ValidateOnboardingStateID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseOnboardingStateID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Onboarding State ID
func (id OnboardingStateId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/providers/Microsoft.SecurityInsights/onboardingStates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, id.OnboardingStateName)
}

// Segments returns a slice of Resource ID Segments which comprise this Onboarding State ID
func (id OnboardingStateId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftOperationalInsights", "Microsoft.OperationalInsights", "Microsoft.OperationalInsights"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceValue"),
		resourceids.StaticSegment("staticProviders2", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSecurityInsights", "Microsoft.SecurityInsights", "Microsoft.SecurityInsights"),
		resourceids.StaticSegment("staticOnboardingStates", "onboardingStates", "onboardingStates"),
		resourceids.UserSpecifiedSegment("onboardingStateName", "onboardingStateValue"),
	}
}

// String returns a human-readable description of this Onboarding State ID
func (id OnboardingStateId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workspace Name: %q", id.WorkspaceName),
		fmt.Sprintf("Onboarding State Name: %q", id.OnboardingStateName),
	}
	return fmt.Sprintf("Onboarding State (%s)", strings.Join(components, "\n"))
}
