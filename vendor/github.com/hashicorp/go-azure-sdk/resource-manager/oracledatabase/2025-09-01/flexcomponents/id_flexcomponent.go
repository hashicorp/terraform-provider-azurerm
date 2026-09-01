package flexcomponents

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&FlexComponentId{})
}

var _ resourceids.ResourceId = &FlexComponentId{}

// FlexComponentId is a struct representing the Resource ID for a Flex Component
type FlexComponentId struct {
	SubscriptionId    string
	LocationName      string
	FlexComponentName string
}

// NewFlexComponentID returns a new FlexComponentId struct
func NewFlexComponentID(subscriptionId string, locationName string, flexComponentName string) FlexComponentId {
	return FlexComponentId{
		SubscriptionId:    subscriptionId,
		LocationName:      locationName,
		FlexComponentName: flexComponentName,
	}
}

// ParseFlexComponentID parses 'input' into a FlexComponentId
func ParseFlexComponentID(input string) (*FlexComponentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&FlexComponentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := FlexComponentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseFlexComponentIDInsensitively parses 'input' case-insensitively into a FlexComponentId
// note: this method should only be used for API response data and not user input
func ParseFlexComponentIDInsensitively(input string) (*FlexComponentId, error) {
	parser := resourceids.NewParserFromResourceIdType(&FlexComponentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := FlexComponentId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *FlexComponentId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.LocationName, ok = input.Parsed["locationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "locationName", input)
	}

	if id.FlexComponentName, ok = input.Parsed["flexComponentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "flexComponentName", input)
	}

	return nil
}

// ValidateFlexComponentID checks that 'input' can be parsed as a Flex Component ID
func ValidateFlexComponentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseFlexComponentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Flex Component ID
func (id FlexComponentId) ID() string {
	fmtString := "/subscriptions/%s/providers/Oracle.Database/locations/%s/flexComponents/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.FlexComponentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Flex Component ID
func (id FlexComponentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticOracleDatabase", "Oracle.Database", "Oracle.Database"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
		resourceids.StaticSegment("staticFlexComponents", "flexComponents", "flexComponents"),
		resourceids.UserSpecifiedSegment("flexComponentName", "flexComponentName"),
	}
}

// String returns a human-readable description of this Flex Component ID
func (id FlexComponentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Flex Component Name: %q", id.FlexComponentName),
	}
	return fmt.Sprintf("Flex Component (%s)", strings.Join(components, "\n"))
}
