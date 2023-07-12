package privateendpoints

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = PrivateEndpointId{}

// PrivateEndpointId is a struct representing the Resource ID for a Private Endpoint
type PrivateEndpointId struct {
	SubscriptionId      string
	ResourceGroupName   string
	PrivateEndpointName string
}

// NewPrivateEndpointID returns a new PrivateEndpointId struct
func NewPrivateEndpointID(subscriptionId string, resourceGroupName string, privateEndpointName string) PrivateEndpointId {
	return PrivateEndpointId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		PrivateEndpointName: privateEndpointName,
	}
}

// ParsePrivateEndpointID parses 'input' into a PrivateEndpointId
func ParsePrivateEndpointID(input string) (*PrivateEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(PrivateEndpointId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PrivateEndpointId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.PrivateEndpointName, ok = parsed.Parsed["privateEndpointName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "privateEndpointName", *parsed)
	}

	return &id, nil
}

// ParsePrivateEndpointIDInsensitively parses 'input' case-insensitively into a PrivateEndpointId
// note: this method should only be used for API response data and not user input
func ParsePrivateEndpointIDInsensitively(input string) (*PrivateEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(PrivateEndpointId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PrivateEndpointId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.PrivateEndpointName, ok = parsed.Parsed["privateEndpointName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "privateEndpointName", *parsed)
	}

	return &id, nil
}

// ValidatePrivateEndpointID checks that 'input' can be parsed as a Private Endpoint ID
func ValidatePrivateEndpointID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePrivateEndpointID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Private Endpoint ID
func (id PrivateEndpointId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/privateEndpoints/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.PrivateEndpointName)
}

// Segments returns a slice of Resource ID Segments which comprise this Private Endpoint ID
func (id PrivateEndpointId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticPrivateEndpoints", "privateEndpoints", "privateEndpoints"),
		resourceids.UserSpecifiedSegment("privateEndpointName", "privateEndpointValue"),
	}
}

// String returns a human-readable description of this Private Endpoint ID
func (id PrivateEndpointId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Private Endpoint Name: %q", id.PrivateEndpointName),
	}
	return fmt.Sprintf("Private Endpoint (%s)", strings.Join(components, "\n"))
}
