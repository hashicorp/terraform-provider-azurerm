package dbversions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DbSystemDbVersionId{})
}

var _ resourceids.ResourceId = &DbSystemDbVersionId{}

// DbSystemDbVersionId is a struct representing the Resource ID for a Db System Db Version
type DbSystemDbVersionId struct {
	SubscriptionId        string
	LocationName          string
	DbSystemDbVersionName string
}

// NewDbSystemDbVersionID returns a new DbSystemDbVersionId struct
func NewDbSystemDbVersionID(subscriptionId string, locationName string, dbSystemDbVersionName string) DbSystemDbVersionId {
	return DbSystemDbVersionId{
		SubscriptionId:        subscriptionId,
		LocationName:          locationName,
		DbSystemDbVersionName: dbSystemDbVersionName,
	}
}

// ParseDbSystemDbVersionID parses 'input' into a DbSystemDbVersionId
func ParseDbSystemDbVersionID(input string) (*DbSystemDbVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DbSystemDbVersionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DbSystemDbVersionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDbSystemDbVersionIDInsensitively parses 'input' case-insensitively into a DbSystemDbVersionId
// note: this method should only be used for API response data and not user input
func ParseDbSystemDbVersionIDInsensitively(input string) (*DbSystemDbVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DbSystemDbVersionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DbSystemDbVersionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DbSystemDbVersionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.LocationName, ok = input.Parsed["locationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "locationName", input)
	}

	if id.DbSystemDbVersionName, ok = input.Parsed["dbSystemDbVersionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "dbSystemDbVersionName", input)
	}

	return nil
}

// ValidateDbSystemDbVersionID checks that 'input' can be parsed as a Db System Db Version ID
func ValidateDbSystemDbVersionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDbSystemDbVersionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Db System Db Version ID
func (id DbSystemDbVersionId) ID() string {
	fmtString := "/subscriptions/%s/providers/Oracle.Database/locations/%s/dbSystemDbVersions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.DbSystemDbVersionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Db System Db Version ID
func (id DbSystemDbVersionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticOracleDatabase", "Oracle.Database", "Oracle.Database"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
		resourceids.StaticSegment("staticDbSystemDbVersions", "dbSystemDbVersions", "dbSystemDbVersions"),
		resourceids.UserSpecifiedSegment("dbSystemDbVersionName", "dbSystemDbVersionName"),
	}
}

// String returns a human-readable description of this Db System Db Version ID
func (id DbSystemDbVersionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Db System Db Version Name: %q", id.DbSystemDbVersionName),
	}
	return fmt.Sprintf("Db System Db Version (%s)", strings.Join(components, "\n"))
}
