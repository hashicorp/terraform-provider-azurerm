package longtermretentionbackups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = LocationLongTermRetentionServerLongTermRetentionDatabaseId{}

// LocationLongTermRetentionServerLongTermRetentionDatabaseId is a struct representing the Resource ID for a Location Long Term Retention Server Long Term Retention Database
type LocationLongTermRetentionServerLongTermRetentionDatabaseId struct {
	SubscriptionId                string
	ResourceGroupName             string
	LocationName                  string
	LongTermRetentionServerName   string
	LongTermRetentionDatabaseName string
}

// NewLocationLongTermRetentionServerLongTermRetentionDatabaseID returns a new LocationLongTermRetentionServerLongTermRetentionDatabaseId struct
func NewLocationLongTermRetentionServerLongTermRetentionDatabaseID(subscriptionId string, resourceGroupName string, locationName string, longTermRetentionServerName string, longTermRetentionDatabaseName string) LocationLongTermRetentionServerLongTermRetentionDatabaseId {
	return LocationLongTermRetentionServerLongTermRetentionDatabaseId{
		SubscriptionId:                subscriptionId,
		ResourceGroupName:             resourceGroupName,
		LocationName:                  locationName,
		LongTermRetentionServerName:   longTermRetentionServerName,
		LongTermRetentionDatabaseName: longTermRetentionDatabaseName,
	}
}

// ParseLocationLongTermRetentionServerLongTermRetentionDatabaseID parses 'input' into a LocationLongTermRetentionServerLongTermRetentionDatabaseId
func ParseLocationLongTermRetentionServerLongTermRetentionDatabaseID(input string) (*LocationLongTermRetentionServerLongTermRetentionDatabaseId, error) {
	parser := resourceids.NewParserFromResourceIdType(LocationLongTermRetentionServerLongTermRetentionDatabaseId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LocationLongTermRetentionServerLongTermRetentionDatabaseId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
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

// ParseLocationLongTermRetentionServerLongTermRetentionDatabaseIDInsensitively parses 'input' case-insensitively into a LocationLongTermRetentionServerLongTermRetentionDatabaseId
// note: this method should only be used for API response data and not user input
func ParseLocationLongTermRetentionServerLongTermRetentionDatabaseIDInsensitively(input string) (*LocationLongTermRetentionServerLongTermRetentionDatabaseId, error) {
	parser := resourceids.NewParserFromResourceIdType(LocationLongTermRetentionServerLongTermRetentionDatabaseId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LocationLongTermRetentionServerLongTermRetentionDatabaseId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
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

// ValidateLocationLongTermRetentionServerLongTermRetentionDatabaseID checks that 'input' can be parsed as a Location Long Term Retention Server Long Term Retention Database ID
func ValidateLocationLongTermRetentionServerLongTermRetentionDatabaseID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLocationLongTermRetentionServerLongTermRetentionDatabaseID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Location Long Term Retention Server Long Term Retention Database ID
func (id LocationLongTermRetentionServerLongTermRetentionDatabaseId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/locations/%s/longTermRetentionServers/%s/longTermRetentionDatabases/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LocationName, id.LongTermRetentionServerName, id.LongTermRetentionDatabaseName)
}

// Segments returns a slice of Resource ID Segments which comprise this Location Long Term Retention Server Long Term Retention Database ID
func (id LocationLongTermRetentionServerLongTermRetentionDatabaseId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
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

// String returns a human-readable description of this Location Long Term Retention Server Long Term Retention Database ID
func (id LocationLongTermRetentionServerLongTermRetentionDatabaseId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Long Term Retention Server Name: %q", id.LongTermRetentionServerName),
		fmt.Sprintf("Long Term Retention Database Name: %q", id.LongTermRetentionDatabaseName),
	}
	return fmt.Sprintf("Location Long Term Retention Server Long Term Retention Database (%s)", strings.Join(components, "\n"))
}
