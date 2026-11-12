package databasesystemshaperesources

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DatabaseSystemShapeId{})
}

var _ resourceids.ResourceId = &DatabaseSystemShapeId{}

// DatabaseSystemShapeId is a struct representing the Resource ID for a Database System Shape
type DatabaseSystemShapeId struct {
	SubscriptionId          string
	LocationName            string
	DatabaseSystemShapeName string
}

// NewDatabaseSystemShapeID returns a new DatabaseSystemShapeId struct
func NewDatabaseSystemShapeID(subscriptionId string, locationName string, databaseSystemShapeName string) DatabaseSystemShapeId {
	return DatabaseSystemShapeId{
		SubscriptionId:          subscriptionId,
		LocationName:            locationName,
		DatabaseSystemShapeName: databaseSystemShapeName,
	}
}

// ParseDatabaseSystemShapeID parses 'input' into a DatabaseSystemShapeId
func ParseDatabaseSystemShapeID(input string) (*DatabaseSystemShapeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DatabaseSystemShapeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DatabaseSystemShapeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDatabaseSystemShapeIDInsensitively parses 'input' case-insensitively into a DatabaseSystemShapeId
// note: this method should only be used for API response data and not user input
func ParseDatabaseSystemShapeIDInsensitively(input string) (*DatabaseSystemShapeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DatabaseSystemShapeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DatabaseSystemShapeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DatabaseSystemShapeId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.LocationName, ok = input.Parsed["locationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "locationName", input)
	}

	if id.DatabaseSystemShapeName, ok = input.Parsed["databaseSystemShapeName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "databaseSystemShapeName", input)
	}

	return nil
}

// ValidateDatabaseSystemShapeID checks that 'input' can be parsed as a Database System Shape ID
func ValidateDatabaseSystemShapeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDatabaseSystemShapeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Database System Shape ID
func (id DatabaseSystemShapeId) ID() string {
	fmtString := "/subscriptions/%s/providers/Oracle.Database/locations/%s/databaseSystemShapes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.DatabaseSystemShapeName)
}

// Segments returns a slice of Resource ID Segments which comprise this Database System Shape ID
func (id DatabaseSystemShapeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticOracleDatabase", "Oracle.Database", "Oracle.Database"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
		resourceids.StaticSegment("staticDatabaseSystemShapes", "databaseSystemShapes", "databaseSystemShapes"),
		resourceids.UserSpecifiedSegment("databaseSystemShapeName", "databaseSystemShapeName"),
	}
}

// String returns a human-readable description of this Database System Shape ID
func (id DatabaseSystemShapeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Database System Shape Name: %q", id.DatabaseSystemShapeName),
	}
	return fmt.Sprintf("Database System Shape (%s)", strings.Join(components, "\n"))
}
