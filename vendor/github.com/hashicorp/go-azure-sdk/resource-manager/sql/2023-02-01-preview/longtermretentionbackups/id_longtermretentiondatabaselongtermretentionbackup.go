package longtermretentionbackups

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = LongTermRetentionDatabaseLongTermRetentionBackupId{}

// LongTermRetentionDatabaseLongTermRetentionBackupId is a struct representing the Resource ID for a Long Term Retention Database Long Term Retention Backup
type LongTermRetentionDatabaseLongTermRetentionBackupId struct {
	SubscriptionId                string
	ResourceGroupName             string
	LocationName                  string
	LongTermRetentionServerName   string
	LongTermRetentionDatabaseName string
	LongTermRetentionBackupName   string
}

// NewLongTermRetentionDatabaseLongTermRetentionBackupID returns a new LongTermRetentionDatabaseLongTermRetentionBackupId struct
func NewLongTermRetentionDatabaseLongTermRetentionBackupID(subscriptionId string, resourceGroupName string, locationName string, longTermRetentionServerName string, longTermRetentionDatabaseName string, longTermRetentionBackupName string) LongTermRetentionDatabaseLongTermRetentionBackupId {
	return LongTermRetentionDatabaseLongTermRetentionBackupId{
		SubscriptionId:                subscriptionId,
		ResourceGroupName:             resourceGroupName,
		LocationName:                  locationName,
		LongTermRetentionServerName:   longTermRetentionServerName,
		LongTermRetentionDatabaseName: longTermRetentionDatabaseName,
		LongTermRetentionBackupName:   longTermRetentionBackupName,
	}
}

// ParseLongTermRetentionDatabaseLongTermRetentionBackupID parses 'input' into a LongTermRetentionDatabaseLongTermRetentionBackupId
func ParseLongTermRetentionDatabaseLongTermRetentionBackupID(input string) (*LongTermRetentionDatabaseLongTermRetentionBackupId, error) {
	parser := resourceids.NewParserFromResourceIdType(LongTermRetentionDatabaseLongTermRetentionBackupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LongTermRetentionDatabaseLongTermRetentionBackupId{}

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

	if id.LongTermRetentionBackupName, ok = parsed.Parsed["longTermRetentionBackupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "longTermRetentionBackupName", *parsed)
	}

	return &id, nil
}

// ParseLongTermRetentionDatabaseLongTermRetentionBackupIDInsensitively parses 'input' case-insensitively into a LongTermRetentionDatabaseLongTermRetentionBackupId
// note: this method should only be used for API response data and not user input
func ParseLongTermRetentionDatabaseLongTermRetentionBackupIDInsensitively(input string) (*LongTermRetentionDatabaseLongTermRetentionBackupId, error) {
	parser := resourceids.NewParserFromResourceIdType(LongTermRetentionDatabaseLongTermRetentionBackupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LongTermRetentionDatabaseLongTermRetentionBackupId{}

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

	if id.LongTermRetentionBackupName, ok = parsed.Parsed["longTermRetentionBackupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "longTermRetentionBackupName", *parsed)
	}

	return &id, nil
}

// ValidateLongTermRetentionDatabaseLongTermRetentionBackupID checks that 'input' can be parsed as a Long Term Retention Database Long Term Retention Backup ID
func ValidateLongTermRetentionDatabaseLongTermRetentionBackupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLongTermRetentionDatabaseLongTermRetentionBackupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Long Term Retention Database Long Term Retention Backup ID
func (id LongTermRetentionDatabaseLongTermRetentionBackupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql/locations/%s/longTermRetentionServers/%s/longTermRetentionDatabases/%s/longTermRetentionBackups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LocationName, id.LongTermRetentionServerName, id.LongTermRetentionDatabaseName, id.LongTermRetentionBackupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Long Term Retention Database Long Term Retention Backup ID
func (id LongTermRetentionDatabaseLongTermRetentionBackupId) Segments() []resourceids.Segment {
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
		resourceids.StaticSegment("staticLongTermRetentionBackups", "longTermRetentionBackups", "longTermRetentionBackups"),
		resourceids.UserSpecifiedSegment("longTermRetentionBackupName", "longTermRetentionBackupValue"),
	}
}

// String returns a human-readable description of this Long Term Retention Database Long Term Retention Backup ID
func (id LongTermRetentionDatabaseLongTermRetentionBackupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Long Term Retention Server Name: %q", id.LongTermRetentionServerName),
		fmt.Sprintf("Long Term Retention Database Name: %q", id.LongTermRetentionDatabaseName),
		fmt.Sprintf("Long Term Retention Backup Name: %q", id.LongTermRetentionBackupName),
	}
	return fmt.Sprintf("Long Term Retention Database Long Term Retention Backup (%s)", strings.Join(components, "\n"))
}
