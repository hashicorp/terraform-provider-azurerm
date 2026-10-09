package fleets

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&FleetspaceId{})
}

var _ resourceids.ResourceId = &FleetspaceId{}

// FleetspaceId is a struct representing the Resource ID for a Fleetspace
type FleetspaceId struct {
	SubscriptionId    string
	ResourceGroupName string
	FleetName         string
	FleetspaceName    string
}

// NewFleetspaceID returns a new FleetspaceId struct
func NewFleetspaceID(subscriptionId string, resourceGroupName string, fleetName string, fleetspaceName string) FleetspaceId {
	return FleetspaceId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		FleetName:         fleetName,
		FleetspaceName:    fleetspaceName,
	}
}

// ParseFleetspaceID parses 'input' into a FleetspaceId
func ParseFleetspaceID(input string) (*FleetspaceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&FleetspaceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := FleetspaceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseFleetspaceIDInsensitively parses 'input' case-insensitively into a FleetspaceId
// note: this method should only be used for API response data and not user input
func ParseFleetspaceIDInsensitively(input string) (*FleetspaceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&FleetspaceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := FleetspaceId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *FleetspaceId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.FleetspaceName, ok = input.Parsed["fleetspaceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "fleetspaceName", input)
	}

	return nil
}

// ValidateFleetspaceID checks that 'input' can be parsed as a Fleetspace ID
func ValidateFleetspaceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseFleetspaceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Fleetspace ID
func (id FleetspaceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/fleets/%s/fleetspaces/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.FleetName, id.FleetspaceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Fleetspace ID
func (id FleetspaceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticFleets", "fleets", "fleets"),
		resourceids.UserSpecifiedSegment("fleetName", "fleetName"),
		resourceids.StaticSegment("staticFleetspaces", "fleetspaces", "fleetspaces"),
		resourceids.UserSpecifiedSegment("fleetspaceName", "fleetspaceName"),
	}
}

// String returns a human-readable description of this Fleetspace ID
func (id FleetspaceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Fleet Name: %q", id.FleetName),
		fmt.Sprintf("Fleetspace Name: %q", id.FleetspaceName),
	}
	return fmt.Sprintf("Fleetspace (%s)", strings.Join(components, "\n"))
}
