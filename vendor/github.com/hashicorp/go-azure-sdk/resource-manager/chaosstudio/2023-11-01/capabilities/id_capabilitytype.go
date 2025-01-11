package capabilities

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&CapabilityTypeId{})
}

var _ resourceids.ResourceId = &CapabilityTypeId{}

// CapabilityTypeId is a struct representing the Resource ID for a Capability Type
type CapabilityTypeId struct {
	SubscriptionId     string
	LocationName       string
	TargetTypeName     string
	CapabilityTypeName string
}

// NewCapabilityTypeID returns a new CapabilityTypeId struct
func NewCapabilityTypeID(subscriptionId string, locationName string, targetTypeName string, capabilityTypeName string) CapabilityTypeId {
	return CapabilityTypeId{
		SubscriptionId:     subscriptionId,
		LocationName:       locationName,
		TargetTypeName:     targetTypeName,
		CapabilityTypeName: capabilityTypeName,
	}
}

// ParseCapabilityTypeID parses 'input' into a CapabilityTypeId
func ParseCapabilityTypeID(input string) (*CapabilityTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CapabilityTypeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CapabilityTypeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseCapabilityTypeIDInsensitively parses 'input' case-insensitively into a CapabilityTypeId
// note: this method should only be used for API response data and not user input
func ParseCapabilityTypeIDInsensitively(input string) (*CapabilityTypeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&CapabilityTypeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := CapabilityTypeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *CapabilityTypeId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.LocationName, ok = input.Parsed["locationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "locationName", input)
	}

	if id.TargetTypeName, ok = input.Parsed["targetTypeName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "targetTypeName", input)
	}

	if id.CapabilityTypeName, ok = input.Parsed["capabilityTypeName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "capabilityTypeName", input)
	}

	return nil
}

// ValidateCapabilityTypeID checks that 'input' can be parsed as a Capability Type ID
func ValidateCapabilityTypeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCapabilityTypeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Capability Type ID
func (id CapabilityTypeId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Chaos/locations/%s/targetTypes/%s/capabilityTypes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.TargetTypeName, id.CapabilityTypeName)
}

// Segments returns a slice of Resource ID Segments which comprise this Capability Type ID
func (id CapabilityTypeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftChaos", "Microsoft.Chaos", "Microsoft.Chaos"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
		resourceids.StaticSegment("staticTargetTypes", "targetTypes", "targetTypes"),
		resourceids.UserSpecifiedSegment("targetTypeName", "targetTypeName"),
		resourceids.StaticSegment("staticCapabilityTypes", "capabilityTypes", "capabilityTypes"),
		resourceids.UserSpecifiedSegment("capabilityTypeName", "capabilityTypeName"),
	}
}

// String returns a human-readable description of this Capability Type ID
func (id CapabilityTypeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Target Type Name: %q", id.TargetTypeName),
		fmt.Sprintf("Capability Type Name: %q", id.CapabilityTypeName),
	}
	return fmt.Sprintf("Capability Type (%s)", strings.Join(components, "\n"))
}
