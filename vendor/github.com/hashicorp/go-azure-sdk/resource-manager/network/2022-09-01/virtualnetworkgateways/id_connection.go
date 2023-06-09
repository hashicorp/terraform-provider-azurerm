package virtualnetworkgateways

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = ConnectionId{}

// ConnectionId is a struct representing the Resource ID for a Connection
type ConnectionId struct {
	SubscriptionId    string
	ResourceGroupName string
	ConnectionName    string
}

// NewConnectionID returns a new ConnectionId struct
func NewConnectionID(subscriptionId string, resourceGroupName string, connectionName string) ConnectionId {
	return ConnectionId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ConnectionName:    connectionName,
	}
}

// ParseConnectionID parses 'input' into a ConnectionId
func ParseConnectionID(input string) (*ConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ConnectionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ConnectionName, ok = parsed.Parsed["connectionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "connectionName", *parsed)
	}

	return &id, nil
}

// ParseConnectionIDInsensitively parses 'input' case-insensitively into a ConnectionId
// note: this method should only be used for API response data and not user input
func ParseConnectionIDInsensitively(input string) (*ConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ConnectionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ConnectionName, ok = parsed.Parsed["connectionName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "connectionName", *parsed)
	}

	return &id, nil
}

// ValidateConnectionID checks that 'input' can be parsed as a Connection ID
func ValidateConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Connection ID
func (id ConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/connections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Connection ID
func (id ConnectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticConnections", "connections", "connections"),
		resourceids.UserSpecifiedSegment("connectionName", "connectionValue"),
	}
}

// String returns a human-readable description of this Connection ID
func (id ConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Connection Name: %q", id.ConnectionName),
	}
	return fmt.Sprintf("Connection (%s)", strings.Join(components, "\n"))
}
