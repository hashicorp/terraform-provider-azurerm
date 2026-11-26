package connectionanalyzers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&ConnectionAnalyzerId{})
}

var _ resourceids.ResourceId = &ConnectionAnalyzerId{}

// ConnectionAnalyzerId is a struct representing the Resource ID for a Connection Analyzer
type ConnectionAnalyzerId struct {
	SubscriptionId         string
	ResourceGroupName      string
	NetworkWatcherName     string
	ConnectionAnalyzerName string
}

// NewConnectionAnalyzerID returns a new ConnectionAnalyzerId struct
func NewConnectionAnalyzerID(subscriptionId string, resourceGroupName string, networkWatcherName string, connectionAnalyzerName string) ConnectionAnalyzerId {
	return ConnectionAnalyzerId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		NetworkWatcherName:     networkWatcherName,
		ConnectionAnalyzerName: connectionAnalyzerName,
	}
}

// ParseConnectionAnalyzerID parses 'input' into a ConnectionAnalyzerId
func ParseConnectionAnalyzerID(input string) (*ConnectionAnalyzerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConnectionAnalyzerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConnectionAnalyzerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseConnectionAnalyzerIDInsensitively parses 'input' case-insensitively into a ConnectionAnalyzerId
// note: this method should only be used for API response data and not user input
func ParseConnectionAnalyzerIDInsensitively(input string) (*ConnectionAnalyzerId, error) {
	parser := resourceids.NewParserFromResourceIdType(&ConnectionAnalyzerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := ConnectionAnalyzerId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *ConnectionAnalyzerId) FromParseResult(input resourceids.ParseResult) error {
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

	if id.ConnectionAnalyzerName, ok = input.Parsed["connectionAnalyzerName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "connectionAnalyzerName", input)
	}

	return nil
}

// ValidateConnectionAnalyzerID checks that 'input' can be parsed as a Connection Analyzer ID
func ValidateConnectionAnalyzerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseConnectionAnalyzerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Connection Analyzer ID
func (id ConnectionAnalyzerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkWatchers/%s/connectionAnalyzers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkWatcherName, id.ConnectionAnalyzerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Connection Analyzer ID
func (id ConnectionAnalyzerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkWatchers", "networkWatchers", "networkWatchers"),
		resourceids.UserSpecifiedSegment("networkWatcherName", "networkWatcherName"),
		resourceids.StaticSegment("staticConnectionAnalyzers", "connectionAnalyzers", "connectionAnalyzers"),
		resourceids.UserSpecifiedSegment("connectionAnalyzerName", "connectionAnalyzerName"),
	}
}

// String returns a human-readable description of this Connection Analyzer ID
func (id ConnectionAnalyzerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Watcher Name: %q", id.NetworkWatcherName),
		fmt.Sprintf("Connection Analyzer Name: %q", id.ConnectionAnalyzerName),
	}
	return fmt.Sprintf("Connection Analyzer (%s)", strings.Join(components, "\n"))
}
