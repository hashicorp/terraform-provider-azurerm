package systemversions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&SystemVersionId{})
}

var _ resourceids.ResourceId = &SystemVersionId{}

// SystemVersionId is a struct representing the Resource ID for a System Version
type SystemVersionId struct {
	SubscriptionId    string
	LocationName      string
	SystemVersionName string
}

// NewSystemVersionID returns a new SystemVersionId struct
func NewSystemVersionID(subscriptionId string, locationName string, systemVersionName string) SystemVersionId {
	return SystemVersionId{
		SubscriptionId:    subscriptionId,
		LocationName:      locationName,
		SystemVersionName: systemVersionName,
	}
}

// ParseSystemVersionID parses 'input' into a SystemVersionId
func ParseSystemVersionID(input string) (*SystemVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SystemVersionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SystemVersionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseSystemVersionIDInsensitively parses 'input' case-insensitively into a SystemVersionId
// note: this method should only be used for API response data and not user input
func ParseSystemVersionIDInsensitively(input string) (*SystemVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&SystemVersionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := SystemVersionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *SystemVersionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.LocationName, ok = input.Parsed["locationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "locationName", input)
	}

	if id.SystemVersionName, ok = input.Parsed["systemVersionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "systemVersionName", input)
	}

	return nil
}

// ValidateSystemVersionID checks that 'input' can be parsed as a System Version ID
func ValidateSystemVersionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSystemVersionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted System Version ID
func (id SystemVersionId) ID() string {
	fmtString := "/subscriptions/%s/providers/Oracle.Database/locations/%s/systemVersions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.SystemVersionName)
}

// Segments returns a slice of Resource ID Segments which comprise this System Version ID
func (id SystemVersionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticOracleDatabase", "Oracle.Database", "Oracle.Database"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
		resourceids.StaticSegment("staticSystemVersions", "systemVersions", "systemVersions"),
		resourceids.UserSpecifiedSegment("systemVersionName", "systemVersionName"),
	}
}

// String returns a human-readable description of this System Version ID
func (id SystemVersionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("System Version Name: %q", id.SystemVersionName),
	}
	return fmt.Sprintf("System Version (%s)", strings.Join(components, "\n"))
}
