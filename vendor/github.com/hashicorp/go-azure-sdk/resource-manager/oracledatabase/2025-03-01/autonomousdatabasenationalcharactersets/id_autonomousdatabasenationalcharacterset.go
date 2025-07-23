package autonomousdatabasenationalcharactersets

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&AutonomousDatabaseNationalCharacterSetId{})
}

var _ resourceids.ResourceId = &AutonomousDatabaseNationalCharacterSetId{}

// AutonomousDatabaseNationalCharacterSetId is a struct representing the Resource ID for a Autonomous Database National Character Set
type AutonomousDatabaseNationalCharacterSetId struct {
	SubscriptionId                             string
	LocationName                               string
	AutonomousDatabaseNationalCharacterSetName string
}

// NewAutonomousDatabaseNationalCharacterSetID returns a new AutonomousDatabaseNationalCharacterSetId struct
func NewAutonomousDatabaseNationalCharacterSetID(subscriptionId string, locationName string, autonomousDatabaseNationalCharacterSetName string) AutonomousDatabaseNationalCharacterSetId {
	return AutonomousDatabaseNationalCharacterSetId{
		SubscriptionId: subscriptionId,
		LocationName:   locationName,
		AutonomousDatabaseNationalCharacterSetName: autonomousDatabaseNationalCharacterSetName,
	}
}

// ParseAutonomousDatabaseNationalCharacterSetID parses 'input' into a AutonomousDatabaseNationalCharacterSetId
func ParseAutonomousDatabaseNationalCharacterSetID(input string) (*AutonomousDatabaseNationalCharacterSetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AutonomousDatabaseNationalCharacterSetId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AutonomousDatabaseNationalCharacterSetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseAutonomousDatabaseNationalCharacterSetIDInsensitively parses 'input' case-insensitively into a AutonomousDatabaseNationalCharacterSetId
// note: this method should only be used for API response data and not user input
func ParseAutonomousDatabaseNationalCharacterSetIDInsensitively(input string) (*AutonomousDatabaseNationalCharacterSetId, error) {
	parser := resourceids.NewParserFromResourceIdType(&AutonomousDatabaseNationalCharacterSetId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := AutonomousDatabaseNationalCharacterSetId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *AutonomousDatabaseNationalCharacterSetId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.LocationName, ok = input.Parsed["locationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "locationName", input)
	}

	if id.AutonomousDatabaseNationalCharacterSetName, ok = input.Parsed["autonomousDatabaseNationalCharacterSetName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "autonomousDatabaseNationalCharacterSetName", input)
	}

	return nil
}

// ValidateAutonomousDatabaseNationalCharacterSetID checks that 'input' can be parsed as a Autonomous Database National Character Set ID
func ValidateAutonomousDatabaseNationalCharacterSetID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAutonomousDatabaseNationalCharacterSetID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Autonomous Database National Character Set ID
func (id AutonomousDatabaseNationalCharacterSetId) ID() string {
	fmtString := "/subscriptions/%s/providers/Oracle.Database/locations/%s/autonomousDatabaseNationalCharacterSets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.AutonomousDatabaseNationalCharacterSetName)
}

// Segments returns a slice of Resource ID Segments which comprise this Autonomous Database National Character Set ID
func (id AutonomousDatabaseNationalCharacterSetId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticOracleDatabase", "Oracle.Database", "Oracle.Database"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
		resourceids.StaticSegment("staticAutonomousDatabaseNationalCharacterSets", "autonomousDatabaseNationalCharacterSets", "autonomousDatabaseNationalCharacterSets"),
		resourceids.UserSpecifiedSegment("autonomousDatabaseNationalCharacterSetName", "autonomousDatabaseNationalCharacterSetName"),
	}
}

// String returns a human-readable description of this Autonomous Database National Character Set ID
func (id AutonomousDatabaseNationalCharacterSetId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Autonomous Database National Character Set Name: %q", id.AutonomousDatabaseNationalCharacterSetName),
	}
	return fmt.Sprintf("Autonomous Database National Character Set (%s)", strings.Join(components, "\n"))
}
