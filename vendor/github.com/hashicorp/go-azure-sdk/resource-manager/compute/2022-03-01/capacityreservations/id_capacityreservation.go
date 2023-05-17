package capacityreservations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = CapacityReservationId{}

// CapacityReservationId is a struct representing the Resource ID for a Capacity Reservation
type CapacityReservationId struct {
	SubscriptionId               string
	ResourceGroupName            string
	CapacityReservationGroupName string
	CapacityReservationName      string
}

// NewCapacityReservationID returns a new CapacityReservationId struct
func NewCapacityReservationID(subscriptionId string, resourceGroupName string, capacityReservationGroupName string, capacityReservationName string) CapacityReservationId {
	return CapacityReservationId{
		SubscriptionId:               subscriptionId,
		ResourceGroupName:            resourceGroupName,
		CapacityReservationGroupName: capacityReservationGroupName,
		CapacityReservationName:      capacityReservationName,
	}
}

// ParseCapacityReservationID parses 'input' into a CapacityReservationId
func ParseCapacityReservationID(input string) (*CapacityReservationId, error) {
	parser := resourceids.NewParserFromResourceIdType(CapacityReservationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CapacityReservationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.CapacityReservationGroupName, ok = parsed.Parsed["capacityReservationGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "capacityReservationGroupName", *parsed)
	}

	if id.CapacityReservationName, ok = parsed.Parsed["capacityReservationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "capacityReservationName", *parsed)
	}

	return &id, nil
}

// ParseCapacityReservationIDInsensitively parses 'input' case-insensitively into a CapacityReservationId
// note: this method should only be used for API response data and not user input
func ParseCapacityReservationIDInsensitively(input string) (*CapacityReservationId, error) {
	parser := resourceids.NewParserFromResourceIdType(CapacityReservationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CapacityReservationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.CapacityReservationGroupName, ok = parsed.Parsed["capacityReservationGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "capacityReservationGroupName", *parsed)
	}

	if id.CapacityReservationName, ok = parsed.Parsed["capacityReservationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "capacityReservationName", *parsed)
	}

	return &id, nil
}

// ValidateCapacityReservationID checks that 'input' can be parsed as a Capacity Reservation ID
func ValidateCapacityReservationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCapacityReservationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Capacity Reservation ID
func (id CapacityReservationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/capacityReservationGroups/%s/capacityReservations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.CapacityReservationGroupName, id.CapacityReservationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Capacity Reservation ID
func (id CapacityReservationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCompute", "Microsoft.Compute", "Microsoft.Compute"),
		resourceids.StaticSegment("staticCapacityReservationGroups", "capacityReservationGroups", "capacityReservationGroups"),
		resourceids.UserSpecifiedSegment("capacityReservationGroupName", "capacityReservationGroupValue"),
		resourceids.StaticSegment("staticCapacityReservations", "capacityReservations", "capacityReservations"),
		resourceids.UserSpecifiedSegment("capacityReservationName", "capacityReservationValue"),
	}
}

// String returns a human-readable description of this Capacity Reservation ID
func (id CapacityReservationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Capacity Reservation Group Name: %q", id.CapacityReservationGroupName),
		fmt.Sprintf("Capacity Reservation Name: %q", id.CapacityReservationName),
	}
	return fmt.Sprintf("Capacity Reservation (%s)", strings.Join(components, "\n"))
}
