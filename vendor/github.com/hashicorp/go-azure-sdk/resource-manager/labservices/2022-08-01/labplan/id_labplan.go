package labplan

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = LabPlanId{}

// LabPlanId is a struct representing the Resource ID for a Lab Plan
type LabPlanId struct {
	SubscriptionId    string
	ResourceGroupName string
	LabPlanName       string
}

// NewLabPlanID returns a new LabPlanId struct
func NewLabPlanID(subscriptionId string, resourceGroupName string, labPlanName string) LabPlanId {
	return LabPlanId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		LabPlanName:       labPlanName,
	}
}

// ParseLabPlanID parses 'input' into a LabPlanId
func ParseLabPlanID(input string) (*LabPlanId, error) {
	parser := resourceids.NewParserFromResourceIdType(LabPlanId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LabPlanId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LabPlanName, ok = parsed.Parsed["labPlanName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "labPlanName", *parsed)
	}

	return &id, nil
}

// ParseLabPlanIDInsensitively parses 'input' case-insensitively into a LabPlanId
// note: this method should only be used for API response data and not user input
func ParseLabPlanIDInsensitively(input string) (*LabPlanId, error) {
	parser := resourceids.NewParserFromResourceIdType(LabPlanId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LabPlanId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LabPlanName, ok = parsed.Parsed["labPlanName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "labPlanName", *parsed)
	}

	return &id, nil
}

// ValidateLabPlanID checks that 'input' can be parsed as a Lab Plan ID
func ValidateLabPlanID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLabPlanID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Lab Plan ID
func (id LabPlanId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.LabServices/labPlans/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LabPlanName)
}

// Segments returns a slice of Resource ID Segments which comprise this Lab Plan ID
func (id LabPlanId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftLabServices", "Microsoft.LabServices", "Microsoft.LabServices"),
		resourceids.StaticSegment("staticLabPlans", "labPlans", "labPlans"),
		resourceids.UserSpecifiedSegment("labPlanName", "labPlanValue"),
	}
}

// String returns a human-readable description of this Lab Plan ID
func (id LabPlanId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Lab Plan Name: %q", id.LabPlanName),
	}
	return fmt.Sprintf("Lab Plan (%s)", strings.Join(components, "\n"))
}
