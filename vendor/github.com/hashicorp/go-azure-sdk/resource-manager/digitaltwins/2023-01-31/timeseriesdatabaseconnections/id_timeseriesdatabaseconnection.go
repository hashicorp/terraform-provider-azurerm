package timeseriesdatabaseconnections

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&TimeSeriesDatabaseConnectionId{})
}

var _ resourceids.ResourceId = &TimeSeriesDatabaseConnectionId{}

// TimeSeriesDatabaseConnectionId is a struct representing the Resource ID for a Time Series Database Connection
type TimeSeriesDatabaseConnectionId struct {
	SubscriptionId                   string
	ResourceGroupName                string
	DigitalTwinsInstanceName         string
	TimeSeriesDatabaseConnectionName string
}

// NewTimeSeriesDatabaseConnectionID returns a new TimeSeriesDatabaseConnectionId struct
func NewTimeSeriesDatabaseConnectionID(subscriptionId string, resourceGroupName string, digitalTwinsInstanceName string, timeSeriesDatabaseConnectionName string) TimeSeriesDatabaseConnectionId {
	return TimeSeriesDatabaseConnectionId{
		SubscriptionId:                   subscriptionId,
		ResourceGroupName:                resourceGroupName,
		DigitalTwinsInstanceName:         digitalTwinsInstanceName,
		TimeSeriesDatabaseConnectionName: timeSeriesDatabaseConnectionName,
	}
}

// ParseTimeSeriesDatabaseConnectionID parses 'input' into a TimeSeriesDatabaseConnectionId
func ParseTimeSeriesDatabaseConnectionID(input string) (*TimeSeriesDatabaseConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TimeSeriesDatabaseConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TimeSeriesDatabaseConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseTimeSeriesDatabaseConnectionIDInsensitively parses 'input' case-insensitively into a TimeSeriesDatabaseConnectionId
// note: this method should only be used for API response data and not user input
func ParseTimeSeriesDatabaseConnectionIDInsensitively(input string) (*TimeSeriesDatabaseConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&TimeSeriesDatabaseConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := TimeSeriesDatabaseConnectionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *TimeSeriesDatabaseConnectionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.DigitalTwinsInstanceName, ok = input.Parsed["digitalTwinsInstanceName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "digitalTwinsInstanceName", input)
	}

	if id.TimeSeriesDatabaseConnectionName, ok = input.Parsed["timeSeriesDatabaseConnectionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "timeSeriesDatabaseConnectionName", input)
	}

	return nil
}

// ValidateTimeSeriesDatabaseConnectionID checks that 'input' can be parsed as a Time Series Database Connection ID
func ValidateTimeSeriesDatabaseConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseTimeSeriesDatabaseConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Time Series Database Connection ID
func (id TimeSeriesDatabaseConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DigitalTwins/digitalTwinsInstances/%s/timeSeriesDatabaseConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DigitalTwinsInstanceName, id.TimeSeriesDatabaseConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Time Series Database Connection ID
func (id TimeSeriesDatabaseConnectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDigitalTwins", "Microsoft.DigitalTwins", "Microsoft.DigitalTwins"),
		resourceids.StaticSegment("staticDigitalTwinsInstances", "digitalTwinsInstances", "digitalTwinsInstances"),
		resourceids.UserSpecifiedSegment("digitalTwinsInstanceName", "digitalTwinsInstanceName"),
		resourceids.StaticSegment("staticTimeSeriesDatabaseConnections", "timeSeriesDatabaseConnections", "timeSeriesDatabaseConnections"),
		resourceids.UserSpecifiedSegment("timeSeriesDatabaseConnectionName", "timeSeriesDatabaseConnectionName"),
	}
}

// String returns a human-readable description of this Time Series Database Connection ID
func (id TimeSeriesDatabaseConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Digital Twins Instance Name: %q", id.DigitalTwinsInstanceName),
		fmt.Sprintf("Time Series Database Connection Name: %q", id.TimeSeriesDatabaseConnectionName),
	}
	return fmt.Sprintf("Time Series Database Connection (%s)", strings.Join(components, "\n"))
}
