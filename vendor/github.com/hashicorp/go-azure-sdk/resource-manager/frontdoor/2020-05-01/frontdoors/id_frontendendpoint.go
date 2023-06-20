package frontdoors

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = FrontendEndpointId{}

// FrontendEndpointId is a struct representing the Resource ID for a Frontend Endpoint
type FrontendEndpointId struct {
	SubscriptionId       string
	ResourceGroupName    string
	FrontDoorName        string
	FrontendEndpointName string
}

// NewFrontendEndpointID returns a new FrontendEndpointId struct
func NewFrontendEndpointID(subscriptionId string, resourceGroupName string, frontDoorName string, frontendEndpointName string) FrontendEndpointId {
	return FrontendEndpointId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		FrontDoorName:        frontDoorName,
		FrontendEndpointName: frontendEndpointName,
	}
}

// ParseFrontendEndpointID parses 'input' into a FrontendEndpointId
func ParseFrontendEndpointID(input string) (*FrontendEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(FrontendEndpointId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FrontendEndpointId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.FrontDoorName, ok = parsed.Parsed["frontDoorName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "frontDoorName", *parsed)
	}

	if id.FrontendEndpointName, ok = parsed.Parsed["frontendEndpointName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "frontendEndpointName", *parsed)
	}

	return &id, nil
}

// ParseFrontendEndpointIDInsensitively parses 'input' case-insensitively into a FrontendEndpointId
// note: this method should only be used for API response data and not user input
func ParseFrontendEndpointIDInsensitively(input string) (*FrontendEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(FrontendEndpointId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FrontendEndpointId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.FrontDoorName, ok = parsed.Parsed["frontDoorName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "frontDoorName", *parsed)
	}

	if id.FrontendEndpointName, ok = parsed.Parsed["frontendEndpointName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "frontendEndpointName", *parsed)
	}

	return &id, nil
}

// ValidateFrontendEndpointID checks that 'input' can be parsed as a Frontend Endpoint ID
func ValidateFrontendEndpointID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseFrontendEndpointID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Frontend Endpoint ID
func (id FrontendEndpointId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/frontDoors/%s/frontendEndpoints/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.FrontDoorName, id.FrontendEndpointName)
}

// Segments returns a slice of Resource ID Segments which comprise this Frontend Endpoint ID
func (id FrontendEndpointId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticFrontDoors", "frontDoors", "frontDoors"),
		resourceids.UserSpecifiedSegment("frontDoorName", "frontDoorValue"),
		resourceids.StaticSegment("staticFrontendEndpoints", "frontendEndpoints", "frontendEndpoints"),
		resourceids.UserSpecifiedSegment("frontendEndpointName", "frontendEndpointValue"),
	}
}

// String returns a human-readable description of this Frontend Endpoint ID
func (id FrontendEndpointId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Front Door Name: %q", id.FrontDoorName),
		fmt.Sprintf("Frontend Endpoint Name: %q", id.FrontendEndpointName),
	}
	return fmt.Sprintf("Frontend Endpoint (%s)", strings.Join(components, "\n"))
}
