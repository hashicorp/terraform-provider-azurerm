package connectionmonitors

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ConnectionMonitorId{})
}

var _ resourceids.ResourceId = &ConnectionMonitorId{}

// ConnectionMonitorId is a struct representing the Resource ID for a Connection Monitor
type ConnectionMonitorId struct {
	SubscriptionId        string
	ResourceGroupName     string
	NetworkWatcherName    string
	ConnectionMonitorName string
}

// NewConnectionMonitorID returns a new ConnectionMonitorId struct
func NewConnectionMonitorID(subscriptionId string, resourceGroupName string, networkWatcherName string, connectionMonitorName string) ConnectionMonitorId {
	return ConnectionMonitorId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		NetworkWatcherName:    networkWatcherName,
		ConnectionMonitorName: connectionMonitorName,
	}
}

// ParseConnectionMonitorID parses 'input' into a ConnectionMonitorId
func ParseConnectionMonitorID(input string) (*ConnectionMonitorId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConnectionMonitorId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConnectionMonitorId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseConnectionMonitorIDInsensitively parses 'input' case-insensitively into a ConnectionMonitorId
// note: this method should only be used for API response data and not user input
func ParseConnectionMonitorIDInsensitively(input string) (*ConnectionMonitorId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConnectionMonitorId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConnectionMonitorId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ConnectionMonitorId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.NetworkWatcherName, ok = input.Parsed["networkWatcherName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "networkWatcherName", input)
	}

	if id.ConnectionMonitorName, ok = input.Parsed["connectionMonitorName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "connectionMonitorName", input)
	}

	return nil
}

// ValidateConnectionMonitorID checks that 'input' can be parsed as a Connection Monitor ID
func ValidateConnectionMonitorID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseConnectionMonitorID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Connection Monitor ID
func (id ConnectionMonitorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkWatchers/%s/connectionMonitors/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkWatcherName, id.ConnectionMonitorName)
}

// Segments returns a slice of Resource ID Segments which comprise this Connection Monitor ID
func (id ConnectionMonitorId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkWatchers", "networkWatchers", "networkWatchers"),
		resourceids.UserSpecifiedSegment("networkWatcherName", "networkWatcherName"),
		resourceids.StaticSegment("staticConnectionMonitors", "connectionMonitors", "connectionMonitors"),
		resourceids.UserSpecifiedSegment("connectionMonitorName", "connectionMonitorName"),
	}
}

// String returns a human-readable description of this Connection Monitor ID
func (id ConnectionMonitorId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Watcher Name: %q", id.NetworkWatcherName),
		fmt.Sprintf("Connection Monitor Name: %q", id.ConnectionMonitorName),
	}
	return fmt.Sprintf("Connection Monitor (%s)", strings.Join(components, "\n"))
}
