package flowlogs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = FlowLogId{}

// FlowLogId is a struct representing the Resource ID for a Flow Log
type FlowLogId struct {
	SubscriptionId     string
	ResourceGroupName  string
	NetworkWatcherName string
	FlowLogName        string
}

// NewFlowLogID returns a new FlowLogId struct
func NewFlowLogID(subscriptionId string, resourceGroupName string, networkWatcherName string, flowLogName string) FlowLogId {
	return FlowLogId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		NetworkWatcherName: networkWatcherName,
		FlowLogName:        flowLogName,
	}
}

// ParseFlowLogID parses 'input' into a FlowLogId
func ParseFlowLogID(input string) (*FlowLogId, error) {
	parser := resourceids.NewParserFromResourceIdType(FlowLogId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FlowLogId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetworkWatcherName, ok = parsed.Parsed["networkWatcherName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkWatcherName", *parsed)
	}

	if id.FlowLogName, ok = parsed.Parsed["flowLogName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "flowLogName", *parsed)
	}

	return &id, nil
}

// ParseFlowLogIDInsensitively parses 'input' case-insensitively into a FlowLogId
// note: this method should only be used for API response data and not user input
func ParseFlowLogIDInsensitively(input string) (*FlowLogId, error) {
	parser := resourceids.NewParserFromResourceIdType(FlowLogId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FlowLogId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.NetworkWatcherName, ok = parsed.Parsed["networkWatcherName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "networkWatcherName", *parsed)
	}

	if id.FlowLogName, ok = parsed.Parsed["flowLogName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "flowLogName", *parsed)
	}

	return &id, nil
}

// ValidateFlowLogID checks that 'input' can be parsed as a Flow Log ID
func ValidateFlowLogID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseFlowLogID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Flow Log ID
func (id FlowLogId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkWatchers/%s/flowLogs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NetworkWatcherName, id.FlowLogName)
}

// Segments returns a slice of Resource ID Segments which comprise this Flow Log ID
func (id FlowLogId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticNetworkWatchers", "networkWatchers", "networkWatchers"),
		resourceids.UserSpecifiedSegment("networkWatcherName", "networkWatcherValue"),
		resourceids.StaticSegment("staticFlowLogs", "flowLogs", "flowLogs"),
		resourceids.UserSpecifiedSegment("flowLogName", "flowLogValue"),
	}
}

// String returns a human-readable description of this Flow Log ID
func (id FlowLogId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Network Watcher Name: %q", id.NetworkWatcherName),
		fmt.Sprintf("Flow Log Name: %q", id.FlowLogName),
	}
	return fmt.Sprintf("Flow Log (%s)", strings.Join(components, "\n"))
}
