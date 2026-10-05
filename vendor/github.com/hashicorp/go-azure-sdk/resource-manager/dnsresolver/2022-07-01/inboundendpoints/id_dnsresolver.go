package inboundendpoints

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/recaser"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

func init() {
	recaser.RegisterResourceId(&DnsResolverId{})
}

var _ resourceids.ResourceId = &DnsResolverId{}

// DnsResolverId is a struct representing the Resource ID for a Dns Resolver
type DnsResolverId struct {
	SubscriptionId    string
	ResourceGroupName string
	DnsResolverName   string
}

// NewDnsResolverID returns a new DnsResolverId struct
func NewDnsResolverID(subscriptionId string, resourceGroupName string, dnsResolverName string) DnsResolverId {
	return DnsResolverId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		DnsResolverName:   dnsResolverName,
	}
}

// ParseDnsResolverID parses 'input' into a DnsResolverId
func ParseDnsResolverID(input string) (*DnsResolverId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DnsResolverId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DnsResolverId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

// ParseDnsResolverIDInsensitively parses 'input' case-insensitively into a DnsResolverId
// note: this method should only be used for API response data and not user input
func ParseDnsResolverIDInsensitively(input string) (*DnsResolverId, error) {
	parser := resourceids.NewParserFromResourceIdType(&DnsResolverId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	id := DnsResolverId{}
	if err = id.FromParseResult(*parsed); err != nil {
		return nil, err
	}

	return &id, nil
}

func (id *DnsResolverId) FromParseResult(input resourceids.ParseResult) error {
	var ok bool

	if id.SubscriptionId, ok = input.Parsed["subscriptionId"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", input)
	}

	if id.ResourceGroupName, ok = input.Parsed["resourceGroupName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", input)
	}

	if id.DnsResolverName, ok = input.Parsed["dnsResolverName"]; !ok {
		return resourceids.NewSegmentNotSpecifiedError(id, "dnsResolverName", input)
	}

	return nil
}

// ValidateDnsResolverID checks that 'input' can be parsed as a Dns Resolver ID
func ValidateDnsResolverID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDnsResolverID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Dns Resolver ID
func (id DnsResolverId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/dnsResolvers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DnsResolverName)
}

// Segments returns a slice of Resource ID Segments which comprise this Dns Resolver ID
func (id DnsResolverId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticDnsResolvers", "dnsResolvers", "dnsResolvers"),
		resourceids.UserSpecifiedSegment("dnsResolverName", "dnsResolverName"),
	}
}

// String returns a human-readable description of this Dns Resolver ID
func (id DnsResolverId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Dns Resolver Name: %q", id.DnsResolverName),
	}
	return fmt.Sprintf("Dns Resolver (%s)", strings.Join(components, "\n"))
}
