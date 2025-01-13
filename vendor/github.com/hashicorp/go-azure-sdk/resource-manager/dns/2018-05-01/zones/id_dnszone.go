package zones

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DnsZoneId{})
}

var _ resourceids.ResourceId = &DnsZoneId{}

// DnsZoneId is a struct representing the Resource ID for a Dns Zone
type DnsZoneId struct {
	SubscriptionId    string
	ResourceGroupName string
	DnsZoneName       string
}

// NewDnsZoneID returns a new DnsZoneId struct
func NewDnsZoneID(subscriptionId string, resourceGroupName string, dnsZoneName string) DnsZoneId {
	return DnsZoneId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		DnsZoneName:       dnsZoneName,
	}
}

// ParseDnsZoneID parses 'input' into a DnsZoneId
func ParseDnsZoneID(input string) (*DnsZoneId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DnsZoneId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DnsZoneId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDnsZoneIDInsensitively parses 'input' case-insensitively into a DnsZoneId
// note: this method should only be used for API response data and not user input
func ParseDnsZoneIDInsensitively(input string) (*DnsZoneId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DnsZoneId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DnsZoneId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DnsZoneId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.DnsZoneName, ok = input.Parsed["dnsZoneName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "dnsZoneName", input)
	}

	return nil
}

// ValidateDnsZoneID checks that 'input' can be parsed as a Dns Zone ID
func ValidateDnsZoneID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDnsZoneID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Dns Zone ID
func (id DnsZoneId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/dnsZones/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DnsZoneName)
}

// Segments returns a slice of Resource ID Segments which comprise this Dns Zone ID
func (id DnsZoneId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticDnsZones", "dnsZones", "dnsZones"),
		resourceids.UserSpecifiedSegment("dnsZoneName", "dnsZoneName"),
	}
}

// String returns a human-readable description of this Dns Zone ID
func (id DnsZoneId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Dns Zone Name: %q", id.DnsZoneName),
	}
	return fmt.Sprintf("Dns Zone (%s)", strings.Join(components, "\n"))
}
