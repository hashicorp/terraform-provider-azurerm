package giversions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&GiVersionId{})
}

var _ resourceids.ResourceId = &GiVersionId{}

// GiVersionId is a struct representing the Resource ID for a Gi Version
type GiVersionId struct {
	SubscriptionId string
	LocationName   string
	GiVersionName  string
}

// NewGiVersionID returns a new GiVersionId struct
func NewGiVersionID(subscriptionId string, locationName string, giVersionName string) GiVersionId {
	return GiVersionId{
		SubscriptionId: subscriptionId,
		LocationName:   locationName,
		GiVersionName:  giVersionName,
	}
}

// ParseGiVersionID parses 'input' into a GiVersionId
func ParseGiVersionID(input string) (*GiVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&GiVersionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GiVersionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseGiVersionIDInsensitively parses 'input' case-insensitively into a GiVersionId
// note: this method should only be used for API response data and not user input
func ParseGiVersionIDInsensitively(input string) (*GiVersionId, error) {
	parser := resourceids.NewParserFromResourceIdType(&GiVersionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := GiVersionId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *GiVersionId) FromParseResult(input resourceids.ParseResult) error {
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

	return nil
}

// ValidateGiVersionID checks that 'input' can be parsed as a Gi Version ID
func ValidateGiVersionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseGiVersionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Gi Version ID
func (id GiVersionId) ID() string {
	fmtString := "/subscriptions/%s/providers/Oracle.Database/locations/%s/giVersions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.GiVersionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Gi Version ID
func (id GiVersionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticOracleDatabase", "Oracle.Database", "Oracle.Database"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
		resourceids.StaticSegment("staticGiVersions", "giVersions", "giVersions"),
		resourceids.UserSpecifiedSegment("giVersionName", "giVersionName"),
	}
}

// String returns a human-readable description of this Gi Version ID
func (id GiVersionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Gi Version Name: %q", id.GiVersionName),
	}
	return fmt.Sprintf("Gi Version (%s)", strings.Join(components, "\n"))
}
