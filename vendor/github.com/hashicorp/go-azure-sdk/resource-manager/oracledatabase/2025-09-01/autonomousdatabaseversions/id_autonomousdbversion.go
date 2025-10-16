package autonomousdatabaseversions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AutonomousDbVersionId{})
}

var _ resourceids.ResourceId = &AutonomousDbVersionId{}

// AutonomousDbVersionId is a struct representing the Resource ID for a Autonomous Db Version
type AutonomousDbVersionId struct {
	SubscriptionId          string
	LocationName            string
	AutonomousDbVersionName string
}

// NewAutonomousDbVersionID returns a new AutonomousDbVersionId struct
func NewAutonomousDbVersionID(subscriptionId string, locationName string, autonomousDbVersionName string) AutonomousDbVersionId {
	return AutonomousDbVersionId{
		SubscriptionId:          subscriptionId,
		LocationName:            locationName,
		AutonomousDbVersionName: autonomousDbVersionName,
	}
}

// ParseAutonomousDbVersionID parses 'input' into a AutonomousDbVersionId
func ParseAutonomousDbVersionID(input string) (*AutonomousDbVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AutonomousDbVersionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AutonomousDbVersionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAutonomousDbVersionIDInsensitively parses 'input' case-insensitively into a AutonomousDbVersionId
// note: this method should only be used for API response data and not user input
func ParseAutonomousDbVersionIDInsensitively(input string) (*AutonomousDbVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AutonomousDbVersionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AutonomousDbVersionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AutonomousDbVersionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.LocationName, ok = input.Parsed["locationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "locationName", input)
	}

	if id.AutonomousDbVersionName, ok = input.Parsed["autonomousDbVersionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "autonomousDbVersionName", input)
	}

	return nil
}

// ValidateAutonomousDbVersionID checks that 'input' can be parsed as a Autonomous Db Version ID
func ValidateAutonomousDbVersionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAutonomousDbVersionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Autonomous Db Version ID
func (id AutonomousDbVersionId) ID() string {
	fmtString := "/subscriptions/%s/providers/Oracle.Database/locations/%s/autonomousDbVersions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.AutonomousDbVersionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Autonomous Db Version ID
func (id AutonomousDbVersionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticOracleDatabase", "Oracle.Database", "Oracle.Database"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
		resourceids.StaticSegment("staticAutonomousDbVersions", "autonomousDbVersions", "autonomousDbVersions"),
		resourceids.UserSpecifiedSegment("autonomousDbVersionName", "autonomousDbVersionName"),
	}
}

// String returns a human-readable description of this Autonomous Db Version ID
func (id AutonomousDbVersionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Autonomous Db Version Name: %q", id.AutonomousDbVersionName),
	}
	return fmt.Sprintf("Autonomous Db Version (%s)", strings.Join(components, "\n"))
}
