package proximityplacementgroups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ProximityPlacementGroupId{})
}

var _ resourceids.ResourceId = &ProximityPlacementGroupId{}

// ProximityPlacementGroupId is a struct representing the Resource ID for a Proximity Placement Group
type ProximityPlacementGroupId struct {
	SubscriptionId              string
	ResourceGroupName           string
	ProximityPlacementGroupName string
}

// NewProximityPlacementGroupID returns a new ProximityPlacementGroupId struct
func NewProximityPlacementGroupID(subscriptionId string, resourceGroupName string, proximityPlacementGroupName string) ProximityPlacementGroupId {
	return ProximityPlacementGroupId{
		SubscriptionId:              subscriptionId,
		ResourceGroupName:           resourceGroupName,
		ProximityPlacementGroupName: proximityPlacementGroupName,
	}
}

// ParseProximityPlacementGroupID parses 'input' into a ProximityPlacementGroupId
func ParseProximityPlacementGroupID(input string) (*ProximityPlacementGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProximityPlacementGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProximityPlacementGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseProximityPlacementGroupIDInsensitively parses 'input' case-insensitively into a ProximityPlacementGroupId
// note: this method should only be used for API response data and not user input
func ParseProximityPlacementGroupIDInsensitively(input string) (*ProximityPlacementGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ProximityPlacementGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ProximityPlacementGroupId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ProximityPlacementGroupId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.ProximityPlacementGroupName, ok = input.Parsed["proximityPlacementGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "proximityPlacementGroupName", input)
	}

	return nil
}

// ValidateProximityPlacementGroupID checks that 'input' can be parsed as a Proximity Placement Group ID
func ValidateProximityPlacementGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProximityPlacementGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Proximity Placement Group ID
func (id ProximityPlacementGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/proximityPlacementGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProximityPlacementGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Proximity Placement Group ID
func (id ProximityPlacementGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCompute", "Microsoft.Compute", "Microsoft.Compute"),
		resourceids.StaticSegment("staticProximityPlacementGroups", "proximityPlacementGroups", "proximityPlacementGroups"),
		resourceids.UserSpecifiedSegment("proximityPlacementGroupName", "proximityPlacementGroupName"),
	}
}

// String returns a human-readable description of this Proximity Placement Group ID
func (id ProximityPlacementGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Proximity Placement Group Name: %q", id.ProximityPlacementGroupName),
	}
	return fmt.Sprintf("Proximity Placement Group (%s)", strings.Join(components, "\n"))
}
