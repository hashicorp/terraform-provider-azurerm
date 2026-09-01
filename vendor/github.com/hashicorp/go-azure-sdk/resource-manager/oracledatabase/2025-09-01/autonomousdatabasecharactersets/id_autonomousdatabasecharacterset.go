package autonomousdatabasecharactersets

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AutonomousDatabaseCharacterSetId{})
}

var _ resourceids.ResourceId = &AutonomousDatabaseCharacterSetId{}

// AutonomousDatabaseCharacterSetId is a struct representing the Resource ID for a Autonomous Database Character Set
type AutonomousDatabaseCharacterSetId struct {
	SubscriptionId                     string
	LocationName                       string
	AutonomousDatabaseCharacterSetName string
}

// NewAutonomousDatabaseCharacterSetID returns a new AutonomousDatabaseCharacterSetId struct
func NewAutonomousDatabaseCharacterSetID(subscriptionId string, locationName string, autonomousDatabaseCharacterSetName string) AutonomousDatabaseCharacterSetId {
	return AutonomousDatabaseCharacterSetId{
		SubscriptionId:                     subscriptionId,
		LocationName:                       locationName,
		AutonomousDatabaseCharacterSetName: autonomousDatabaseCharacterSetName,
	}
}

// ParseAutonomousDatabaseCharacterSetID parses 'input' into a AutonomousDatabaseCharacterSetId
func ParseAutonomousDatabaseCharacterSetID(input string) (*AutonomousDatabaseCharacterSetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AutonomousDatabaseCharacterSetId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AutonomousDatabaseCharacterSetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAutonomousDatabaseCharacterSetIDInsensitively parses 'input' case-insensitively into a AutonomousDatabaseCharacterSetId
// note: this method should only be used for API response data and not user input
func ParseAutonomousDatabaseCharacterSetIDInsensitively(input string) (*AutonomousDatabaseCharacterSetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AutonomousDatabaseCharacterSetId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AutonomousDatabaseCharacterSetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AutonomousDatabaseCharacterSetId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.LocationName, ok = input.Parsed["locationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "locationName", input)
	}

	if id.AutonomousDatabaseCharacterSetName, ok = input.Parsed["autonomousDatabaseCharacterSetName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "autonomousDatabaseCharacterSetName", input)
	}

	return nil
}

// ValidateAutonomousDatabaseCharacterSetID checks that 'input' can be parsed as a Autonomous Database Character Set ID
func ValidateAutonomousDatabaseCharacterSetID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAutonomousDatabaseCharacterSetID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Autonomous Database Character Set ID
func (id AutonomousDatabaseCharacterSetId) ID() string {
	fmtString := "/subscriptions/%s/providers/Oracle.Database/locations/%s/autonomousDatabaseCharacterSets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.AutonomousDatabaseCharacterSetName)
}

// Segments returns a slice of Resource ID Segments which comprise this Autonomous Database Character Set ID
func (id AutonomousDatabaseCharacterSetId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticOracleDatabase", "Oracle.Database", "Oracle.Database"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
		resourceids.StaticSegment("staticAutonomousDatabaseCharacterSets", "autonomousDatabaseCharacterSets", "autonomousDatabaseCharacterSets"),
		resourceids.UserSpecifiedSegment("autonomousDatabaseCharacterSetName", "autonomousDatabaseCharacterSetName"),
	}
}

// String returns a human-readable description of this Autonomous Database Character Set ID
func (id AutonomousDatabaseCharacterSetId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Autonomous Database Character Set Name: %q", id.AutonomousDatabaseCharacterSetName),
	}
	return fmt.Sprintf("Autonomous Database Character Set (%s)", strings.Join(components, "\n"))
}
