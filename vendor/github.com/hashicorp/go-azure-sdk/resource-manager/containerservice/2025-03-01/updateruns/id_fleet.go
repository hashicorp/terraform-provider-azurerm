package updateruns

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&FleetId{})
}

var _ resourceids.ResourceId = &FleetId{}

// FleetId is a struct representing the Resource ID for a Fleet
type FleetId struct {
	SubscriptionId    string
	ResourceGroupName string
	FleetName         string
}

// NewFleetID returns a new FleetId struct
func NewFleetID(subscriptionId string, resourceGroupName string, fleetName string) FleetId {
	return FleetId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		FleetName:         fleetName,
	}
}

// ParseFleetID parses 'input' into a FleetId
func ParseFleetID(input string) (*FleetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&FleetId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := FleetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseFleetIDInsensitively parses 'input' case-insensitively into a FleetId
// note: this method should only be used for API response data and not user input
func ParseFleetIDInsensitively(input string) (*FleetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&FleetId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := FleetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *FleetId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.FleetName, ok = input.Parsed["fleetName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "fleetName", input)
	}

	return nil
}

// ValidateFleetID checks that 'input' can be parsed as a Fleet ID
func ValidateFleetID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseFleetID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Fleet ID
func (id FleetId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerService/fleets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.FleetName)
}

// Segments returns a slice of Resource ID Segments which comprise this Fleet ID
func (id FleetId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerService", "Microsoft.ContainerService", "Microsoft.ContainerService"),
		resourceids.StaticSegment("staticFleets", "fleets", "fleets"),
		resourceids.UserSpecifiedSegment("fleetName", "fleetName"),
	}
}

// String returns a human-readable description of this Fleet ID
func (id FleetId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Fleet Name: %q", id.FleetName),
	}
	return fmt.Sprintf("Fleet (%s)", strings.Join(components, "\n"))
}
