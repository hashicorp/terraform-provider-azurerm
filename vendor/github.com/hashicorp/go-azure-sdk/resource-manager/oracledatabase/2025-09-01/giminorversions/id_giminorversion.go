package giminorversions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&GiMinorVersionId{})
}

var _ resourceids.ResourceId = &GiMinorVersionId{}

// GiMinorVersionId is a struct representing the Resource ID for a Gi Minor Version
type GiMinorVersionId struct {
	SubscriptionId     string
	LocationName       string
	GiVersionName      string
	GiMinorVersionName string
}

// NewGiMinorVersionID returns a new GiMinorVersionId struct
func NewGiMinorVersionID(subscriptionId string, locationName string, giVersionName string, giMinorVersionName string) GiMinorVersionId {
	return GiMinorVersionId{
		SubscriptionId:     subscriptionId,
		LocationName:       locationName,
		GiVersionName:      giVersionName,
		GiMinorVersionName: giMinorVersionName,
	}
}

// ParseGiMinorVersionID parses 'input' into a GiMinorVersionId
func ParseGiMinorVersionID(input string) (*GiMinorVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&GiMinorVersionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GiMinorVersionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseGiMinorVersionIDInsensitively parses 'input' case-insensitively into a GiMinorVersionId
// note: this method should only be used for API response data and not user input
func ParseGiMinorVersionIDInsensitively(input string) (*GiMinorVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&GiMinorVersionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GiMinorVersionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *GiMinorVersionId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.LocationName, ok = input.Parsed["locationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "locationName", input)
	}

	if id.GiVersionName, ok = input.Parsed["giVersionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "giVersionName", input)
	}

	if id.GiMinorVersionName, ok = input.Parsed["giMinorVersionName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "giMinorVersionName", input)
	}

	return nil
}

// ValidateGiMinorVersionID checks that 'input' can be parsed as a Gi Minor Version ID
func ValidateGiMinorVersionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseGiMinorVersionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Gi Minor Version ID
func (id GiMinorVersionId) ID() string {
	fmtString := "/subscriptions/%s/providers/Oracle.Database/locations/%s/giVersions/%s/giMinorVersions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.GiVersionName, id.GiMinorVersionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Gi Minor Version ID
func (id GiMinorVersionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticOracleDatabase", "Oracle.Database", "Oracle.Database"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
		resourceids.StaticSegment("staticGiVersions", "giVersions", "giVersions"),
		resourceids.UserSpecifiedSegment("giVersionName", "giVersionName"),
		resourceids.StaticSegment("staticGiMinorVersions", "giMinorVersions", "giMinorVersions"),
		resourceids.UserSpecifiedSegment("giMinorVersionName", "giMinorVersionName"),
	}
}

// String returns a human-readable description of this Gi Minor Version ID
func (id GiMinorVersionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Gi Version Name: %q", id.GiVersionName),
		fmt.Sprintf("Gi Minor Version Name: %q", id.GiMinorVersionName),
	}
	return fmt.Sprintf("Gi Minor Version (%s)", strings.Join(components, "\n"))
}
