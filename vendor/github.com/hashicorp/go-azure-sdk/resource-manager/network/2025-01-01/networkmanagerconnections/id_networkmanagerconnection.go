package networkmanagerconnections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&NetworkManagerConnectionId{})
}

var _ resourceids.ResourceId = &NetworkManagerConnectionId{}

// NetworkManagerConnectionId is a struct representing the Resource ID for a Network Manager Connection
type NetworkManagerConnectionId struct {
	SubscriptionId               string
	NetworkManagerConnectionName string
}

// NewNetworkManagerConnectionID returns a new NetworkManagerConnectionId struct
func NewNetworkManagerConnectionID(subscriptionId string, networkManagerConnectionName string) NetworkManagerConnectionId {
	return NetworkManagerConnectionId{
		SubscriptionId:               subscriptionId,
		NetworkManagerConnectionName: networkManagerConnectionName,
	}
}

// ParseNetworkManagerConnectionID parses 'input' into a NetworkManagerConnectionId
func ParseNetworkManagerConnectionID(input string) (*NetworkManagerConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NetworkManagerConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NetworkManagerConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseNetworkManagerConnectionIDInsensitively parses 'input' case-insensitively into a NetworkManagerConnectionId
// note: this method should only be used for API response data and not user input
func ParseNetworkManagerConnectionIDInsensitively(input string) (*NetworkManagerConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&NetworkManagerConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := NetworkManagerConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *NetworkManagerConnectionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.NetworkManagerConnectionName, ok = input.Parsed["networkManagerConnectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "networkManagerConnectionName", input)
	}

	return nil
}

// ValidateNetworkManagerConnectionID checks that 'input' can be parsed as a Network Manager Connection ID
func ValidateNetworkManagerConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseNetworkManagerConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Network Manager Connection ID
func (id NetworkManagerConnectionId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Network/networkManagerConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.NetworkManagerConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Network Manager Connection ID
func (id NetworkManagerConnectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkManagerConnections", "networkManagerConnections", "networkManagerConnections"),
		resourceids.UserSpecifiedSegment("networkManagerConnectionName", "networkManagerConnectionName"),
	}
}

// String returns a human-readable description of this Network Manager Connection ID
func (id NetworkManagerConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Network Manager Connection Name: %q", id.NetworkManagerConnectionName),
	}
	return fmt.Sprintf("Network Manager Connection (%s)", strings.Join(components, "\n"))
}
