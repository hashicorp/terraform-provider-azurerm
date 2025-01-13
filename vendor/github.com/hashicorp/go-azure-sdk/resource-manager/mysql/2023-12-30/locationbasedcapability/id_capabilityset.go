package locationbasedcapability

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&CapabilitySetId{})
}

var _ resourceids.ResourceId = &CapabilitySetId{}

// CapabilitySetId is a struct representing the Resource ID for a Capability Set
type CapabilitySetId struct {
	SubscriptionId    string
	LocationName      string
	CapabilitySetName string
}

// NewCapabilitySetID returns a new CapabilitySetId struct
func NewCapabilitySetID(subscriptionId string, locationName string, capabilitySetName string) CapabilitySetId {
	return CapabilitySetId{
		SubscriptionId:    subscriptionId,
		LocationName:      locationName,
		CapabilitySetName: capabilitySetName,
	}
}

// ParseCapabilitySetID parses 'input' into a CapabilitySetId
func ParseCapabilitySetID(input string) (*CapabilitySetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CapabilitySetId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CapabilitySetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCapabilitySetIDInsensitively parses 'input' case-insensitively into a CapabilitySetId
// note: this method should only be used for API response data and not user input
func ParseCapabilitySetIDInsensitively(input string) (*CapabilitySetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CapabilitySetId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CapabilitySetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CapabilitySetId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.LocationName, ok = input.Parsed["locationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "locationName", input)
	}

	if id.CapabilitySetName, ok = input.Parsed["capabilitySetName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "capabilitySetName", input)
	}

	return nil
}

// ValidateCapabilitySetID checks that 'input' can be parsed as a Capability Set ID
func ValidateCapabilitySetID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCapabilitySetID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Capability Set ID
func (id CapabilitySetId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.DBforMySQL/locations/%s/capabilitySets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.CapabilitySetName)
}

// Segments returns a slice of Resource ID Segments which comprise this Capability Set ID
func (id CapabilitySetId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDBforMySQL", "Microsoft.DBforMySQL", "Microsoft.DBforMySQL"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
		resourceids.StaticSegment("staticCapabilitySets", "capabilitySets", "capabilitySets"),
		resourceids.UserSpecifiedSegment("capabilitySetName", "capabilitySetName"),
	}
}

// String returns a human-readable description of this Capability Set ID
func (id CapabilitySetId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Capability Set Name: %q", id.CapabilitySetName),
	}
	return fmt.Sprintf("Capability Set (%s)", strings.Join(components, "\n"))
}
