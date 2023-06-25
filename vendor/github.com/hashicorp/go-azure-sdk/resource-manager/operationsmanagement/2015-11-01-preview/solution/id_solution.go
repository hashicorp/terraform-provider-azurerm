package solution

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = SolutionId{}

// SolutionId is a struct representing the Resource ID for a Solution
type SolutionId struct {
	SubscriptionId    string
	ResourceGroupName string
	SolutionName      string
}

// NewSolutionID returns a new SolutionId struct
func NewSolutionID(subscriptionId string, resourceGroupName string, solutionName string) SolutionId {
	return SolutionId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		SolutionName:      solutionName,
	}
}

// ParseSolutionID parses 'input' into a SolutionId
func ParseSolutionID(input string) (*SolutionId, error) {
	parser := resourceids.NewParserFromResourceIdType(SolutionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SolutionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SolutionName, ok = parsed.Parsed["solutionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "solutionName", *parsed)
	}

	return &id, nil
}

// ParseSolutionIDInsensitively parses 'input' case-insensitively into a SolutionId
// note: this method should only be used for API response data and not user input
func ParseSolutionIDInsensitively(input string) (*SolutionId, error) {
	parser := resourceids.NewParserFromResourceIdType(SolutionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SolutionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.SolutionName, ok = parsed.Parsed["solutionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "solutionName", *parsed)
	}

	return &id, nil
}

// ValidateSolutionID checks that 'input' can be parsed as a Solution ID
func ValidateSolutionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSolutionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Solution ID
func (id SolutionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationsManagement/solutions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SolutionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Solution ID
func (id SolutionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftOperationsManagement", "Microsoft.OperationsManagement", "Microsoft.OperationsManagement"),
		resourceids.StaticSegment("staticSolutions", "solutions", "solutions"),
		resourceids.UserSpecifiedSegment("solutionName", "solutionValue"),
	}
}

// String returns a human-readable description of this Solution ID
func (id SolutionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Solution Name: %q", id.SolutionName),
	}
	return fmt.Sprintf("Solution (%s)", strings.Join(components, "\n"))
}
