package capacityreservationgroups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = CapacityReservationGroupId{}

// CapacityReservationGroupId is a struct representing the Resource ID for a Capacity Reservation Group
type CapacityReservationGroupId struct {
	SubscriptionId               string
	ResourceGroupName            string
	CapacityReservationGroupName string
}

// NewCapacityReservationGroupID returns a new CapacityReservationGroupId struct
func NewCapacityReservationGroupID(subscriptionId string, resourceGroupName string, capacityReservationGroupName string) CapacityReservationGroupId {
	return CapacityReservationGroupId{
		SubscriptionId:               subscriptionId,
		ResourceGroupName:            resourceGroupName,
		CapacityReservationGroupName: capacityReservationGroupName,
	}
}

// ParseCapacityReservationGroupID parses 'input' into a CapacityReservationGroupId
func ParseCapacityReservationGroupID(input string) (*CapacityReservationGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(CapacityReservationGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CapacityReservationGroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.CapacityReservationGroupName, ok = parsed.Parsed["capacityReservationGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "capacityReservationGroupName", *parsed)
	}

	return &id, nil
}

// ParseCapacityReservationGroupIDInsensitively parses 'input' case-insensitively into a CapacityReservationGroupId
// note: this method should only be used for API response data and not user input
func ParseCapacityReservationGroupIDInsensitively(input string) (*CapacityReservationGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(CapacityReservationGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CapacityReservationGroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.CapacityReservationGroupName, ok = parsed.Parsed["capacityReservationGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "capacityReservationGroupName", *parsed)
	}

	return &id, nil
}

// ValidateCapacityReservationGroupID checks that 'input' can be parsed as a Capacity Reservation Group ID
func ValidateCapacityReservationGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCapacityReservationGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Capacity Reservation Group ID
func (id CapacityReservationGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/capacityReservationGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.CapacityReservationGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Capacity Reservation Group ID
func (id CapacityReservationGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCompute", "Microsoft.Compute", "Microsoft.Compute"),
		resourceids.StaticSegment("staticCapacityReservationGroups", "capacityReservationGroups", "capacityReservationGroups"),
		resourceids.UserSpecifiedSegment("capacityReservationGroupName", "capacityReservationGroupValue"),
	}
}

// String returns a human-readable description of this Capacity Reservation Group ID
func (id CapacityReservationGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Capacity Reservation Group Name: %q", id.CapacityReservationGroupName),
	}
	return fmt.Sprintf("Capacity Reservation Group (%s)", strings.Join(components, "\n"))
}
