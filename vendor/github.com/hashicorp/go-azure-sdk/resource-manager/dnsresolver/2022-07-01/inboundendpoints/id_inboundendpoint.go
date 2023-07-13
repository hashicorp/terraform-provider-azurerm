package inboundendpoints

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ resourceids.ResourceId = InboundEndpointId{}

// InboundEndpointId is a struct representing the Resource ID for a Inbound Endpoint
type InboundEndpointId struct {
	SubscriptionId      string
	ResourceGroupName   string
	DnsResolverName     string
	InboundEndpointName string
}

// NewInboundEndpointID returns a new InboundEndpointId struct
func NewInboundEndpointID(subscriptionId string, resourceGroupName string, dnsResolverName string, inboundEndpointName string) InboundEndpointId {
	return InboundEndpointId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		DnsResolverName:     dnsResolverName,
		InboundEndpointName: inboundEndpointName,
	}
}

// ParseInboundEndpointID parses 'input' into a InboundEndpointId
func ParseInboundEndpointID(input string) (*InboundEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(InboundEndpointId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := InboundEndpointId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DnsResolverName, ok = parsed.Parsed["dnsResolverName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "dnsResolverName", *parsed)
	}

	if id.InboundEndpointName, ok = parsed.Parsed["inboundEndpointName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "inboundEndpointName", *parsed)
	}

	return &id, nil
}

// ParseInboundEndpointIDInsensitively parses 'input' case-insensitively into a InboundEndpointId
// note: this method should only be used for API response data and not user input
func ParseInboundEndpointIDInsensitively(input string) (*InboundEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(InboundEndpointId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := InboundEndpointId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "subscriptionId", *parsed)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "resourceGroupName", *parsed)
	}

	if id.DnsResolverName, ok = parsed.Parsed["dnsResolverName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "dnsResolverName", *parsed)
	}

	if id.InboundEndpointName, ok = parsed.Parsed["inboundEndpointName"]; !ok {
		return nil, resourceids.NewSegmentNotSpecifiedError(id, "inboundEndpointName", *parsed)
	}

	return &id, nil
}

// ValidateInboundEndpointID checks that 'input' can be parsed as a Inbound Endpoint ID
func ValidateInboundEndpointID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseInboundEndpointID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Inbound Endpoint ID
func (id InboundEndpointId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/dnsResolvers/%s/inboundEndpoints/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DnsResolverName, id.InboundEndpointName)
}

// Segments returns a slice of Resource ID Segments which comprise this Inbound Endpoint ID
func (id InboundEndpointId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticDnsResolvers", "dnsResolvers", "dnsResolvers"),
		resourceids.UserSpecifiedSegment("dnsResolverName", "dnsResolverValue"),
		resourceids.StaticSegment("staticInboundEndpoints", "inboundEndpoints", "inboundEndpoints"),
		resourceids.UserSpecifiedSegment("inboundEndpointName", "inboundEndpointValue"),
	}
}

// String returns a human-readable description of this Inbound Endpoint ID
func (id InboundEndpointId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Dns Resolver Name: %q", id.DnsResolverName),
		fmt.Sprintf("Inbound Endpoint Name: %q", id.InboundEndpointName),
	}
	return fmt.Sprintf("Inbound Endpoint (%s)", strings.Join(components, "\n"))
}
