package virtualmachineimages

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&EdgeZoneId{})
}

var _ resourceids.ResourceId = &EdgeZoneId{}

// EdgeZoneId is a struct representing the Resource ID for a Edge Zone
type EdgeZoneId struct {
	SubscriptionId string
	LocationName   string
	EdgeZoneName   string
}

// NewEdgeZoneID returns a new EdgeZoneId struct
func NewEdgeZoneID(subscriptionId string, locationName string, edgeZoneName string) EdgeZoneId {
	return EdgeZoneId{
		SubscriptionId: subscriptionId,
		LocationName:   locationName,
		EdgeZoneName:   edgeZoneName,
	}
}

// ParseEdgeZoneID parses 'input' into a EdgeZoneId
func ParseEdgeZoneID(input string) (*EdgeZoneId, error) {
	parser := resourceids.NewParserFromResourceIdType(&EdgeZoneId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := EdgeZoneId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseEdgeZoneIDInsensitively parses 'input' case-insensitively into a EdgeZoneId
// note: this method should only be used for API response data and not user input
func ParseEdgeZoneIDInsensitively(input string) (*EdgeZoneId, error) {
	parser := resourceids.NewParserFromResourceIdType(&EdgeZoneId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := EdgeZoneId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *EdgeZoneId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.LocationName, ok = input.Parsed["locationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "locationName", input)
	}

	if id.EdgeZoneName, ok = input.Parsed["edgeZoneName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "edgeZoneName", input)
	}

	return nil
}

// ValidateEdgeZoneID checks that 'input' can be parsed as a Edge Zone ID
func ValidateEdgeZoneID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseEdgeZoneID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Edge Zone ID
func (id EdgeZoneId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Compute/locations/%s/edgeZones/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.EdgeZoneName)
}

// Segments returns a slice of Resource ID Segments which comprise this Edge Zone ID
func (id EdgeZoneId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCompute", "Microsoft.Compute", "Microsoft.Compute"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
		resourceids.StaticSegment("staticEdgeZones", "edgeZones", "edgeZones"),
		resourceids.UserSpecifiedSegment("edgeZoneName", "edgeZoneName"),
	}
}

// String returns a human-readable description of this Edge Zone ID
func (id EdgeZoneId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Edge Zone Name: %q", id.EdgeZoneName),
	}
	return fmt.Sprintf("Edge Zone (%s)", strings.Join(components, "\n"))
}
