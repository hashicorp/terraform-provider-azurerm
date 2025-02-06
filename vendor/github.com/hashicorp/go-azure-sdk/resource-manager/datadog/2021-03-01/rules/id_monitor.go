package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&MonitorId{})
}

var _ resourceids.ResourceId = &MonitorId{}

// MonitorId is a struct representing the Resource ID for a Monitor
type MonitorId struct {
	SubscriptionId    string
	ResourceGroupName string
	MonitorName       string
}

// NewMonitorID returns a new MonitorId struct
func NewMonitorID(subscriptionId string, resourceGroupName string, monitorName string) MonitorId {
	return MonitorId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		MonitorName:       monitorName,
	}
}

// ParseMonitorID parses 'input' into a MonitorId
func ParseMonitorID(input string) (*MonitorId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MonitorId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MonitorId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseMonitorIDInsensitively parses 'input' case-insensitively into a MonitorId
// note: this method should only be used for API response data and not user input
func ParseMonitorIDInsensitively(input string) (*MonitorId, error) {
	parser := resourceids.NewParserFromResourceIdType(&MonitorId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := MonitorId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *MonitorId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.MonitorName, ok = input.Parsed["monitorName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "monitorName", input)
	}

	return nil
}

// ValidateMonitorID checks that 'input' can be parsed as a Monitor ID
func ValidateMonitorID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseMonitorID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Monitor ID
func (id MonitorId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Datadog/monitors/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.MonitorName)
}

// Segments returns a slice of Resource ID Segments which comprise this Monitor ID
func (id MonitorId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDatadog", "Microsoft.Datadog", "Microsoft.Datadog"),
		resourceids.StaticSegment("staticMonitors", "monitors", "monitors"),
		resourceids.UserSpecifiedSegment("monitorName", "monitorName"),
	}
}

// String returns a human-readable description of this Monitor ID
func (id MonitorId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Monitor Name: %q", id.MonitorName),
	}
	return fmt.Sprintf("Monitor (%s)", strings.Join(components, "\n"))
}
