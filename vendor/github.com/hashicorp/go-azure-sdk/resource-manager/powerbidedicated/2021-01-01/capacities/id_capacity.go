package capacities

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = CapacityId{}

// CapacityId is a struct representing the Resource ID for a Capacity
type CapacityId struct {
	SubscriptionId    string
	ResourceGroupName string
	CapacityName      string
}

// NewCapacityID returns a new CapacityId struct
func NewCapacityID(subscriptionId string, resourceGroupName string, capacityName string) CapacityId {
	return CapacityId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		CapacityName:      capacityName,
	}
}

// ParseCapacityID parses 'input' into a CapacityId
func ParseCapacityID(input string) (*CapacityId, error) {
	parser := resourceids.NewParserFromResourceIdType(CapacityId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CapacityId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.CapacityName, ok = parsed.Parsed["capacityName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "capacityName", *parsed)
	}

	return &id, nil
}

// ParseCapacityIDInsensitively parses 'input' case-insensitively into a CapacityId
// note: this method should only be used for API response data and not user input
func ParseCapacityIDInsensitively(input string) (*CapacityId, error) {
	parser := resourceids.NewParserFromResourceIdType(CapacityId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CapacityId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.CapacityName, ok = parsed.Parsed["capacityName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "capacityName", *parsed)
	}

	return &id, nil
}

// ValidateCapacityID checks that 'input' can be parsed as a Capacity ID
func ValidateCapacityID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCapacityID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Capacity ID
func (id CapacityId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.PowerBIDedicated/capacities/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.CapacityName)
}

// Segments returns a slice of Resource ID Segments which comprise this Capacity ID
func (id CapacityId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftPowerBIDedicated", "Microsoft.PowerBIDedicated", "Microsoft.PowerBIDedicated"),
		resourceids.StaticSegment("staticCapacities", "capacities", "capacities"),
		resourceids.UserSpecifiedSegment("capacityName", "capacityValue"),
	}
}

// String returns a human-readable description of this Capacity ID
func (id CapacityId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Capacity Name: %q", id.CapacityName),
	}
	return fmt.Sprintf("Capacity (%s)", strings.Join(components, "\n"))
}
