package longtermretentionbackups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = LongTermRetentionBackupId{}

// LongTermRetentionBackupId is a struct representing the Resource ID for a Long Term Retention Backup
type LongTermRetentionBackupId struct {
	SubscriptionId                string
	LocationName                  string
	LongTermRetentionServerName   string
	LongTermRetentionDatabaseName string
	LongTermRetentionBackupName   string
}

// NewLongTermRetentionBackupID returns a new LongTermRetentionBackupId struct
func NewLongTermRetentionBackupID(subscriptionId string, locationName string, longTermRetentionServerName string, longTermRetentionDatabaseName string, longTermRetentionBackupName string) LongTermRetentionBackupId {
	return LongTermRetentionBackupId{
		SubscriptionId:                subscriptionId,
		LocationName:                  locationName,
		LongTermRetentionServerName:   longTermRetentionServerName,
		LongTermRetentionDatabaseName: longTermRetentionDatabaseName,
		LongTermRetentionBackupName:   longTermRetentionBackupName,
	}
}

// ParseLongTermRetentionBackupID parses 'input' into a LongTermRetentionBackupId
func ParseLongTermRetentionBackupID(input string) (*LongTermRetentionBackupId, error) {
	parser := resourceids.NewParserFromResourceIdType(LongTermRetentionBackupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LongTermRetentionBackupId{}

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

	if id.LongTermRetentionBackupName, ok = parsed.Parsed["longTermRetentionBackupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "longTermRetentionBackupName", *parsed)
	}

	return &id, nil
}

// ParseLongTermRetentionBackupIDInsensitively parses 'input' case-insensitively into a LongTermRetentionBackupId
// note: this method should only be used for API response data and not user input
func ParseLongTermRetentionBackupIDInsensitively(input string) (*LongTermRetentionBackupId, error) {
	parser := resourceids.NewParserFromResourceIdType(LongTermRetentionBackupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LongTermRetentionBackupId{}

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

	if id.LongTermRetentionBackupName, ok = parsed.Parsed["longTermRetentionBackupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "longTermRetentionBackupName", *parsed)
	}

	return &id, nil
}

// ValidateLongTermRetentionBackupID checks that 'input' can be parsed as a Long Term Retention Backup ID
func ValidateLongTermRetentionBackupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLongTermRetentionBackupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Long Term Retention Backup ID
func (id LongTermRetentionBackupId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Sql/locations/%s/longTermRetentionServers/%s/longTermRetentionDatabases/%s/longTermRetentionBackups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.LongTermRetentionServerName, id.LongTermRetentionDatabaseName, id.LongTermRetentionBackupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Long Term Retention Backup ID
func (id LongTermRetentionBackupId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticLongTermRetentionBackups", "longTermRetentionBackups", "longTermRetentionBackups"),
		resourceids.UserSpecifiedSegment("longTermRetentionBackupName", "longTermRetentionBackupValue"),
	}
}

// String returns a human-readable description of this Long Term Retention Backup ID
func (id LongTermRetentionBackupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Long Term Retention Server Name: %q", id.LongTermRetentionServerName),
		fmt.Sprintf("Long Term Retention Database Name: %q", id.LongTermRetentionDatabaseName),
		fmt.Sprintf("Long Term Retention Backup Name: %q", id.LongTermRetentionBackupName),
	}
	return fmt.Sprintf("Long Term Retention Backup (%s)", strings.Join(components, "\n"))
}
