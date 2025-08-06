package dnsprivateviews

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DnsPrivateViewId{})
}

var _ resourceids.ResourceId = &DnsPrivateViewId{}

// DnsPrivateViewId is a struct representing the Resource ID for a Dns Private View
type DnsPrivateViewId struct {
	SubscriptionId     string
	LocationName       string
	DnsPrivateViewName string
}

// NewDnsPrivateViewID returns a new DnsPrivateViewId struct
func NewDnsPrivateViewID(subscriptionId string, locationName string, dnsPrivateViewName string) DnsPrivateViewId {
	return DnsPrivateViewId{
		SubscriptionId:     subscriptionId,
		LocationName:       locationName,
		DnsPrivateViewName: dnsPrivateViewName,
	}
}

// ParseDnsPrivateViewID parses 'input' into a DnsPrivateViewId
func ParseDnsPrivateViewID(input string) (*DnsPrivateViewId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DnsPrivateViewId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DnsPrivateViewId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDnsPrivateViewIDInsensitively parses 'input' case-insensitively into a DnsPrivateViewId
// note: this method should only be used for API response data and not user input
func ParseDnsPrivateViewIDInsensitively(input string) (*DnsPrivateViewId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DnsPrivateViewId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DnsPrivateViewId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DnsPrivateViewId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.LocationName, ok = input.Parsed["locationName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "locationName", input)
	}

	if id.DnsPrivateViewName, ok = input.Parsed["dnsPrivateViewName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "dnsPrivateViewName", input)
	}

	return nil
}

// ValidateDnsPrivateViewID checks that 'input' can be parsed as a Dns Private View ID
func ValidateDnsPrivateViewID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDnsPrivateViewID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Dns Private View ID
func (id DnsPrivateViewId) ID() string {
	fmtString := "/subscriptions/%s/providers/Oracle.Database/locations/%s/dnsPrivateViews/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.LocationName, id.DnsPrivateViewName)
}

// Segments returns a slice of Resource ID Segments which comprise this Dns Private View ID
func (id DnsPrivateViewId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticOracleDatabase", "Oracle.Database", "Oracle.Database"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationName"),
		resourceids.StaticSegment("staticDnsPrivateViews", "dnsPrivateViews", "dnsPrivateViews"),
		resourceids.UserSpecifiedSegment("dnsPrivateViewName", "dnsPrivateViewName"),
	}
}

// String returns a human-readable description of this Dns Private View ID
func (id DnsPrivateViewId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Location Name: %q", id.LocationName),
		fmt.Sprintf("Dns Private View Name: %q", id.DnsPrivateViewName),
	}
	return fmt.Sprintf("Dns Private View (%s)", strings.Join(components, "\n"))
}
