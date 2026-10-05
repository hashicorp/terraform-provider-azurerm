package frontdoors

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&FrontDoorId{})
}

var _ resourceids.ResourceId = &FrontDoorId{}

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
	parser := resourceids.NewParserFromResourceIdType(&FrontDoorId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := FrontDoorId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseFrontDoorIDInsensitively parses 'input' case-insensitively into a FrontDoorId
// note: this method should only be used for API response data and not user input
func ParseFrontDoorIDInsensitively(input string) (*FrontDoorId, error) {
	parser := resourceids.NewParserFromResourceIdType(&FrontDoorId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := FrontDoorId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *FrontDoorId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.FrontDoorName, ok = input.Parsed["frontDoorName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "frontDoorName", input)
	}

	return nil
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
		resourceids.UserSpecifiedSegment("frontDoorName", "frontDoorName"),
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
