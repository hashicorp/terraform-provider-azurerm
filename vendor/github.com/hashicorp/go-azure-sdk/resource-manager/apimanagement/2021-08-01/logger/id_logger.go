package logger

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = LoggerId{}

// LoggerId is a struct representing the Resource ID for a Logger
type LoggerId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	LoggerId          string
}

// NewLoggerID returns a new LoggerId struct
func NewLoggerID(subscriptionId string, resourceGroupName string, serviceName string, loggerId string) LoggerId {
	return LoggerId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		LoggerId:          loggerId,
	}
}

// ParseLoggerID parses 'input' into a LoggerId
func ParseLoggerID(input string) (*LoggerId, error) {
	parser := resourceids.NewParserFromResourceIdType(LoggerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LoggerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServiceName, ok = parsed.Parsed["serviceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serviceName", *parsed)
	}

	if id.LoggerId, ok = parsed.Parsed["loggerId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "loggerId", *parsed)
	}

	return &id, nil
}

// ParseLoggerIDInsensitively parses 'input' case-insensitively into a LoggerId
// note: this method should only be used for API response data and not user input
func ParseLoggerIDInsensitively(input string) (*LoggerId, error) {
	parser := resourceids.NewParserFromResourceIdType(LoggerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LoggerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.ServiceName, ok = parsed.Parsed["serviceName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "serviceName", *parsed)
	}

	if id.LoggerId, ok = parsed.Parsed["loggerId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "loggerId", *parsed)
	}

	return &id, nil
}

// ValidateLoggerID checks that 'input' can be parsed as a Logger ID
func ValidateLoggerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLoggerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Logger ID
func (id LoggerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ApiManagement/service/%s/loggers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.LoggerId)
}

// Segments returns a slice of Resource ID Segments which comprise this Logger ID
func (id LoggerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApiManagement", "Microsoft.ApiManagement", "Microsoft.ApiManagement"),
		resourceids.StaticSegment("staticService", "service", "service"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceValue"),
		resourceids.StaticSegment("staticLoggers", "loggers", "loggers"),
		resourceids.UserSpecifiedSegment("loggerId", "loggerIdValue"),
	}
}

// String returns a human-readable description of this Logger ID
func (id LoggerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Logger: %q", id.LoggerId),
	}
	return fmt.Sprintf("Logger (%s)", strings.Join(components, "\n"))
}
