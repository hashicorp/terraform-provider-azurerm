package dbsystemshapes

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DbSystemShapeId{})
}

var _ resourceids.ResourceId = &DbSystemShapeId{}

// DbSystemShapeId is a struct representing the Resource ID for a Db System Shape
type DbSystemShapeId struct {
	SubscriptionId    string
	LocationName      string
	DbSystemShapeName string
}

// NewDbSystemShapeID returns a new DbSystemShapeId struct
func NewDbSystemShapeID(subscriptionId string, locationName string, dbSystemShapeName string) DbSystemShapeId {
	return DbSystemShapeId{
		SubscriptionId:    subscriptionId,
		LocationName:      locationName,
		DbSystemShapeName: dbSystemShapeName,
	}
}

// ParseDbSystemShapeID parses 'input' into a DbSystemShapeId
func ParseDbSystemShapeID(input string) (*DbSystemShapeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DbSystemShapeId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DbSystemShapeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDbSystemShapeIDInsensitively parses 'input' case-insensitively into a DbSystemShapeId
// note: this method should only be used for API response data and not user input
func ParseDbSystemShapeIDInsensitively(input string) (*DbSystemShapeId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DbSystemShapeId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DbSystemShapeId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DbSystemShapeId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.LocationName, ok = input.Parsed["locationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "locationName", input)
	}

	if id.DbSystemShapeName, ok = input.Parsed["dbSystemShapeName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "dbSystemShapeName", input)
	}

	return nil
}

// ValidateDbSystemShapeID checks that 'input' can be parsed as a Db System Shape ID
func ValidateDbSystemShapeID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDbSystemShapeID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Db System Shape ID
func (id DbSystemShapeId) ID() string {
	fmtString := "/subscriptions/%s/providers/Oracle.Database/locations/%s/dbSystemShapes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.DbSystemShapeName)
}

// Segments returns a slice of Resource ID Segments which comprise this Db System Shape ID
func (id DbSystemShapeId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticOracleDatabase", "Oracle.Database", "Oracle.Database"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
		resourceids.StaticSegment("staticDbSystemShapes", "dbSystemShapes", "dbSystemShapes"),
		resourceids.UserSpecifiedSegment("dbSystemShapeName", "dbSystemShapeName"),
	}
}

// String returns a human-readable description of this Db System Shape ID
func (id DbSystemShapeId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Db System Shape Name: %q", id.DbSystemShapeName),
	}
	return fmt.Sprintf("Db System Shape (%s)", strings.Join(components, "\n"))
}
