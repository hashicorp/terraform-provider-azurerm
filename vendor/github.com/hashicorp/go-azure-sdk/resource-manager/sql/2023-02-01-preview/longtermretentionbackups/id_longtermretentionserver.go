package longtermretentionbackups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = LongTermRetentionServerId{}

// LongTermRetentionServerId is a struct representing the Resource ID for a Long Term Retention Server
type LongTermRetentionServerId struct {
	SubscriptionId              string
	LocationName                string
	LongTermRetentionServerName string
}

// NewLongTermRetentionServerID returns a new LongTermRetentionServerId struct
func NewLongTermRetentionServerID(subscriptionId string, locationName string, longTermRetentionServerName string) LongTermRetentionServerId {
	return LongTermRetentionServerId{
		SubscriptionId:              subscriptionId,
		LocationName:                locationName,
		LongTermRetentionServerName: longTermRetentionServerName,
	}
}

// ParseLongTermRetentionServerID parses 'input' into a LongTermRetentionServerId
func ParseLongTermRetentionServerID(input string) (*LongTermRetentionServerId, error) {
	parser := resourceids.NewParserFromResourceIdType(LongTermRetentionServerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LongTermRetentionServerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.LocationName, ok = parsed.Parsed["locationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "locationName", *parsed)
	}

	if id.LongTermRetentionServerName, ok = parsed.Parsed["longTermRetentionServerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "longTermRetentionServerName", *parsed)
	}

	return &id, nil
}

// ParseLongTermRetentionServerIDInsensitively parses 'input' case-insensitively into a LongTermRetentionServerId
// note: this method should only be used for API response data and not user input
func ParseLongTermRetentionServerIDInsensitively(input string) (*LongTermRetentionServerId, error) {
	parser := resourceids.NewParserFromResourceIdType(LongTermRetentionServerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LongTermRetentionServerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.LocationName, ok = parsed.Parsed["locationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "locationName", *parsed)
	}

	if id.LongTermRetentionServerName, ok = parsed.Parsed["longTermRetentionServerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "longTermRetentionServerName", *parsed)
	}

	return &id, nil
}

// ValidateLongTermRetentionServerID checks that 'input' can be parsed as a Long Term Retention Server ID
func ValidateLongTermRetentionServerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLongTermRetentionServerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Long Term Retention Server ID
func (id LongTermRetentionServerId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Sql/locations/%s/longTermRetentionServers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.LongTermRetentionServerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Long Term Retention Server ID
func (id LongTermRetentionServerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSql", "Microsoft.Sql", "Microsoft.Sql"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationValue"),
		resourceids.StaticSegment("staticLongTermRetentionServers", "longTermRetentionServers", "longTermRetentionServers"),
		resourceids.UserSpecifiedSegment("longTermRetentionServerName", "longTermRetentionServerValue"),
	}
}

// String returns a human-readable description of this Long Term Retention Server ID
func (id LongTermRetentionServerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Long Term Retention Server Name: %q", id.LongTermRetentionServerName),
	}
	return fmt.Sprintf("Long Term Retention Server (%s)", strings.Join(components, "\n"))
}
