package eventsources

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = &EventSourceId{}

// EventSourceId is a struct representing the Resource ID for a Event Source
type EventSourceId struct {
	SubscriptionId    string
	ResourceGroupName string
	EnvironmentName   string
	EventSourceName   string
}

// NewEventSourceID returns a new EventSourceId struct
func NewEventSourceID(subscriptionId string, resourceGroupName string, environmentName string, eventSourceName string) EventSourceId {
	return EventSourceId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		EnvironmentName:   environmentName,
		EventSourceName:   eventSourceName,
	}
}

// ParseEventSourceID parses 'input' into a EventSourceId
func ParseEventSourceID(input string) (*EventSourceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&EventSourceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := EventSourceId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseEventSourceIDInsensitively parses 'input' case-insensitively into a EventSourceId
// note: this method should only be used for API response data and not user input
func ParseEventSourceIDInsensitively(input string) (*EventSourceId, error) {
	parser := resourceids.NewParserFromResourceIdType(&EventSourceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := EventSourceId{}
	if err := id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *EventSourceId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.EnvironmentName, ok = input.Parsed["environmentName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "environmentName", input)
	}

	if id.EventSourceName, ok = input.Parsed["eventSourceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "eventSourceName", input)
	}

	return nil
}

// ValidateEventSourceID checks that 'input' can be parsed as a Event Source ID
func ValidateEventSourceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseEventSourceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Event Source ID
func (id EventSourceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.TimeSeriesInsights/environments/%s/eventSources/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.EnvironmentName, id.EventSourceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Event Source ID
func (id EventSourceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftTimeSeriesInsights", "Microsoft.TimeSeriesInsights", "Microsoft.TimeSeriesInsights"),
		resourceids.StaticSegment("staticEnvironments", "environments", "environments"),
		resourceids.UserSpecifiedSegment("environmentName", "environmentValue"),
		resourceids.StaticSegment("staticEventSources", "eventSources", "eventSources"),
		resourceids.UserSpecifiedSegment("eventSourceName", "eventSourceValue"),
	}
}

// String returns a human-readable description of this Event Source ID
func (id EventSourceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Environment Name: %q", id.EnvironmentName),
		fmt.Sprintf("Event Source Name: %q", id.EventSourceName),
	}
	return fmt.Sprintf("Event Source (%s)", strings.Join(components, "\n"))
}
