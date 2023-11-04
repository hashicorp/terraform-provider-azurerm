package longtermretentionbackups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = LocationLongTermRetentionServerId{}

// LocationLongTermRetentionServerId is a struct representing the Resource ID for a Location Long Term Retention Server
type LocationLongTermRetentionServerId struct {
	SubscriptionId              string
	ResourceGroupName           string
	LocationName                string
	LongTermRetentionServerName string
}

// NewLocationLongTermRetentionServerID returns a new LocationLongTermRetentionServerId struct
func NewLocationLongTermRetentionServerID(subscriptionId string, resourceGroupName string, locationName string, longTermRetentionServerName string) LocationLongTermRetentionServerId {
	return LocationLongTermRetentionServerId{
		SubscriptionId:              subscriptionId,
		ResourceGroupName:           resourceGroupName,
		LocationName:                locationName,
		LongTermRetentionServerName: longTermRetentionServerName,
	}
}

// ParseLocationLongTermRetentionServerID parses 'input' into a LocationLongTermRetentionServerId
func ParseLocationLongTermRetentionServerID(input string) (*LocationLongTermRetentionServerId, error) {
	parser := resourceids.NewParserFromResourceIdType(LocationLongTermRetentionServerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LocationLongTermRetentionServerId{}

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

	return &id, nil
}

// ParseLocationLongTermRetentionServerIDInsensitively parses 'input' case-insensitively into a LocationLongTermRetentionServerId
// note: this method should only be used for API response data and not user input
func ParseLocationLongTermRetentionServerIDInsensitively(input string) (*LocationLongTermRetentionServerId, error) {
	parser := resourceids.NewParserFromResourceIdType(LocationLongTermRetentionServerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LocationLongTermRetentionServerId{}

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

	return &id, nil
}

// ValidateLocationLongTermRetentionServerID checks that 'input' can be parsed as a Location Long Term Retention Server ID
func ValidateLocationLongTermRetentionServerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLocationLongTermRetentionServerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Location Long Term Retention Server ID
func (id LocationLongTermRetentionServerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/locations/%s/longTermRetentionServers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LocationName, id.LongTermRetentionServerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Location Long Term Retention Server ID
func (id LocationLongTermRetentionServerId) Segments() []resourceids.Segment {
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
	}
}

// String returns a human-readable description of this Location Long Term Retention Server ID
func (id LocationLongTermRetentionServerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Long Term Retention Server Name: %q", id.LongTermRetentionServerName),
	}
	return fmt.Sprintf("Location Long Term Retention Server (%s)", strings.Join(components, "\n"))
}
