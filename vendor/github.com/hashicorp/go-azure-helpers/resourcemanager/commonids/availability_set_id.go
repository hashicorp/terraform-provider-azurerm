// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = &AvailabilitySetId{}

// AvailabilitySetId is a struct representing the Resource ID for a Availability Set
type AvailabilitySetId struct {
	SubscriptionId      string
	ResourceGroupName   string
	AvailabilitySetName string
}

// NewAvailabilitySetID returns a new AvailabilitySetId struct
func NewAvailabilitySetID(subscriptionId string, resourceGroupName string, availabilitySetName string) AvailabilitySetId {
	return AvailabilitySetId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		AvailabilitySetName: availabilitySetName,
	}
}

// ParseAvailabilitySetID parses 'input' into a AvailabilitySetId
func ParseAvailabilitySetID(input string) (*AvailabilitySetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AvailabilitySetId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AvailabilitySetId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAvailabilitySetIDInsensitively parses 'input' case-insensitively into a AvailabilitySetId
// note: this method should only be used for API response data and not user input
func ParseAvailabilitySetIDInsensitively(input string) (*AvailabilitySetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AvailabilitySetId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AvailabilitySetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AvailabilitySetId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.AvailabilitySetName, ok = input.Parsed["availabilitySetName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "availabilitySetName", input)
	}

	return nil
}

// ValidateAvailabilitySetID checks that 'input' can be parsed as a Availability Set ID
func ValidateAvailabilitySetID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAvailabilitySetID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Availability Set ID
func (id AvailabilitySetId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/availabilitySets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AvailabilitySetName)
}

// Segments returns a slice of Resource ID Segments which comprise this Availability Set ID
func (id AvailabilitySetId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCompute", "Microsoft.Compute", "Microsoft.Compute"),
		resourceids.StaticSegment("staticAvailabilitySets", "availabilitySets", "availabilitySets"),
		resourceids.UserSpecifiedSegment("availabilitySetName", "availabilitySetValue"),
	}
}

// String returns a human-readable description of this Availability Set ID
func (id AvailabilitySetId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Availability Set Name: %q", id.AvailabilitySetName),
	}
	return fmt.Sprintf("Availability Set (%s)", strings.Join(components, "\n"))
}
