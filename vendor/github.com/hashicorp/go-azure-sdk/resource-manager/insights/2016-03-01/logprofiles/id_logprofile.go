package logprofiles

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = LogProfileId{}

// LogProfileId is a struct representing the Resource ID for a Log Profile
type LogProfileId struct {
	SubscriptionId string
	LogProfileName string
}

// NewLogProfileID returns a new LogProfileId struct
func NewLogProfileID(subscriptionId string, logProfileName string) LogProfileId {
	return LogProfileId{
		SubscriptionId: subscriptionId,
		LogProfileName: logProfileName,
	}
}

// ParseLogProfileID parses 'input' into a LogProfileId
func ParseLogProfileID(input string) (*LogProfileId, error) {
	parser := resourceids.NewParserFromResourceIdType(LogProfileId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LogProfileId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.LogProfileName, ok = parsed.Parsed["logProfileName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "logProfileName", *parsed)
	}

	return &id, nil
}

// ParseLogProfileIDInsensitively parses 'input' case-insensitively into a LogProfileId
// note: this method should only be used for API response data and not user input
func ParseLogProfileIDInsensitively(input string) (*LogProfileId, error) {
	parser := resourceids.NewParserFromResourceIdType(LogProfileId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LogProfileId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.LogProfileName, ok = parsed.Parsed["logProfileName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "logProfileName", *parsed)
	}

	return &id, nil
}

// ValidateLogProfileID checks that 'input' can be parsed as a Log Profile ID
func ValidateLogProfileID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLogProfileID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Log Profile ID
func (id LogProfileId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Insights/logProfiles/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LogProfileName)
}

// Segments returns a slice of Resource ID Segments which comprise this Log Profile ID
func (id LogProfileId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftInsights", "Microsoft.Insights", "Microsoft.Insights"),
		resourceids.StaticSegment("staticLogProfiles", "logProfiles", "logProfiles"),
		resourceids.UserSpecifiedSegment("logProfileName", "logProfileValue"),
	}
}

// String returns a human-readable description of this Log Profile ID
func (id LogProfileId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Log Profile Name: %q", id.LogProfileName),
	}
	return fmt.Sprintf("Log Profile (%s)", strings.Join(components, "\n"))
}
