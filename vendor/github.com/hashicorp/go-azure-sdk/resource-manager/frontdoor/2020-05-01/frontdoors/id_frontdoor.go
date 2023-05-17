package frontdoors

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = FrontDoorId{}

// FrontDoorId is a struct representing the Resource ID for a Front Door
type FrontDoorId struct {
	SubscriptionId    string
	ResourceGroupName string
	FrontDoorName     string
}

// NewFrontDoorID returns a new FrontDoorId struct
func NewFrontDoorID(subscriptionId string, resourceGroupName string, frontDoorName string) FrontDoorId {
	return FrontDoorId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		FrontDoorName:     frontDoorName,
	}
}

// ParseFrontDoorID parses 'input' into a FrontDoorId
func ParseFrontDoorID(input string) (*FrontDoorId, error) {
	parser := resourceids.NewParserFromResourceIdType(FrontDoorId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FrontDoorId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.FrontDoorName, ok = parsed.Parsed["frontDoorName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "frontDoorName", *parsed)
	}

	return &id, nil
}

// ParseFrontDoorIDInsensitively parses 'input' case-insensitively into a FrontDoorId
// note: this method should only be used for API response data and not user input
func ParseFrontDoorIDInsensitively(input string) (*FrontDoorId, error) {
	parser := resourceids.NewParserFromResourceIdType(FrontDoorId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FrontDoorId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.FrontDoorName, ok = parsed.Parsed["frontDoorName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "frontDoorName", *parsed)
	}

	return &id, nil
}

// ValidateFrontDoorID checks that 'input' can be parsed as a Front Door ID
func ValidateFrontDoorID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseFrontDoorID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Front Door ID
func (id FrontDoorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/frontDoors/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.FrontDoorName)
}

// Segments returns a slice of Resource ID Segments which comprise this Front Door ID
func (id FrontDoorId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticFrontDoors", "frontDoors", "frontDoors"),
		resourceids.UserSpecifiedSegment("frontDoorName", "frontDoorValue"),
	}
}

// String returns a human-readable description of this Front Door ID
func (id FrontDoorId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Front Door Name: %q", id.FrontDoorName),
	}
	return fmt.Sprintf("Front Door (%s)", strings.Join(components, "\n"))
}
