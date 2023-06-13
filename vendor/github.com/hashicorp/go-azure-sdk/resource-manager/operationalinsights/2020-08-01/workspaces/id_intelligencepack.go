package workspaces

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = IntelligencePackId{}

// IntelligencePackId is a struct representing the Resource ID for a Intelligence Pack
type IntelligencePackId struct {
	SubscriptionId       string
	ResourceGroupName    string
	WorkspaceName        string
	IntelligencePackName string
}

// NewIntelligencePackID returns a new IntelligencePackId struct
func NewIntelligencePackID(subscriptionId string, resourceGroupName string, workspaceName string, intelligencePackName string) IntelligencePackId {
	return IntelligencePackId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		WorkspaceName:        workspaceName,
		IntelligencePackName: intelligencePackName,
	}
}

// ParseIntelligencePackID parses 'input' into a IntelligencePackId
func ParseIntelligencePackID(input string) (*IntelligencePackId, error) {
	parser := resourceids.NewParserFromResourceIdType(IntelligencePackId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := IntelligencePackId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.IntelligencePackName, ok = parsed.Parsed["intelligencePackName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "intelligencePackName", *parsed)
	}

	return &id, nil
}

// ParseIntelligencePackIDInsensitively parses 'input' case-insensitively into a IntelligencePackId
// note: this method should only be used for API response data and not user input
func ParseIntelligencePackIDInsensitively(input string) (*IntelligencePackId, error) {
	parser := resourceids.NewParserFromResourceIdType(IntelligencePackId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := IntelligencePackId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "workspaceName", *parsed)
	}

	if id.IntelligencePackName, ok = parsed.Parsed["intelligencePackName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "intelligencePackName", *parsed)
	}

	return &id, nil
}

// ValidateIntelligencePackID checks that 'input' can be parsed as a Intelligence Pack ID
func ValidateIntelligencePackID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseIntelligencePackID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Intelligence Pack ID
func (id IntelligencePackId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/intelligencePacks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, id.IntelligencePackName)
}

// Segments returns a slice of Resource ID Segments which comprise this Intelligence Pack ID
func (id IntelligencePackId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftOperationalInsights", "Microsoft.OperationalInsights", "Microsoft.OperationalInsights"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceValue"),
		resourceids.StaticSegment("staticIntelligencePacks", "intelligencePacks", "intelligencePacks"),
		resourceids.UserSpecifiedSegment("intelligencePackName", "intelligencePackValue"),
	}
}

// String returns a human-readable description of this Intelligence Pack ID
func (id IntelligencePackId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workspace Name: %q", id.WorkspaceName),
		fmt.Sprintf("Intelligence Pack Name: %q", id.IntelligencePackName),
	}
	return fmt.Sprintf("Intelligence Pack (%s)", strings.Join(components, "\n"))
}
