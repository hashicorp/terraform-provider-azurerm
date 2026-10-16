package databaseeditions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DatabaseEditionId{})
}

var _ resourceids.ResourceId = &DatabaseEditionId{}

// DatabaseEditionId is a struct representing the Resource ID for a Database Edition
type DatabaseEditionId struct {
	SubscriptionId      string
	LocationName        string
	DatabaseEditionName string
}

// NewDatabaseEditionID returns a new DatabaseEditionId struct
func NewDatabaseEditionID(subscriptionId string, locationName string, databaseEditionName string) DatabaseEditionId {
	return DatabaseEditionId{
		SubscriptionId:      subscriptionId,
		LocationName:        locationName,
		DatabaseEditionName: databaseEditionName,
	}
}

// ParseDatabaseEditionID parses 'input' into a DatabaseEditionId
func ParseDatabaseEditionID(input string) (*DatabaseEditionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DatabaseEditionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DatabaseEditionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDatabaseEditionIDInsensitively parses 'input' case-insensitively into a DatabaseEditionId
// note: this method should only be used for API response data and not user input
func ParseDatabaseEditionIDInsensitively(input string) (*DatabaseEditionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DatabaseEditionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DatabaseEditionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DatabaseEditionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.LocationName, ok = input.Parsed["locationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "locationName", input)
	}

	if id.DatabaseEditionName, ok = input.Parsed["databaseEditionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "databaseEditionName", input)
	}

	return nil
}

// ValidateDatabaseEditionID checks that 'input' can be parsed as a Database Edition ID
func ValidateDatabaseEditionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDatabaseEditionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Database Edition ID
func (id DatabaseEditionId) ID() string {
	fmtString := "/subscriptions/%s/providers/Oracle.Database/locations/%s/databaseEditions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.DatabaseEditionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Database Edition ID
func (id DatabaseEditionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticOracleDatabase", "Oracle.Database", "Oracle.Database"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
		resourceids.StaticSegment("staticDatabaseEditions", "databaseEditions", "databaseEditions"),
		resourceids.UserSpecifiedSegment("databaseEditionName", "databaseEditionName"),
	}
}

// String returns a human-readable description of this Database Edition ID
func (id DatabaseEditionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Database Edition Name: %q", id.DatabaseEditionName),
	}
	return fmt.Sprintf("Database Edition (%s)", strings.Join(components, "\n"))
}
