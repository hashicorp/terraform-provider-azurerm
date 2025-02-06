package profiles

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&TrafficManagerProfileId{})
}

var _ resourceids.ResourceId = &TrafficManagerProfileId{}

// TrafficManagerProfileId is a struct representing the Resource ID for a Traffic Manager Profile
type TrafficManagerProfileId struct {
	SubscriptionId            string
	ResourceGroupName         string
	TrafficManagerProfileName string
}

// NewTrafficManagerProfileID returns a new TrafficManagerProfileId struct
func NewTrafficManagerProfileID(subscriptionId string, resourceGroupName string, trafficManagerProfileName string) TrafficManagerProfileId {
	return TrafficManagerProfileId{
		SubscriptionId:            subscriptionId,
		ResourceGroupName:         resourceGroupName,
		TrafficManagerProfileName: trafficManagerProfileName,
	}
}

// ParseTrafficManagerProfileID parses 'input' into a TrafficManagerProfileId
func ParseTrafficManagerProfileID(input string) (*TrafficManagerProfileId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TrafficManagerProfileId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TrafficManagerProfileId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseTrafficManagerProfileIDInsensitively parses 'input' case-insensitively into a TrafficManagerProfileId
// note: this method should only be used for API response data and not user input
func ParseTrafficManagerProfileIDInsensitively(input string) (*TrafficManagerProfileId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TrafficManagerProfileId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TrafficManagerProfileId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *TrafficManagerProfileId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.TrafficManagerProfileName, ok = input.Parsed["trafficManagerProfileName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "trafficManagerProfileName", input)
	}

	return nil
}

// ValidateTrafficManagerProfileID checks that 'input' can be parsed as a Traffic Manager Profile ID
func ValidateTrafficManagerProfileID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTrafficManagerProfileID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Traffic Manager Profile ID
func (id TrafficManagerProfileId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/trafficManagerProfiles/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.TrafficManagerProfileName)
}

// Segments returns a slice of Resource ID Segments which comprise this Traffic Manager Profile ID
func (id TrafficManagerProfileId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticTrafficManagerProfiles", "trafficManagerProfiles", "trafficManagerProfiles"),
		resourceids.UserSpecifiedSegment("trafficManagerProfileName", "trafficManagerProfileName"),
	}
}

// String returns a human-readable description of this Traffic Manager Profile ID
func (id TrafficManagerProfileId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Traffic Manager Profile Name: %q", id.TrafficManagerProfileName),
	}
	return fmt.Sprintf("Traffic Manager Profile (%s)", strings.Join(components, "\n"))
}
