package fleetupdatestrategies

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&UpdateStrategyId{})
}

var _ resourceids.ResourceId = &UpdateStrategyId{}

// UpdateStrategyId is a struct representing the Resource ID for a Update Strategy
type UpdateStrategyId struct {
	SubscriptionId     string
	ResourceGroupName  string
	FleetName          string
	UpdateStrategyName string
}

// NewUpdateStrategyID returns a new UpdateStrategyId struct
func NewUpdateStrategyID(subscriptionId string, resourceGroupName string, fleetName string, updateStrategyName string) UpdateStrategyId {
	return UpdateStrategyId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		FleetName:          fleetName,
		UpdateStrategyName: updateStrategyName,
	}
}

// ParseUpdateStrategyID parses 'input' into a UpdateStrategyId
func ParseUpdateStrategyID(input string) (*UpdateStrategyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&UpdateStrategyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := UpdateStrategyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseUpdateStrategyIDInsensitively parses 'input' case-insensitively into a UpdateStrategyId
// note: this method should only be used for API response data and not user input
func ParseUpdateStrategyIDInsensitively(input string) (*UpdateStrategyId, error) {
	parser := resourceids.NewParserFromResourceIdType(&UpdateStrategyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := UpdateStrategyId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *UpdateStrategyId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.UpdateStrategyName, ok = input.Parsed["updateStrategyName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "updateStrategyName", input)
	}

	return nil
}

// ValidateUpdateStrategyID checks that 'input' can be parsed as a Update Strategy ID
func ValidateUpdateStrategyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseUpdateStrategyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Update Strategy ID
func (id UpdateStrategyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerService/fleets/%s/updateStrategies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.FleetName, id.UpdateStrategyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Update Strategy ID
func (id UpdateStrategyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerService", "Microsoft.ContainerService", "Microsoft.ContainerService"),
		resourceids.StaticSegment("staticFleets", "fleets", "fleets"),
		resourceids.UserSpecifiedSegment("fleetName", "fleetName"),
		resourceids.StaticSegment("staticUpdateStrategies", "updateStrategies", "updateStrategies"),
		resourceids.UserSpecifiedSegment("updateStrategyName", "updateStrategyName"),
	}
}

// String returns a human-readable description of this Update Strategy ID
func (id UpdateStrategyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Fleet Name: %q", id.FleetName),
		fmt.Sprintf("Update Strategy Name: %q", id.UpdateStrategyName),
	}
	return fmt.Sprintf("Update Strategy (%s)", strings.Join(components, "\n"))
}
