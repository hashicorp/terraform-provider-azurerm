package capacityreservations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&CapacityReservationId{})
}

var _ resourceids.ResourceId = &CapacityReservationId{}

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
	parser := resourceids.NewParserFromResourceIdType(&CapacityReservationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CapacityReservationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCapacityReservationIDInsensitively parses 'input' case-insensitively into a CapacityReservationId
// note: this method should only be used for API response data and not user input
func ParseCapacityReservationIDInsensitively(input string) (*CapacityReservationId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CapacityReservationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CapacityReservationId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CapacityReservationId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.CapacityReservationGroupName, ok = input.Parsed["capacityReservationGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "capacityReservationGroupName", input)
	}

	if id.CapacityReservationName, ok = input.Parsed["capacityReservationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "capacityReservationName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("capacityReservationGroupName", "capacityReservationGroupName"),
		resourceids.StaticSegment("staticCapacityReservations", "capacityReservations", "capacityReservations"),
		resourceids.UserSpecifiedSegment("capacityReservationName", "capacityReservationName"),
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
