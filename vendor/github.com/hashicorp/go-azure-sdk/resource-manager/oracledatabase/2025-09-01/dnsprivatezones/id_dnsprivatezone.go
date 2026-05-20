package dnsprivatezones

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DnsPrivateZoneId{})
}

var _ resourceids.ResourceId = &DnsPrivateZoneId{}

// DnsPrivateZoneId is a struct representing the Resource ID for a Dns Private Zone
type DnsPrivateZoneId struct {
	SubscriptionId     string
	LocationName       string
	DnsPrivateZoneName string
}

// NewDnsPrivateZoneID returns a new DnsPrivateZoneId struct
func NewDnsPrivateZoneID(subscriptionId string, locationName string, dnsPrivateZoneName string) DnsPrivateZoneId {
	return DnsPrivateZoneId{
		SubscriptionId:     subscriptionId,
		LocationName:       locationName,
		DnsPrivateZoneName: dnsPrivateZoneName,
	}
}

// ParseDnsPrivateZoneID parses 'input' into a DnsPrivateZoneId
func ParseDnsPrivateZoneID(input string) (*DnsPrivateZoneId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DnsPrivateZoneId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DnsPrivateZoneId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDnsPrivateZoneIDInsensitively parses 'input' case-insensitively into a DnsPrivateZoneId
// note: this method should only be used for API response data and not user input
func ParseDnsPrivateZoneIDInsensitively(input string) (*DnsPrivateZoneId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DnsPrivateZoneId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DnsPrivateZoneId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DnsPrivateZoneId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.LocationName, ok = input.Parsed["locationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "locationName", input)
	}

	if id.DnsPrivateZoneName, ok = input.Parsed["dnsPrivateZoneName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "dnsPrivateZoneName", input)
	}

	return nil
}

// ValidateDnsPrivateZoneID checks that 'input' can be parsed as a Dns Private Zone ID
func ValidateDnsPrivateZoneID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDnsPrivateZoneID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Dns Private Zone ID
func (id DnsPrivateZoneId) ID() string {
	fmtString := "/subscriptions/%s/providers/Oracle.Database/locations/%s/dnsPrivateZones/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.DnsPrivateZoneName)
}

// Segments returns a slice of Resource ID Segments which comprise this Dns Private Zone ID
func (id DnsPrivateZoneId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticOracleDatabase", "Oracle.Database", "Oracle.Database"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
		resourceids.StaticSegment("staticDnsPrivateZones", "dnsPrivateZones", "dnsPrivateZones"),
		resourceids.UserSpecifiedSegment("dnsPrivateZoneName", "dnsPrivateZoneName"),
	}
}

// String returns a human-readable description of this Dns Private Zone ID
func (id DnsPrivateZoneId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Dns Private Zone Name: %q", id.DnsPrivateZoneName),
	}
	return fmt.Sprintf("Dns Private Zone (%s)", strings.Join(components, "\n"))
}
