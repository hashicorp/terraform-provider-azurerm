package firewalls

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&FirewallId{})
}

var _ resourceids.ResourceId = &FirewallId{}

// FirewallId is a struct representing the Resource ID for a Firewall
type FirewallId struct {
	SubscriptionId    string
	ResourceGroupName string
	FirewallName      string
}

// NewFirewallID returns a new FirewallId struct
func NewFirewallID(subscriptionId string, resourceGroupName string, firewallName string) FirewallId {
	return FirewallId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		FirewallName:      firewallName,
	}
}

// ParseFirewallID parses 'input' into a FirewallId
func ParseFirewallID(input string) (*FirewallId, error) {
	parser := resourceids.NewParserFromResourceIdType(&FirewallId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := FirewallId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseFirewallIDInsensitively parses 'input' case-insensitively into a FirewallId
// note: this method should only be used for API response data and not user input
func ParseFirewallIDInsensitively(input string) (*FirewallId, error) {
	parser := resourceids.NewParserFromResourceIdType(&FirewallId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := FirewallId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *FirewallId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.FirewallName, ok = input.Parsed["firewallName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "firewallName", input)
	}

	return nil
}

// ValidateFirewallID checks that 'input' can be parsed as a Firewall ID
func ValidateFirewallID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseFirewallID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Firewall ID
func (id FirewallId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/PaloAltoNetworks.Cloudngfw/firewalls/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.FirewallName)
}

// Segments returns a slice of Resource ID Segments which comprise this Firewall ID
func (id FirewallId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticPaloAltoNetworksCloudngfw", "PaloAltoNetworks.Cloudngfw", "PaloAltoNetworks.Cloudngfw"),
		resourceids.StaticSegment("staticFirewalls", "firewalls", "firewalls"),
		resourceids.UserSpecifiedSegment("firewallName", "firewallName"),
	}
}

// String returns a human-readable description of this Firewall ID
func (id FirewallId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Firewall Name: %q", id.FirewallName),
	}
	return fmt.Sprintf("Firewall (%s)", strings.Join(components, "\n"))
}
