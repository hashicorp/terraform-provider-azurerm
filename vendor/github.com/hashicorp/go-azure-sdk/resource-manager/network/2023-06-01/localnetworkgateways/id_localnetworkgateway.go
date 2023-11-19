package localnetworkgateways

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = LocalNetworkGatewayId{}

// LocalNetworkGatewayId is a struct representing the Resource ID for a Local Network Gateway
type LocalNetworkGatewayId struct {
	SubscriptionId          string
	ResourceGroupName       string
	LocalNetworkGatewayName string
}

// NewLocalNetworkGatewayID returns a new LocalNetworkGatewayId struct
func NewLocalNetworkGatewayID(subscriptionId string, resourceGroupName string, localNetworkGatewayName string) LocalNetworkGatewayId {
	return LocalNetworkGatewayId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		LocalNetworkGatewayName: localNetworkGatewayName,
	}
}

// ParseLocalNetworkGatewayID parses 'input' into a LocalNetworkGatewayId
func ParseLocalNetworkGatewayID(input string) (*LocalNetworkGatewayId, error) {
	parser := resourceids.NewParserFromResourceIdType(LocalNetworkGatewayId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LocalNetworkGatewayId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LocalNetworkGatewayName, ok = parsed.Parsed["localNetworkGatewayName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "localNetworkGatewayName", *parsed)
	}

	return &id, nil
}

// ParseLocalNetworkGatewayIDInsensitively parses 'input' case-insensitively into a LocalNetworkGatewayId
// note: this method should only be used for API response data and not user input
func ParseLocalNetworkGatewayIDInsensitively(input string) (*LocalNetworkGatewayId, error) {
	parser := resourceids.NewParserFromResourceIdType(LocalNetworkGatewayId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LocalNetworkGatewayId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.LocalNetworkGatewayName, ok = parsed.Parsed["localNetworkGatewayName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "localNetworkGatewayName", *parsed)
	}

	return &id, nil
}

// ValidateLocalNetworkGatewayID checks that 'input' can be parsed as a Local Network Gateway ID
func ValidateLocalNetworkGatewayID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLocalNetworkGatewayID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Local Network Gateway ID
func (id LocalNetworkGatewayId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/localNetworkGateways/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LocalNetworkGatewayName)
}

// Segments returns a slice of Resource ID Segments which comprise this Local Network Gateway ID
func (id LocalNetworkGatewayId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticLocalNetworkGateways", "localNetworkGateways", "localNetworkGateways"),
		resourceids.UserSpecifiedSegment("localNetworkGatewayName", "localNetworkGatewayValue"),
	}
}

// String returns a human-readable description of this Local Network Gateway ID
func (id LocalNetworkGatewayId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Local Network Gateway Name: %q", id.LocalNetworkGatewayName),
	}
	return fmt.Sprintf("Local Network Gateway (%s)", strings.Join(components, "\n"))
}
