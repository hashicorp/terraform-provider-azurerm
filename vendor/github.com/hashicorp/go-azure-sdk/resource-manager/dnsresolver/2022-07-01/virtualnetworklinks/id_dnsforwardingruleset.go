package virtualnetworklinks

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = DnsForwardingRulesetId{}

// DnsForwardingRulesetId is a struct representing the Resource ID for a Dns Forwarding Ruleset
type DnsForwardingRulesetId struct {
	SubscriptionId           string
	ResourceGroupName        string
	DnsForwardingRulesetName string
}

// NewDnsForwardingRulesetID returns a new DnsForwardingRulesetId struct
func NewDnsForwardingRulesetID(subscriptionId string, resourceGroupName string, dnsForwardingRulesetName string) DnsForwardingRulesetId {
	return DnsForwardingRulesetId{
		SubscriptionId:           subscriptionId,
		ResourceGroupName:        resourceGroupName,
		DnsForwardingRulesetName: dnsForwardingRulesetName,
	}
}

// ParseDnsForwardingRulesetID parses 'input' into a DnsForwardingRulesetId
func ParseDnsForwardingRulesetID(input string) (*DnsForwardingRulesetId, error) {
	parser := resourceids.NewParserFromResourceIdType(DnsForwardingRulesetId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DnsForwardingRulesetId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DnsForwardingRulesetName, ok = parsed.Parsed["dnsForwardingRulesetName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "dnsForwardingRulesetName", *parsed)
	}

	return &id, nil
}

// ParseDnsForwardingRulesetIDInsensitively parses 'input' case-insensitively into a DnsForwardingRulesetId
// note: this method should only be used for API response data and not user input
func ParseDnsForwardingRulesetIDInsensitively(input string) (*DnsForwardingRulesetId, error) {
	parser := resourceids.NewParserFromResourceIdType(DnsForwardingRulesetId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DnsForwardingRulesetId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DnsForwardingRulesetName, ok = parsed.Parsed["dnsForwardingRulesetName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "dnsForwardingRulesetName", *parsed)
	}

	return &id, nil
}

// ValidateDnsForwardingRulesetID checks that 'input' can be parsed as a Dns Forwarding Ruleset ID
func ValidateDnsForwardingRulesetID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDnsForwardingRulesetID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Dns Forwarding Ruleset ID
func (id DnsForwardingRulesetId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/dnsForwardingRulesets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DnsForwardingRulesetName)
}

// Segments returns a slice of Resource ID Segments which comprise this Dns Forwarding Ruleset ID
func (id DnsForwardingRulesetId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticDnsForwardingRulesets", "dnsForwardingRulesets", "dnsForwardingRulesets"),
		resourceids.UserSpecifiedSegment("dnsForwardingRulesetName", "dnsForwardingRulesetValue"),
	}
}

// String returns a human-readable description of this Dns Forwarding Ruleset ID
func (id DnsForwardingRulesetId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Dns Forwarding Ruleset Name: %q", id.DnsForwardingRulesetName),
	}
	return fmt.Sprintf("Dns Forwarding Ruleset (%s)", strings.Join(components, "\n"))
}
