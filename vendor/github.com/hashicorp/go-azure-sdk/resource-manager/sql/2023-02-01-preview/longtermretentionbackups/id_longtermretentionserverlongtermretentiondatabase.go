package longtermretentionbackups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = LongTermRetentionServerLongTermRetentionDatabaseId{}

// LongTermRetentionServerLongTermRetentionDatabaseId is a struct representing the Resource ID for a Long Term Retention Server Long Term Retention Database
type LongTermRetentionServerLongTermRetentionDatabaseId struct {
	SubscriptionId                string
	LocationName                  string
	LongTermRetentionServerName   string
	LongTermRetentionDatabaseName string
}

// NewLongTermRetentionServerLongTermRetentionDatabaseID returns a new LongTermRetentionServerLongTermRetentionDatabaseId struct
func NewLongTermRetentionServerLongTermRetentionDatabaseID(subscriptionId string, locationName string, longTermRetentionServerName string, longTermRetentionDatabaseName string) LongTermRetentionServerLongTermRetentionDatabaseId {
	return LongTermRetentionServerLongTermRetentionDatabaseId{
		SubscriptionId:                subscriptionId,
		LocationName:                  locationName,
		LongTermRetentionServerName:   longTermRetentionServerName,
		LongTermRetentionDatabaseName: longTermRetentionDatabaseName,
	}
}

// ParseLongTermRetentionServerLongTermRetentionDatabaseID parses 'input' into a LongTermRetentionServerLongTermRetentionDatabaseId
func ParseLongTermRetentionServerLongTermRetentionDatabaseID(input string) (*LongTermRetentionServerLongTermRetentionDatabaseId, error) {
	parser := resourceids.NewParserFromResourceIdType(LongTermRetentionServerLongTermRetentionDatabaseId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LongTermRetentionServerLongTermRetentionDatabaseId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.LocationName, ok = parsed.Parsed["locationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "locationName", *parsed)
	}

	if id.LongTermRetentionServerName, ok = parsed.Parsed["longTermRetentionServerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "longTermRetentionServerName", *parsed)
	}

	if id.LongTermRetentionDatabaseName, ok = parsed.Parsed["longTermRetentionDatabaseName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "longTermRetentionDatabaseName", *parsed)
	}

	return &id, nil
}

// ParseLongTermRetentionServerLongTermRetentionDatabaseIDInsensitively parses 'input' case-insensitively into a LongTermRetentionServerLongTermRetentionDatabaseId
// note: this method should only be used for API response data and not user input
func ParseLongTermRetentionServerLongTermRetentionDatabaseIDInsensitively(input string) (*LongTermRetentionServerLongTermRetentionDatabaseId, error) {
	parser := resourceids.NewParserFromResourceIdType(LongTermRetentionServerLongTermRetentionDatabaseId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LongTermRetentionServerLongTermRetentionDatabaseId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.LocationName, ok = parsed.Parsed["locationName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "locationName", *parsed)
	}

	if id.LongTermRetentionServerName, ok = parsed.Parsed["longTermRetentionServerName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "longTermRetentionServerName", *parsed)
	}

	if id.LongTermRetentionDatabaseName, ok = parsed.Parsed["longTermRetentionDatabaseName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "longTermRetentionDatabaseName", *parsed)
	}

	return &id, nil
}

// ValidateLongTermRetentionServerLongTermRetentionDatabaseID checks that 'input' can be parsed as a Long Term Retention Server Long Term Retention Database ID
func ValidateLongTermRetentionServerLongTermRetentionDatabaseID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLongTermRetentionServerLongTermRetentionDatabaseID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Long Term Retention Server Long Term Retention Database ID
func (id LongTermRetentionServerLongTermRetentionDatabaseId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Sql/locations/%s/longTermRetentionServers/%s/longTermRetentionDatabases/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.LongTermRetentionServerName, id.LongTermRetentionDatabaseName)
}

// Segments returns a slice of Resource ID Segments which comprise this Long Term Retention Server Long Term Retention Database ID
func (id LongTermRetentionServerLongTermRetentionDatabaseId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSql", "Microsoft.Sql", "Microsoft.Sql"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationValue"),
		resourceids.StaticSegment("staticLongTermRetentionServers", "longTermRetentionServers", "longTermRetentionServers"),
		resourceids.UserSpecifiedSegment("longTermRetentionServerName", "longTermRetentionServerValue"),
		resourceids.StaticSegment("staticLongTermRetentionDatabases", "longTermRetentionDatabases", "longTermRetentionDatabases"),
		resourceids.UserSpecifiedSegment("longTermRetentionDatabaseName", "longTermRetentionDatabaseValue"),
	}
}

// String returns a human-readable description of this Long Term Retention Server Long Term Retention Database ID
func (id LongTermRetentionServerLongTermRetentionDatabaseId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Long Term Retention Server Name: %q", id.LongTermRetentionServerName),
		fmt.Sprintf("Long Term Retention Database Name: %q", id.LongTermRetentionDatabaseName),
	}
	return fmt.Sprintf("Long Term Retention Server Long Term Retention Database (%s)", strings.Join(components, "\n"))
}
